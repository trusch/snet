package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/trusch/snet/pkg/events"
	"github.com/trusch/snet/pkg/snet"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type coreHandler struct {
	db           *sql.DB
	privateKeys  map[string]ed25519.PrivateKey
	eventManager events.Manager
}

func NewCoreHandler(db *sql.DB, privateKeys []ed25519.PrivateKey) (snet.CoreServer, error) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS messages(
  id VARCHAR PRIMARY KEY,
  type VARCHAR,
  feed VARCHAR,
  content BYTEA,
  created_at TIMESTAMP,
  previous VARCHAR,
  signature VARCHAR,
  sequence INT
);
`)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]ed25519.PrivateKey)
	for _, key := range privateKeys {
		keys[hexencode(key.Public().(ed25519.PublicKey))] = key
	}

	em := events.NewManager(context.Background(), 32)
	return &coreHandler{
		db:           db,
		privateKeys:  keys,
		eventManager: em,
	}, nil
}

func (h *coreHandler) Get(ctx context.Context, req *snet.GetRequest) (*snet.Message, error) {
	row := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(h.db).
		Select(messageColumns()...).
		From("messages").
		Where(squirrel.Eq{"id": req.GetId()}).
		QueryRowContext(ctx)
	msg := &snet.Message{}
	ts := time.Time{}
	err := row.Scan(&msg.Id, &msg.Type, &msg.Feed, &msg.Content, &ts, &msg.Previous, &msg.Signature, &msg.Sequence)
	if err != nil {
		return nil, err
	}
	msg.CreatedAt, err = ptypes.TimestampProto(ts)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (h *coreHandler) Create(ctx context.Context, req *snet.CreateRequest) (msg *snet.Message, err error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	privateKey, ok := h.privateKeys[req.Feed]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "feed not managed by this server")
	}

	latestMessage, err := h.getLastMessage(ctx, tx, req.Feed)
	if err != nil {
		logrus.Warnf("first message? error getting last msg: %v", err)
		latestMessage = &snet.Message{}
	}

	msg = &snet.Message{
		Type:      req.Type,
		Feed:      req.Feed,
		Content:   req.Content,
		CreatedAt: ptypes.TimestampNow(),
		Previous:  latestMessage.Id,
		Sequence:  latestMessage.Sequence + 1,
	}
	bs, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	hash := make([]byte, 32)
	sha3.ShakeSum256(hash, bs)
	msg.Id = hexencode(hash)
	sig := ed25519.Sign(privateKey, hash)
	msg.Signature = hexencode(sig)

	ts, _ := ptypes.Timestamp(msg.CreatedAt)
	_, err = squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("messages").Columns(messageColumns()...).
		Values(msg.Id, msg.Type, msg.Feed, msg.Content, ts, msg.Previous, msg.Signature, msg.Sequence).
		RunWith(tx).
		ExecContext(ctx)
	if err != nil {
		return nil, err
	}
	go h.eventManager.Publish("create", msg)
	return msg, nil
}

func (h *coreHandler) List(req *snet.ListRequest, resp snet.Core_ListServer) error {
	stmt := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(h.db).
		Select(messageColumns()...).
		From("messages")

	if req.Feed != "" {
		stmt = stmt.Where(squirrel.Eq{"feed": req.Feed})
	}
	if req.Type != "" {
		stmt = stmt.Where(squirrel.Eq{"type": req.Type})
	}
	if req.StartSequence > 0 {
		stmt = stmt.Where(squirrel.GtOrEq{"sequence": req.StartSequence})
	}
	if req.EndSequence > 0 {
		stmt = stmt.Where(squirrel.LtOrEq{"sequence": req.EndSequence})
	}
	if req.StartTimestamp != nil {
		ts, err := ptypes.Timestamp(req.StartTimestamp)
		if err != nil {
			return err
		}
		stmt = stmt.Where(squirrel.GtOrEq{"created_at": ts})
	}
	if req.EndTimestamp != nil {
		ts, err := ptypes.Timestamp(req.EndTimestamp)
		if err != nil {
			return err
		}
		stmt = stmt.Where(squirrel.LtOrEq{"created_at": ts})
	}
	rows, err := stmt.QueryContext(resp.Context())
	if err != nil {
		return err
	}

	for rows.Next() {
		msg := &snet.Message{}
		ts := time.Time{}
		err = rows.Scan(&msg.Id, &msg.Type, &msg.Feed, &msg.Content, &ts, &msg.Previous, &msg.Signature, &msg.Sequence)
		if err != nil {
			return err
		}
		msg.CreatedAt, err = ptypes.TimestampProto(ts)
		if err != nil {
			return err
		}
		err = resp.Send(msg)
		if err != nil {
			return err
		}
	}

	if req.Stream {
		id, err := h.eventManager.Subscribe("create", func(event interface{}) error {
			return resp.Send(event.(*snet.Message))
		})
		if err != nil {
			return err
		}
		defer h.eventManager.Unsubscribe("create", id)
		<-resp.Context().Done()
	}

	return nil
}

func (h *coreHandler) getLastMessage(ctx context.Context, tx *sql.Tx, feed string) (msg *snet.Message, err error) {
	row := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		Select(messageColumns()...).
		From("messages").
		Where(squirrel.Eq{"feed": feed}).
		OrderBy("sequence DESC").
		Limit(1).
		QueryRowContext(ctx)
	msg = &snet.Message{}
	ts := time.Time{}
	err = row.Scan(&msg.Id, &msg.Type, &msg.Feed, &msg.Content, &ts, &msg.Previous, &msg.Signature, &msg.Sequence)
	if err != nil {
		return nil, err
	}
	msg.CreatedAt, err = ptypes.TimestampProto(ts)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func messageColumns() []string {
	return []string{"id", "type", "feed", "content", "created_at", "previous", "signature", "sequence"}
}

func hexencode(data []byte) string {
	return fmt.Sprintf("%x", data)
}

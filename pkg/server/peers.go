package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"github.com/trusch/snet/pkg/events"
	"github.com/trusch/snet/pkg/snet"
)

type peersHandler struct {
	db           *sql.DB
	eventManager events.Manager
}

func NewPeersHandler(db *sql.DB, eventManager events.Manager) (snet.PeersServer, error) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS peers(
  id VARCHAR PRIMARY KEY,
  name VARCHAR,
  address VARCHAR,
  feeds VARCHAR[],
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
`)
	if err != nil {
		return nil, err
	}
	return &peersHandler{db, eventManager}, nil
}

func (h *peersHandler) Create(ctx context.Context, req *snet.CreatePeerRequest) (*snet.Peer, error) {
	_, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("peers").Columns(peerColumns()...).
		Values(req.Id, req.Name, req.Address, pq.Array(req.Feeds), time.Now(), time.Now()).
		RunWith(h.db).
		ExecContext(ctx)
	if err != nil {
		return nil, err
	}
	peer, err := h.Get(ctx, &snet.GetPeerRequest{Id: req.Id, Name: req.Name})
	if err != nil {
		return nil, err
	}
	h.eventManager.Publish("peer-create", peer)
	return peer, nil
}

func (h *peersHandler) Get(ctx context.Context, req *snet.GetPeerRequest) (*snet.Peer, error) {
	stmt := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select(peerColumns()...).
		From("peers")
	if req.Id != "" {
		stmt = stmt.Where(squirrel.Eq{"id": req.Id})
	}
	if req.Name != "" {
		stmt = stmt.Where(squirrel.Eq{"name": req.Name})
	}
	row := stmt.RunWith(h.db).QueryRowContext(ctx)
	peer := &snet.Peer{}
	var createdAt, updatedAt time.Time
	err := row.Scan(&peer.Id, &peer.Name, &peer.Address, pq.Array(&peer.Feeds), &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	peer.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		return nil, err
	}
	peer.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
	if err != nil {
		return nil, err
	}
	return peer, nil
}

func (h *peersHandler) Update(ctx context.Context, req *snet.UpdatePeerRequest) (*snet.Peer, error) {
	stmt := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Update("peers").
		Where(squirrel.Eq{"id": req.Id}).
		Set("updated_at", time.Now()).
		RunWith(h.db)
	if req.Name != "" {
		stmt = stmt.Set("name", req.Name)
	}
	if req.Address != "" {
		stmt = stmt.Set("address", req.Address)
	}
	if req.Feeds != nil {
		stmt = stmt.Set("feeds", pq.Array(req.Feeds))
	}
	_, err := stmt.ExecContext(ctx)
	if err != nil {
		return nil, err
	}
	peer, err := h.Get(ctx, &snet.GetPeerRequest{Id: req.Id, Name: req.Name})
	if err != nil {
		return nil, err
	}
	h.eventManager.Publish("peer-update", peer)
	return peer, nil
}

func (h *peersHandler) Delete(ctx context.Context, req *snet.GetPeerRequest) (*snet.Peer, error) {
	peer, err := h.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	stmt := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Delete("peers")
	if req.Id != "" {
		stmt = stmt.Where(squirrel.Eq{"id": req.Id})
	}
	if req.Name != "" {
		stmt = stmt.Where(squirrel.Eq{"name": req.Name})
	}
	_, err = stmt.RunWith(h.db).ExecContext(ctx)
	if err != nil {
		return nil, err
	}
	h.eventManager.Publish("peer-delete", peer)
	return peer, nil
}

func (h *peersHandler) List(_ *snet.ListPeerRequest, resp snet.Peers_ListServer) error {
	stmt := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(h.db).
		Select(peerColumns()...).
		From("peers")

	rows, err := stmt.QueryContext(resp.Context())
	if err != nil {
		return err
	}
	for rows.Next() {
		peer := &snet.Peer{}
		var createdAt, updatedAt time.Time
		err = rows.Scan(&peer.Id, &peer.Name, &peer.Address, pq.Array(&peer.Feeds), &createdAt, &updatedAt)
		if err != nil {
			return err
		}
		peer.CreatedAt, err = ptypes.TimestampProto(createdAt)
		if err != nil {
			return err
		}
		peer.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
		if err != nil {
			return err
		}
		err = resp.Send(peer)
		if err != nil {
			return err
		}
	}
	return nil
}

func peerColumns() []string {
	return []string{"id", "name", "address", "feeds", "created_at", "updated_at"}
}

package fetcher

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/trusch/snet/pkg/events"
	"github.com/trusch/snet/pkg/snet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Fetcher struct {
	db               *sql.DB
	eventManager     events.Manager
	connectionsMutex sync.Mutex
	connections      map[string]*grpc.ClientConn
}

func New(db *sql.DB, eventManager events.Manager) *Fetcher {
	fetcher := &Fetcher{
		db:           db,
		eventManager: eventManager,
		connections:  make(map[string]*grpc.ClientConn),
	}
	return fetcher
}

func (f *Fetcher) Work(ctx context.Context) error {
	peers, err := f.getPeers(ctx)
	if err != nil {
		return err
	}
	for _, peer := range peers {
		f.fetchFromPeer(ctx, peer)
	}
	f.eventManager.Subscribe("peer-create", func(obj interface{}) error {
		peer, ok := obj.(*snet.Peer)
		if !ok {
			return errors.New("can not cast peer")
		}
		f.fetchFromPeer(ctx, peer)
		return nil
	})
	f.eventManager.Subscribe("peer-delete", func(obj interface{}) error {
		peer, ok := obj.(*snet.Peer)
		if !ok {
			return errors.New("can not cast peer")
		}
		f.stopFetching(ctx, peer)
		return nil
	})
	f.eventManager.Subscribe("peer-update", func(obj interface{}) error {
		peer, ok := obj.(*snet.Peer)
		if !ok {
			return errors.New("can not cast peer")
		}
		f.stopFetching(ctx, peer)
		f.fetchFromPeer(ctx, peer)
		return nil
	})
	return nil
}

func (f *Fetcher) stopFetching(ctx context.Context, peer *snet.Peer) {
	f.connectionsMutex.Lock()
	defer f.connectionsMutex.Unlock()
	if conn, ok := f.connections[peer.Address]; ok {
		err := conn.Close()
		if err != nil {
			logrus.Warn(err)
		}
		delete(f.connections, peer.Address)
	}
}

func (f *Fetcher) fetchFromPeer(ctx context.Context, peer *snet.Peer) {
	for _, feed := range peer.Feeds {
		go func() {
			seq := f.getLatestLocalSequenceNumber(ctx, feed)
			err := f.fetchFeed(ctx, peer.Address, feed, seq+1, peer.UseTLS)
			if err != nil {
				logrus.Warn(err)
			}
		}()
	}
}

func (f *Fetcher) dial(ctx context.Context, addr string, useTLS bool) (*grpc.ClientConn, error) {
	f.connectionsMutex.Lock()
	defer f.connectionsMutex.Unlock()
	if conn, ok := f.connections[addr]; ok {
		return conn, nil
	}
	opt := grpc.WithInsecure()
	if useTLS {
		opt = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		}))
	}
	conn, err := grpc.DialContext(ctx, addr, opt)
	if err != nil {
		return nil, err
	}
	f.connections[addr] = conn
	return conn, nil
}

func (f *Fetcher) fetchFeed(ctx context.Context, addr string, feed string, startSequence int64, useTLS bool) error {
	conn, err := f.dial(ctx, addr, useTLS)
	if err != nil {
		return err
	}
	cli := snet.NewCoreClient(conn)
	iter, err := cli.List(ctx, &snet.ListRequest{
		Feed:          feed,
		StartSequence: startSequence,
		Stream:        true,
	})
	if err != nil {
		return err
	}
	for {
		msg, err := iter.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		ts, _ := ptypes.Timestamp(msg.CreatedAt)
		_, err = squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			Insert("messages").Columns(messageColumns()...).
			Values(msg.Id, msg.Type, msg.Feed, msg.Content, ts, msg.Previous, msg.Signature, msg.Sequence).
			RunWith(f.db).
			ExecContext(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Fetcher) getPeers(ctx context.Context) (peers []*snet.Peer, err error) {
	stmt := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(f.db).
		Select(peerColumns()...).
		From("peers")

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		peer := &snet.Peer{}
		var createdAt, updatedAt time.Time
		err = rows.Scan(&peer.Id, &peer.Name, &peer.Address, pq.Array(&peer.Feeds), &createdAt, &updatedAt)
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
		peers = append(peers, peer)
	}
	return peers, nil
}

func (f *Fetcher) getLatestLocalSequenceNumber(ctx context.Context, feed string) int64 {
	var seq int64
	err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(f.db).
		Select("sequence").
		From("messages").
		Where(squirrel.Eq{"feed": feed}).
		OrderBy("sequence DESC").
		Limit(1).
		QueryRowContext(ctx).
		Scan(&seq)
	if err != nil {
		logrus.Warn(err)
	}
	return seq

}

func messageColumns() []string {
	return []string{"id", "type", "feed", "content", "created_at", "previous", "signature", "sequence"}
}

func peerColumns() []string {
	return []string{"id", "name", "address", "feeds", "created_at", "updated_at"}
}

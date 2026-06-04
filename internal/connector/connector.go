package connector

import (
	"context"
	"fmt"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

type Conn interface {
	Close(ctx context.Context) error
	ReceiveMessage(ctx context.Context) (pgproto3.BackendMessage, error)
}
 
type Dialer func(ctx context.Context, url string) (Conn, error)

func defaultDialer(ctx context.Context, url string) (Conn, error) {
	return pgconn.Connect(ctx, url)
}

type Connector struct {
	conn     Conn
	slotName string
	pubName  string
	rawConn  *pgconn.PgConn
}

type Option func(*Connector)

func New(ctx context.Context, databaseURL string, slotName string, pubName string, opts ...Option) (*Connector, error) {
	return newWithDialer(ctx, databaseURL, slotName, pubName, defaultDialer, opts...)
}

func newWithDialer(ctx context.Context, databaseURL string, slotName string, pubName string, dialer Dialer, opts ...Option) (*Connector, error) {
	conn, err := dialer(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	c := &Connector{
		conn:     conn,
		slotName: slotName,
		pubName:  pubName,
	}

	if raw, ok := conn.(*pgconn.PgConn); ok {
		c.rawConn = raw
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Connector) CreateSlot(ctx context.Context) error {
	if c.rawConn == nil {
		return fmt.Errorf("replication requires a real postgres connection")
	}
	_, err := pglogrepl.CreateReplicationSlot(
		ctx,
		c.rawConn,
		c.slotName,
		"pgoutput",
		pglogrepl.CreateReplicationSlotOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to create replication slot: %w", err)
	}
	fmt.Println("replication slot created:", c.slotName)
	return nil
}

func (c *Connector) Start(ctx context.Context) error {
	if c.rawConn == nil {
		return fmt.Errorf("replication requires a real postgres connection")
	}
	err := pglogrepl.StartReplication(
		ctx,
		c.rawConn,
		c.slotName,
		0,
		pglogrepl.StartReplicationOptions{
			PluginArgs: []string{
				"proto_version '1'",
				fmt.Sprintf("publication_names '%s'", c.pubName),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start replication: %w", err)
	}
	fmt.Println("replication started")
	return nil
}

func (c *Connector) Close(ctx context.Context) {
	c.conn.Close(ctx)
}

func (c *Connector) GetConn() Conn {
	return c.conn
}

func (c *Connector) GetRawConn() *pgconn.PgConn {
	return c.rawConn
}

// test for demo video
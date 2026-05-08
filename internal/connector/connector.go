package connector

import (
	"context"
	"fmt"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
)

type Connector struct {
	conn *pgconn.PgConn
	slotName string
	pubName string
}

func New(databaseURL string, slotName string, pubName string) (*Connector, error) {
	conn, err := pgconn.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Connector{
		conn: conn,
		slotName: slotName,
		pubName: pubName,
	}, nil
}

func (c *Connector) CreateSlot() error {
	_, err := pglogrepl.CreateReplicationSlot(
		context.Background(),
		c.conn,
		c.slotName,
		"pgoutput",
		pglogrepl.CreateReplicationSlotOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to create replication slot: %w", err)
	}
	fmt.Println("Replication slot created successfully", c.slotName)
	return nil
}

func (c *Connector) Start() error {
	err := pglogrepl.StartReplication(
		context.Background(),
		c.conn,
		c.slotName,
		0,
		pglogrepl.StartReplicationOptions{
			PluginArgs: []string{
				"proto_version '1'", fmt.Sprintf("publication_names '%s'", c.pubName),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start replication: %w", err)
	}
	fmt.Println("Replication started successfully")
	return nil
}
func (c *Connector) Close() {
	c.conn.Close(context.Background())
}

func (c *Connector) GetConn() *pgconn.PgConn {
	return c.conn
}
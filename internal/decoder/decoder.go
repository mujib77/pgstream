package decoder

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/mujib77/pgstream/internal/models"
)

type Decoder struct {
	conn      *pgconn.PgConn
	relations map[uint32]*pglogrepl.RelationMessage
}

func New(conn *pgconn.PgConn) *Decoder {
	return &Decoder{
		conn:      conn,
		relations: make(map[uint32]*pglogrepl.RelationMessage),
	}
}

func (d *Decoder) NextEvent(ctx context.Context) (*models.WalEvent, error) {
	rawMsg, err := d.conn.ReceiveMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message: %w", err)
	}

	msg, ok := rawMsg.(*pgproto3.CopyData)
	if !ok {
		return nil, nil
	}

	if len(msg.Data) == 0 {
		return nil, nil
	}

	switch msg.Data[0] {
	case pglogrepl.PrimaryKeepaliveMessageByteID:
		pka, err := pglogrepl.ParsePrimaryKeepaliveMessage(msg.Data[1:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse keepalive: %w", err)
		}
		if pka.ReplyRequested {
			err = pglogrepl.SendStandbyStatusUpdate(ctx, d.conn,
				pglogrepl.StandbyStatusUpdate{
					WALWritePosition: pka.ServerWALEnd,
					ClientTime:       time.Now(),
				})
			if err != nil {
				return nil, fmt.Errorf("failed to send standby status: %w", err)
			}
		}
		return nil, nil

	case pglogrepl.XLogDataByteID:
		xld, err := pglogrepl.ParseXLogData(msg.Data[1:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse xlog data: %w", err)
		}

		logMsg, err := pglogrepl.Parse(xld.WALData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse wal data: %w", err)
		}

		return d.handleMessage(logMsg, xld.WALStart)
	}

	return nil, nil
}

func (d *Decoder) handleMessage(
	msg pglogrepl.Message,
	lsn pglogrepl.LSN,
) (*models.WalEvent, error) {
	switch m := msg.(type) {
	case *pglogrepl.RelationMessage:
		d.relations[m.RelationID] = m
		return nil, nil

	case *pglogrepl.InsertMessage:
		rel, ok := d.relations[m.RelationID]
		if !ok {
			return nil, nil
		}
		data := decodeColumns(rel, m.Tuple)
		return &models.WalEvent{
			Table:     rel.RelationName,
			Operation: models.OperationInsert,
			Data:      data,
			LSN:       lsn.String(),
		}, nil

	case *pglogrepl.UpdateMessage:
		rel, ok := d.relations[m.RelationID]
		if !ok {
			return nil, nil
		}
		newData := decodeColumns(rel, m.NewTuple)
		var oldData map[string]interface{}
		if m.OldTuple != nil {
			oldData = decodeColumns(rel, m.OldTuple)
		}
		return &models.WalEvent{
			Table:     rel.RelationName,
			Operation: models.OperationUpdate,
			Data:      newData,
			OldData:   oldData,
			LSN:       lsn.String(),
		}, nil

	case *pglogrepl.DeleteMessage:
		rel, ok := d.relations[m.RelationID]
		if !ok {
			return nil, nil
		}
		var oldData map[string]interface{}
		if m.OldTuple != nil {
			oldData = decodeColumns(rel, m.OldTuple)
		}
		return &models.WalEvent{
			Table:     rel.RelationName,
			Operation: models.OperationDelete,
			OldData:   oldData,
			LSN:       lsn.String(),
		}, nil
	}

	return nil, nil
}

func decodeColumns(
	rel *pglogrepl.RelationMessage,
	tuple *pglogrepl.TupleData,
) map[string]interface{} {
	data := make(map[string]interface{})
	if tuple == nil {
		return data
	}
	for i, col := range tuple.Columns {
		if i >= len(rel.Columns) {
			break
		}
		colName := rel.Columns[i].Name
		switch col.DataType {
		case 'n':
			data[colName] = nil
		case 't':
			data[colName] = string(col.Data)
		}
	}
	return data
}
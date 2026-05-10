package handler

import (
	"testing"

	"github.com/mujib77/pgstream/internal/models"
)

func TestHandleInsert(t *testing.T) {
	h := New()  

	event := &models.WalEvent{
		Table:	 "users",
		Operation: models.OperationInsert,
		Data: map[string]interface{}{
			"id":   1,
			"name": "Mujib",
			"email": "mujib@example.com",
		},
		LSN: "0/16C752F8",
	}

	err := h.Handle(event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandleUpdate(t *testing.T) {
	h := New()

	event := &models.WalEvent{
		Table:	 "users",
		Operation: models.OperationUpdate,
		Data: map[string]interface{}{
			"name": "Alice",
		},
		OldData: map[string]interface{}{
			"name": "Mujib",
		},
		LSN: "0/16C75410",
	}

	err := h.Handle(event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
		}

func TestHandleDelete(t *testing.T) {
	h := New()

	event := &models.WalEvent{
		Table:	 "users",
		Operation: models.OperationDelete,
		OldData: map[string]interface{}{
			"id":   1,
			"name": "Mujib",
		},
		LSN: "0/16C754A8",
	}

	err := h.Handle(event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHandleNilEvent(t *testing.T) {
	h := New()

	err := h.Handle(nil)
	if err != nil {
		t.Errorf("expected no error for nil event, got %v", err)
	}
}
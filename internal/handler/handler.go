package handler

import (
	"encoding/json"
	"fmt"

	"github.com/mujib77/pgstream/internal/models"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(event *models.WalEvent) error {
	if event == nil {
		return nil
	}

	data, err := json.MarshalIndent(event.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	switch event.Operation {
	case models.OperationInsert:
		fmt.Printf("\n[INSERT] table=%s lsn=%s\n", event.Table, event.LSN)
		fmt.Printf("data=%s\n", string(data))

	case models.OperationUpdate:
		fmt.Printf("\n[UPDATE] table=%s lsn=%s\n", event.Table, event.LSN)
		if event.OldData != nil {
			old, _ := json.MarshalIndent(event.OldData, "", "  ")
			fmt.Printf("old=%s\n", string(old))
		}
		fmt.Printf("new=%s\n", string(data))

	case models.OperationDelete:
		fmt.Printf("\n[DELETE] table=%s lsn=%s\n", event.Table, event.LSN)
		if event.OldData != nil {
			old, _ := json.MarshalIndent(event.OldData, "", "  ")
			fmt.Printf("data=%s\n", string(old))
		}
	}

	return nil
}
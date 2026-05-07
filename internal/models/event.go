package models

type OperationType string

const (
	OperationInsert OperationType = "INSERT"
	OperationUpdate OperationType = "UPDATE"
	OperationDelete OperationType = "DELETE"
)

type WalEvent struct {
	Table string
	Operation OperationType
	Data map[string]interface{}
	OldData map[string]interface{}
	LSN string
}
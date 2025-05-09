package models

type Expression struct {
	ID         string  `gorm:"primaryKey;column:id"`
	OwnerID    string  `gorm:"column:owner_id"`
	RootTaskID string  `gorm:"column:root_task_id;not null"`
	Status     Status  `gorm:"column:status;not null"`
	Result     float64 `gorm:"column:result;not null"`
}
type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusError   Status = "error"
	StatusPending Status = "pending"
	StatusWaiting Status = "waiting"
	StatusReady   Status = "ready"
	StatusRunning Status = "running"
)

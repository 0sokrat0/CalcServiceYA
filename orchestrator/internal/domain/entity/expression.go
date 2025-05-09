package entity

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

type Expression struct {
	ID         string
	OwnerID    string
	RootTaskID string
	Status     Status
	Result     float64
}

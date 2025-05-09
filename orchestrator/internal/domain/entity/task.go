package entity

type Task struct {
	ID            string
	Arg1          float64
	Arg2          float64
	LeftTaskID    string
	RightTaskID   string
	Operation     string
	OperationTime int64
	Result        float64
	Status        Status
}

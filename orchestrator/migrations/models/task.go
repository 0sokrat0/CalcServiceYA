package models

type Task struct {
	ID            string  `gorm:"primaryKey;column:id"`
	Arg1          float64 `gorm:"column:arg1;not null"`
	Arg2          float64 `gorm:"column:arg2;not null"`
	LeftTaskID    string  `gorm:"column:left_task_id"`
	RightTaskID   string  `gorm:"column:right_task_id"`
	Operation     string  `gorm:"column:operation;not null"`
	OperationTime int64   `gorm:"column:operation_time;not null"`
	Result        float64 `gorm:"column:result;not null"`
	Status        Status  `gorm:"column:status;not null"`
}

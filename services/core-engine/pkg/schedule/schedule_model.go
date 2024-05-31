//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

package schedule

import "context"

const (
	// Add new value according to your usecase
	UpdateCampaign Job = "UPDATE"
)

type (
	Job string

	JobHandler func(context.Context, Schedule) error

	JobHandlerMapping map[Job]JobHandler

	Schedule struct {
		ID            int    `gorm:"column:id"`
		Job           Job    `gorm:"column:job"`
		ExecutionTime int64  `gorm:"column:execution_time"` // Unix UTC
		IsDone        bool   `gorm:"column:is_done"`
		Payload       []byte `gorm:"column:payload"`
	}
)

//counterfeiter:generate -o ./mock . Scheduler
type Scheduler interface {
	NewSchedule(job Job, executionTime int64, payload any) error
}

func (Schedule) TableName() string {
	return "schedule"
}

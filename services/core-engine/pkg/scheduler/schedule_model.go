package scheduler

const (
	// Add new value according to your usecase
	UpdateCampaign Job = "UPDATE_CAMPAIGN"
)

type (
	Job string

	JobFunc map[Job]func([]byte) error

	Scheduler interface {
		NewSchedule(job Job, executionTime int64, payload any) error
		RestoreScheduler() error
	}

	Schedule struct {
		ID            int    `gorm:"column:id"`
		Job           Job    `gorm:"column:job"`
		ExecutionTime int64  `gorm:"column:execution_time"` // Unix UTC
		IsDone        bool   `gorm:"column:is_done"`
		Payload       []byte `gorm:"column:payload"`
	}
)

func (Schedule) TableName() string {
	return "schedule"
}

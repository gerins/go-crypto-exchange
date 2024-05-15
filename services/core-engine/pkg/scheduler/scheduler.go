package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gerins/log"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type scheduler struct {
	redis          *redis.Client
	db             *gorm.DB
	jobFuncMapping JobFunc
}

func NewScheduler(redis *redis.Client, writeDB *gorm.DB) *scheduler {
	newScheduler := &scheduler{
		redis: redis,
		db:    writeDB,
	}

	newScheduler.jobFuncMapping = JobFunc{
		UpdateCampaign: newScheduler.updateCampaign,
	}

	return newScheduler
}

func (s *scheduler) NewSchedule(job Job, executionTime int64, payload any) error {
	task, found := s.jobFuncMapping[job]
	if !found {
		return errors.New("job specification not found")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	newSchedule := Schedule{
		Job:           job,
		ExecutionTime: executionTime,
		Payload:       payloadJSON,
	}

	// Insert to table schedule
	if err := s.db.Save(&newSchedule).Error; err != nil {
		return err
	}

	go s.execute(newSchedule, task)

	return nil
}

func (s *scheduler) RestoreScheduler() error {
	schedules := make([]Schedule, 0)
	if err := s.db.Find(&schedules, "is_done = false").Error; err != nil {
		return err
	}

	for _, schedule := range schedules {
		task := s.jobFuncMapping[schedule.Job]
		go s.execute(schedule, task)
	}

	return nil
}

func (s *scheduler) execute(schedule Schedule, task func([]byte) error) {
	sleepDuration := schedule.ExecutionTime - time.Now().Unix()
	if sleepDuration <= 0 {
		sleepDuration = 0
	}

	time.Sleep(time.Duration(sleepDuration) * time.Second)

	// Make sure only one process processing the task
	key := fmt.Sprintf("service-scheduler:%v", schedule.ID)
	defer s.redis.Expire(context.TODO(), key, 10*time.Second)
	if countVal, _ := s.redis.Incr(context.TODO(), key).Result(); countVal != 1 {
		return // There is already other process
	}

	// Run task
	if err := task(schedule.Payload); err != nil {
		// TODO case failed
		log.Error(err)
	}

	schedule.IsDone = true
	if err := s.db.Save(&schedule).Error; err != nil {
		log.Error(err)
	}
}

func (s *scheduler) updateCampaign(payloadJSON []byte) error {
	return nil
}

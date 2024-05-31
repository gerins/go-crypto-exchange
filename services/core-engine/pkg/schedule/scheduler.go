package schedule

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
	redis      *redis.Client
	readDB     *gorm.DB
	writeDB    *gorm.DB
	jobMapping JobHandlerMapping
}

func New(redis *redis.Client, readDB *gorm.DB, writeDB *gorm.DB) *scheduler {
	newScheduler := &scheduler{
		redis:   redis,
		readDB:  readDB,
		writeDB: writeDB,
	}

	return newScheduler
}

func (s *scheduler) RegisterHandler(jobMapping JobHandlerMapping) {
	s.jobMapping = jobMapping
}

func (s *scheduler) Restore() error {
	schedules := make([]Schedule, 0)
	if err := s.readDB.Find(&schedules, "is_done = false").Error; err != nil {
		return err
	}

	for _, schedule := range schedules {
		if job, found := s.jobMapping[schedule.Job]; found {
			go s.execute(schedule, job)
		}
	}

	return nil
}

func (s *scheduler) NewSchedule(jobCode Job, executionTime int64, payload any) error {
	job, found := s.jobMapping[jobCode]
	if !found {
		return errors.New("job specification not found")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	newSchedule := Schedule{
		Job:           jobCode,
		ExecutionTime: executionTime,
		Payload:       payloadJSON,
	}

	// Insert to table schedule
	if err := s.writeDB.Save(&newSchedule).Error; err != nil {
		return err
	}

	go s.execute(newSchedule, job)

	return nil
}

func (s *scheduler) execute(schedule Schedule, job JobHandler) {
	sleepDuration := schedule.ExecutionTime - time.Now().Unix()
	if sleepDuration <= 0 {
		sleepDuration = 0
	}

	time.Sleep(time.Duration(sleepDuration) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make sure only one process processing the job
	key := fmt.Sprintf("service-scheduler:%v", schedule.ID)
	defer s.redis.Expire(ctx, key, 10*time.Second)
	if countVal, _ := s.redis.Incr(ctx, key).Result(); countVal != 1 {
		return // There is already other process
	}

	// Init logging
	log := log.NewRequest()
	log.ReqBody = schedule
	log.Method = "SCHEDULER"
	log.URL = string(schedule.Job)
	defer log.Save()

	// Run job
	if err := job(log.SaveToContext(ctx), schedule); err != nil {
		log.Error(err)
	}

	schedule.IsDone = true
	if err := s.writeDB.Save(&schedule).Error; err != nil {
		log.Error(err)
	}
}

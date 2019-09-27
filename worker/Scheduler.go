package worker

import (
	"github.com/c0ding/crontab/common"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent

	jobPlanTable map[string]*common.JobSchedulePlan // 任务调度计划表
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (s *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {

	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExisted bool
		err             error

	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan

	case common.JOB_EVENT_DELETE:
		if jobSchedulePlan, jobExisted = s.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(s.jobPlanTable,jobEvent.Job.Name)
		}
	}

	//TODO
}

func (s *Scheduler) schedulerLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)

	scheduleAfter =

	for {
		select {
		case jobEvent = <-s.jobEventChan:
			s.handleJobEvent(jobEvent)
		}
	}
}

func (s *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	s.jobEventChan <- jobEvent
}

func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
	}
	go G_scheduler.schedulerLoop()
	return
}

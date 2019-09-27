package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`      //任务名
	Commond  string `json:"commond"`   //shell命令
	CronExpr string `json:"cron_expr"` // 任务表达式

}

// HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type JobEvent struct {
	EventType int //save , delete
	Job       *Job
}

type JobSchedulePlan struct {
	Job      *Job
	Expr     *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time            //下次调度时间
}

func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {

	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}

// 应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 1, 定义一个response
	var (
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	// 2, 序列化json
	resp, err = json.Marshal(response)
	return
}

// 反序列化
func UnpackJob(value []byte) (ret *Job, err error) {

	var (
		job *Job
	)

	job = new(Job)
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return

}

// 任务变化事件，有两种 1）更新   2)删除任务
func BuildJobEvent(job *Job, eventType int) *JobEvent {

	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// 提取任务名
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

func ExtractKillerName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_KILLER_DIR)
}

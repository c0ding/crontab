package worker

import (
	"context"
	"fmt"
	"github.com/c0ding/crontab/common"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

// 任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	// 单例
	G_jobMgr *JobMgr
)

func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp  *clientv3.GetResponse
		keyValue *mvccpb.KeyValue
		job      *common.Job
		jobEvent *common.JobEvent

		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
	)
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, keyValue = range getResp.Kvs {

		if job, err = common.UnpackJob(keyValue.Value); err == nil {
			// 给调度模块推送 任务事件
			//1，先组成任务事件
			jobEvent = common.BuildJobEvent(job, common.JOB_EVENT_SAVE)
			//TODO:2，推送给调度模块
			fmt.Println(jobEvent)
		}
	}

	func() {
		watchStartRevision = getResp.Header.Revision + 1

		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision))

		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(job, common.JOB_EVENT_SAVE)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					jobEvent = common.BuildJobEvent(job, common.JOB_EVENT_DELETE)
				}

				//TODO:2，推送给调度模块
				// 变化推给scheduler
				//G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()

	return
}

func InitJobMgr() (err error) {

	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, // 连接超时
	}
	if client, err = clientv3.New(config); err != nil {
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	G_jobMgr.watchJobs()
	return
}

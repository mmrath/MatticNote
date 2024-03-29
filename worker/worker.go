package worker

import (
	"github.com/MatticNote/MatticNote/config"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	Enqueue *work.Enqueuer
	Worker  *work.WorkerPool
	Client  *work.Client
)

const (
	workerNamespace = "mn_worker"
)

const (
	JobInboxProcess = "inbox_process"
	JobDeliver      = "deliver"
	JobExportData   = "export_data"
	JobNotePreJob   = "note_pre_job"
)

func getRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   config.Config.Job.MaxIdle,
		MaxActive: config.Config.Job.MaxActive,
		Wait:      true,
		Dial:      config.GetRedisDial,
	}
}

func InitEnqueue() {
	redisPool := getRedisPool()
	Enqueue = work.NewEnqueuer(workerNamespace, redisPool)
	Client = work.NewClient(workerNamespace, redisPool)
}

func InitWorker() {
	redisPool := getRedisPool()
	Worker = work.NewWorkerPool(Context{}, uint(config.Config.Job.MaxActive), workerNamespace, redisPool)

	Worker.JobWithOptions(
		JobInboxProcess,
		work.JobOptions{
			Priority: 50,
			MaxFails: 10,
		},
		(*Context).ProcessInbox,
	)
	Worker.JobWithOptions(
		JobDeliver,
		work.JobOptions{
			Priority: 40,
			MaxFails: 10,
		},
		(*Context).Deliver,
	)
	Worker.JobWithOptions(
		JobExportData,
		work.JobOptions{
			Priority: 30,
			MaxFails: 1,
		},
		(*Context).ExportData,
	)
	Worker.JobWithOptions(
		JobNotePreJob,
		work.JobOptions{
			Priority: 20,
			MaxFails: 1,
		},
		(*Context).NotePreJob,
	)
}

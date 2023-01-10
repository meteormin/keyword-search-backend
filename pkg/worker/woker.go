package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

type JobStatus string

var ctx = context.Background()

const (
	SUCCESS JobStatus = "success"
	FAIL    JobStatus = "fail"
	WAIT    JobStatus = "wait"
)

type Job struct {
	JobId     string              `json:"job_id"`
	Status    JobStatus           `json:"status"`
	Closure   func(job Job) error `json:"-"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

func newJob(jobId string, closure func(job Job) error) Job {
	return Job{
		JobId:   jobId,
		Closure: closure,
		Status:  WAIT,
	}
}

type Queue interface {
	Enqueue(job Job)
	Dequeue() Job
	Clear()
	Setup() int
	GetChan() Job
}

type JobQueue struct {
	queue   []Job
	jobChan chan Job
}

func NewQueue() Queue {
	return &JobQueue{
		queue:   make([]Job, 0),
		jobChan: make(chan Job),
	}
}

func (q *JobQueue) GetChan() Job {
	return <-q.jobChan
}

func (q *JobQueue) Setup() int {
	count := len(q.queue)
	q.jobChan = make(chan Job, count)

	for _, j := range q.queue {
		if &j != nil {
			q.jobChan <- j
		}
	}

	return count
}

func (q *JobQueue) Enqueue(job Job) {
	q.queue = append(q.queue, job)
}

func (q *JobQueue) Dequeue() Job {
	job := q.queue[0]
	q.queue = q.queue[1:]
	return job
}

func (q *JobQueue) Clear() {
	q.queue = make([]Job, 0)
	q.jobChan = make(chan Job)
}

type Worker interface {
	GetName() string
	Start()
	Stop()
	AddJob(job Job)
	SaveJob(key string, job Job)
	GetJob(key string) (*Job, error)
}

type JobWorker struct {
	Name     string
	queue    Queue
	jobChan  chan Job
	quitChan chan bool
	redis    *redis.Client
}

func (w *JobWorker) SaveJob(key string, job Job) {
	bytes, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}

	err = w.redis.Set(ctx, key, string(bytes), time.Minute).Err()
	if err != nil {
		panic(err)
	}
}

func (w *JobWorker) GetJob(key string) (*Job, error) {
	val, err := w.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		var convJob *Job
		bytes := []byte(val)
		err = json.Unmarshal(bytes, convJob)
		if err != nil {
			return nil, err
		}
		return convJob, nil
	}
}

func NewWorker(name string, redis *redis.Client) Worker {
	return &JobWorker{
		Name:     name,
		queue:    NewQueue(),
		jobChan:  make(chan Job),
		quitChan: make(chan bool),
		redis:    redis,
	}
}

func (w *JobWorker) GetName() string {
	return w.Name
}

func (w *JobWorker) Start() {
	count := w.queue.Setup()
	w.jobChan = make(chan Job, count)

	go func() {
		for {
			w.jobChan <- w.queue.GetChan()
			select {
			case job := <-w.jobChan:
				key := fmt.Sprintf("%s.%s", w.Name, job.JobId)
				convJob, err := w.GetJob(key)
				if err != nil {
					panic(err)
				}

				if convJob != nil {
					if convJob.Status != SUCCESS {
						continue
					}
				}

				job.CreatedAt = time.Now()
				log.Printf("worker %s, job %s \n", w.Name, job.JobId)
				err = job.Closure(job)
				if err != nil {
					job.Status = FAIL
				} else {
					job.Status = SUCCESS
				}
				log.Printf("end job id %s status %s", job.JobId, job.Status)
				time.Sleep(time.Second * 3)

				job.UpdatedAt = time.Now()

				w.SaveJob(key, job)
			case <-w.quitChan:
				log.Printf("worker %s stopping\n", w.Name)
				return
			}
		}

	}()
}

func (w *JobWorker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}

func (w *JobWorker) AddJob(job Job) {
	w.queue.Enqueue(job)
}

type Dispatcher interface {
	Dispatch(jobId string, closure func(j Job) error)
	Run()
	SelectWorker(name string) Dispatcher
	GetWorkers() []Worker
	GetRedis() *redis.Client
}

type JobDispatcher struct {
	Workers    []Worker
	workerPool chan chan Job
	worker     Worker
	Redis      *redis.Client
}

type Option struct {
	Name string
}

type DispatcherOption struct {
	WorkerOptions []Option
	Redis         *redis.Client
}

func NewDispatcher(opt DispatcherOption) Dispatcher {
	workers := make([]Worker, 0)

	for _, o := range opt.WorkerOptions {
		workers = append(workers, NewWorker(o.Name, opt.Redis))
	}

	return &JobDispatcher{
		Workers:    workers,
		workerPool: make(chan chan Job, len(workers)),
		worker:     nil,
		Redis:      opt.Redis,
	}
}

func (d *JobDispatcher) GetRedis() *redis.Client {
	return d.Redis
}

func (d *JobDispatcher) GetWorkers() []Worker {
	return d.Workers
}

func (d *JobDispatcher) SelectWorker(name string) Dispatcher {
	if name == "" {
		for _, w := range d.Workers {
			if w.GetName() == "default" {
				d.worker = w
			}
		}

	}

	for _, w := range d.Workers {
		if w.GetName() == name {
			d.worker = w
		}
	}

	return d
}

func (d *JobDispatcher) Dispatch(jobId string, closure func(j Job) error) {
	if d.worker == nil {
		for _, w := range d.Workers {
			if w.GetName() == "default" {
				d.worker = w
			}
		}
	}

	d.worker.AddJob(newJob(jobId, closure))
}

func (d *JobDispatcher) Run() {
	for _, w := range d.Workers {
		w.Start()
	}
}

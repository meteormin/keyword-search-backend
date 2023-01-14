package worker

import (
	"context"
	"encoding/json"
	"errors"
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
const DefaultWorker = "default"

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
	Enqueue(job Job) error
	Dequeue() (*Job, error)
	Clear()
}

type JobQueue struct {
	queue       []Job
	jobChan     chan Job
	maxJobCount int
}

func NewQueue(maxJobCount int) Queue {
	return &JobQueue{
		queue:       make([]Job, 0),
		jobChan:     make(chan Job, maxJobCount),
		maxJobCount: maxJobCount,
	}
}

func (q *JobQueue) Enqueue(job Job) error {
	if q.maxJobCount > len(q.queue) {
		q.queue = append(q.queue, job)
		q.jobChan <- job
		return nil
	}

	return errors.New("can't enqueue job queue: over queue size")
}

func (q *JobQueue) Dequeue() (*Job, error) {
	if len(q.queue) == 0 {
		job := <-q.jobChan
		return &job, nil
	}

	job := q.queue[0]
	q.queue = q.queue[1:]
	jobChan := <-q.jobChan
	if job.JobId == jobChan.JobId {
		return &jobChan, nil
	}

	return nil, errors.New("can't match job id")
}

func (q *JobQueue) Clear() {
	q.queue = make([]Job, 0)
	q.jobChan = make(chan Job)
}

type Worker interface {
	GetName() string
	Start()
	Stop()
	AddJob(job Job) error
}

type JobWorker struct {
	Name        string
	queue       Queue
	jobChan     chan Job
	quitChan    chan bool
	redis       func() *redis.Client
	maxJobCount int
}

func saveJob(r *redis.Client, key string, job Job) {
	bytes, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}

	err = r.Set(ctx, key, string(bytes), time.Minute).Err()
	if err != nil {
		panic(err)
	}
}

func getJob(r *redis.Client, key string) (*Job, error) {
	val, err := r.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		var convJob *Job
		bytes := []byte(val)
		err = json.Unmarshal(bytes, &convJob)
		if err != nil {
			return nil, err
		}
		return convJob, nil
	}
}

func NewWorker(name string, redis func() *redis.Client, maxJobCount int) Worker {
	return &JobWorker{
		Name:        name,
		queue:       NewQueue(maxJobCount),
		jobChan:     make(chan Job, maxJobCount),
		quitChan:    make(chan bool),
		redis:       redis,
		maxJobCount: maxJobCount,
	}
}

func (w *JobWorker) GetName() string {
	return w.Name
}

func (w *JobWorker) Start() {
	go func() {
		for {
			r := w.redis()
			jobChan, err := w.queue.Dequeue()
			if err != nil {
				log.Print(err)
				continue
			}

			w.jobChan <- *jobChan
			select {
			case job := <-w.jobChan:
				key := fmt.Sprintf("%s.%s", w.Name, job.JobId)
				convJob, err := getJob(r, key)
				if err != nil {
					log.Println(err)
				}

				if convJob != nil {
					if convJob.Status != SUCCESS {
						err := w.queue.Enqueue(job)
						if err != nil {
							log.Println(err)
							continue
						}
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

				saveJob(r, key, job)
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

func (w *JobWorker) AddJob(job Job) error {
	err := w.queue.Enqueue(job)
	if err != nil {
		return err
	}
	return nil
}

type Dispatcher interface {
	Dispatch(jobId string, closure func(j Job) error) error
	Run()
	SelectWorker(name string) Dispatcher
	GetWorkers() []Worker
	GetRedis() func() *redis.Client
	AddWorker(option Option)
	RemoveWorker(nam string)
}

type JobDispatcher struct {
	workers    []Worker
	workerPool chan chan Job
	worker     Worker
	Redis      func() *redis.Client
}

type Option struct {
	Name        string
	MaxJobCount int
}

type DispatcherOption struct {
	WorkerOptions []Option
	Redis         func() *redis.Client
}

var defaultWorkerOption = []Option{
	{
		Name:        DefaultWorker,
		MaxJobCount: 10,
	},
}

func NewDispatcher(opt DispatcherOption) Dispatcher {
	workers := make([]Worker, 0)

	if len(opt.WorkerOptions) == 0 {
		opt.WorkerOptions = defaultWorkerOption
	}

	for _, o := range opt.WorkerOptions {
		workers = append(workers, NewWorker(o.Name, opt.Redis, o.MaxJobCount))
	}

	return &JobDispatcher{
		workers:    workers,
		workerPool: make(chan chan Job, len(workers)),
		worker:     nil,
		Redis:      opt.Redis,
	}
}

func (d *JobDispatcher) AddWorker(option Option) {
	d.workers = append(d.workers, NewWorker(option.Name, d.Redis, option.MaxJobCount))
}

func (d *JobDispatcher) RemoveWorker(name string) {
	var rmIndex *int = nil
	for i, worker := range d.workers {
		if worker.GetName() == name {
			rmIndex = &i
		}
	}

	if rmIndex != nil {
		d.workers = append(d.workers[:*rmIndex], d.workers[*rmIndex+1:]...)
	}
}

func (d *JobDispatcher) GetRedis() func() *redis.Client {
	return d.Redis
}

func (d *JobDispatcher) GetWorkers() []Worker {
	return d.workers
}

func (d *JobDispatcher) SelectWorker(name string) Dispatcher {
	if name == "" {
		for _, w := range d.workers {
			if w.GetName() == "default" {
				d.worker = w
			}
		}

	}

	for _, w := range d.workers {
		if w.GetName() == name {
			d.worker = w
		}
	}

	return d
}

func (d *JobDispatcher) Dispatch(jobId string, closure func(j Job) error) error {
	if d.worker == nil {
		for _, w := range d.workers {
			if w.GetName() == DefaultWorker {
				d.worker = w
			}
		}
	}

	err := d.worker.AddJob(newJob(jobId, closure))
	if err != nil {
		return err
	}

	return nil
}

func (d *JobDispatcher) Run() {
	for _, w := range d.workers {
		w.Start()
	}
}

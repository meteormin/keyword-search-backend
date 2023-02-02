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
	JobId     string               `json:"job_id"`
	Status    JobStatus            `json:"status"`
	Closure   func(job *Job) error `json:"-"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

func (j *Job) Marshal() (string, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return "", err
	}

	return string(marshal), nil
}

func (j *Job) UnMarshal(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), &j)
	if err != nil {
		return err
	}
	return nil
}

func newJob(jobId string, closure func(job *Job) error) Job {
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
	Count() int
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

func (q *JobQueue) Count() int {
	return len(q.queue)
}

type Worker interface {
	GetName() string
	Run()
	Stop()
	AddJob(job Job) error
	IsRunning() bool
	MaxJobCount() int
	JobCount() int
	BeforeJob(fn func(j *Job))
	AfterJob(fn func(j *Job, err error))
}

type JobWorker struct {
	Name        string
	queue       Queue
	jobChan     chan Job
	quitChan    chan bool
	redis       func() *redis.Client
	maxJobCount int
	isRunning   bool
	beforeJob   func(j *Job)
	afterJob    func(j *Job, err error)
}

func saveJob(r *redis.Client, key string, job Job) {
	jsonJob, err := job.Marshal()
	if err != nil {
		panic(err)
	}

	err = r.Set(ctx, key, jsonJob, time.Minute).Err()
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
		err = convJob.UnMarshal(val)

		if err != nil {
			return nil, err
		}

		return convJob, nil
	}
}

func work(w *JobWorker, job *Job) {
	job.CreatedAt = time.Now()

	log.Printf("worker %s, job %s \n", w.Name, job.JobId)
	if w.beforeJob != nil {
		w.beforeJob(job)
	}

	err := job.Closure(job)
	if err != nil {
		job.Status = FAIL
	} else {
		job.Status = SUCCESS
	}

	job.UpdatedAt = time.Now()

	if w.afterJob != nil {
		w.afterJob(job, err)
	}

	jsonJob, err := job.Marshal()
	if err != nil {
		log.Print(err)
	}

	log.Printf("end job: %s", jsonJob)
}

type Config struct {
	Name        string
	Redis       func() *redis.Client
	MaxJobCount int
	BeforeJob   func(j *Job)
	AfterJob    func(j *Job, err error)
}

func NewWorker(cfg Config) Worker {
	return &JobWorker{
		Name:        cfg.Name,
		queue:       NewQueue(cfg.MaxJobCount),
		jobChan:     make(chan Job, cfg.MaxJobCount),
		quitChan:    make(chan bool),
		redis:       cfg.Redis,
		maxJobCount: cfg.MaxJobCount,
		beforeJob:   cfg.BeforeJob,
		afterJob:    cfg.AfterJob,
	}
}

func (w *JobWorker) GetName() string {
	return w.Name
}

func (w *JobWorker) Run() {
	if w.isRunning {
		log.Printf("%s worker is running", w.Name)
		return
	}

	go func() {
		log.Printf("Start Worker(%s):", w.Name)
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
						err = w.queue.Enqueue(job)
						if err != nil {
							log.Println(err)
							continue
						}
						continue
					}
				}

				work(w, &job)
				saveJob(r, key, job)
			case <-w.quitChan:
				log.Printf("worker %s stopping\n", w.Name)
				return
			}
		}
	}()
}

func (w *JobWorker) Stop() {
	if !w.isRunning {
		log.Printf("%s worker is not running", w.Name)
		return
	}

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

func (w *JobWorker) IsRunning() bool {
	return w.isRunning
}

func (w *JobWorker) MaxJobCount() int {
	return w.maxJobCount
}

func (w *JobWorker) JobCount() int {
	return w.queue.Count()
}

func (w *JobWorker) BeforeJob(fn func(j *Job)) {
	w.beforeJob = fn
}

func (w *JobWorker) AfterJob(fn func(j *Job, err error)) {
	w.afterJob = fn
}

package worker

import (
	"github.com/go-redis/redis/v9"
	"log"
	"strconv"
)

type Dispatcher interface {
	Dispatch(jobId string, closure func(j *Job) error) error
	Run(names ...string)
	Stop(names ...string)
	SelectWorker(name string) Dispatcher
	GetWorkers() []Worker
	GetRedis() func() *redis.Client
	AddWorker(option Option)
	RemoveWorker(nam string)
	Status(isConsole bool) *Status
}

type JobDispatcher struct {
	workers []Worker
	worker  Worker
	Redis   func() *redis.Client
}

type Option struct {
	Name        string
	MaxJobCount int
	BeforeJob   func(j *Job)
	AfterJob    func(j *Job, err error)
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
		workers = append(workers, NewWorker(Config{
			o.Name,
			opt.Redis,
			o.MaxJobCount,
			o.BeforeJob,
			o.AfterJob,
		}))
	}

	return &JobDispatcher{
		workers: workers,
		worker:  nil,
		Redis:   opt.Redis,
	}
}

func (d *JobDispatcher) AddWorker(option Option) {
	d.workers = append(d.workers, NewWorker(Config{
		option.Name,
		d.Redis,
		option.MaxJobCount,
		option.BeforeJob,
		option.AfterJob,
	}))
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

func (d *JobDispatcher) BeforeJob(fn func(j *Job)) {
	if d.worker == nil {
		d.SelectWorker(DefaultWorker)
	}

	d.worker.BeforeJob(fn)
}

func (d *JobDispatcher) AfterJob(fn func(j *Job, err error)) {
	if d.worker == nil {
		d.SelectWorker(DefaultWorker)
	}

	d.worker.AfterJob(fn)
}

func (d *JobDispatcher) Dispatch(jobId string, closure func(j *Job) error) error {
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

func (d *JobDispatcher) Run(workerNames ...string) {
	workers := make([]Worker, 0)

	if len(workerNames) == 0 {
		workers = d.workers
	} else {
		for _, w := range d.workers {
			for _, wn := range workerNames {
				if wn == w.GetName() {
					workers = append(workers, w)
				}
			}
		}
	}

	for _, w := range workers {
		w.Run()
	}
}

func (d *JobDispatcher) Stop(workerNames ...string) {
	workers := make([]Worker, 0)

	if len(workerNames) == 0 {
		workers = d.workers
	} else {
		for _, w := range d.workers {
			for _, wn := range workerNames {
				if wn == w.GetName() {
					workers = append(workers, w)
				}
			}
		}
	}

	for _, w := range workers {
		w.Stop()
	}
}

type Status struct {
	Workers     []map[string]string `json:"workers"`
	WorkerCount int                 `json:"worker_count"`
}

func (d *JobDispatcher) Status(isConsole bool) *Status {

	workers := make([]map[string]string, 0)
	for _, w := range d.workers {
		workerInfo := map[string]string{
			"name":          w.GetName(),
			"is_running":    strconv.FormatBool(w.IsRunning()),
			"job_count":     strconv.Itoa(w.JobCount()),
			"max_job_count": strconv.Itoa(w.MaxJobCount()),
		}

		workers = append(workers, workerInfo)
	}

	if isConsole {
		for _, w := range workers {
			log.Printf("[worker name]: %s", w["name"])
			log.Printf("[worker is running]: %s", w["is_running"])
			log.Printf("[worker's job count]: %s", w["job_count"])
			log.Printf("[worker's max job count]:  %s", w["max_job_count"])
		}
	}

	return &Status{
		Workers:     workers,
		WorkerCount: len(workers),
	}
}

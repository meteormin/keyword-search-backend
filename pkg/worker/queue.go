package worker

type Job struct {
	JobId   string
	Status  bool
	Closure func()
}

type Queue interface {
	Enqueue(job Job)
	Dequeue() Job
	Clear()
	Exist(jobId string)
}

type JobQueue struct {
	queue []Job
}

func NewQueue() Queue {
	return &JobQueue{
		queue: make([]Job, 0),
	}
}

func (j *JobQueue) Enqueue(job Job) {
	j.queue = append(j.queue, job)
}

func (j *JobQueue) Dequeue() Job {
	return j.queue[len(j.queue)-1]
}

func (j *JobQueue) Clear() {
	j.queue = make([]Job, 0)
}

func (j *JobQueue) Exist(jobId string) {

}

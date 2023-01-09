package worker

type Job struct {
	Chan    chan error
	JobId   string
	Closure func()
}

type Queue interface {
	Enqueue(job Job)
	Dequeue()
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

func (j JobQueue) Enqueue(job Job) {
	//TODO implement me
	panic("implement me")
}

func (j JobQueue) Dequeue() {
	//TODO implement me
	panic("implement me")
}

func (j JobQueue) Clear() {
	//TODO implement me
	panic("implement me")
}

func (j JobQueue) Exist(jobId string) {
	//TODO implement me
	panic("implement me")
}

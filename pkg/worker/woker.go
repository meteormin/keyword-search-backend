package worker

type Job struct {
	Chan    chan error
	JobId   string
	Closure func()
}

type Worker interface {
	Enqueue(job Job)
	Dequeue()
	Clear()
	Exist(jobId string)
}

type JobQueueWorker struct {
	queue []Job
}

func NewWorker() Worker {
	return &JobQueueWorker{
		queue: make([]Job, 0),
	}
}

func (j JobQueueWorker) Enqueue(job Job) {
	//TODO implement me
	panic("implement me")
}

func (j JobQueueWorker) Dequeue() {
	//TODO implement me
	panic("implement me")
}

func (j JobQueueWorker) Clear() {
	//TODO implement me
	panic("implement me")
}

func (j JobQueueWorker) Exist(jobId string) {
	//TODO implement me
	panic("implement me")
}

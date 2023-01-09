package worker

type Worker interface {
	StartWorker()
	EndWorker()
	AddJob()
	RemoveJob()
}

type JobWorker struct {
	queue Queue
}

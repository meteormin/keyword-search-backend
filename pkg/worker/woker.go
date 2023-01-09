package worker

type Worker interface {
	StartWorker()
	EndWorker()
	AddJob()
	RemoveJob()
}

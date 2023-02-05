package api_jobs

import "github.com/miniyus/keyword-search-backend/pkg/worker"

type GetJobs struct {
	Jobs []worker.Job `json:"api_jobs"`
}

type GetJob struct {
	worker.Job
}

type GetStatus struct {
	worker.StatusInfo
}

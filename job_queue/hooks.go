package job_queue

import (
	"errors"
	"fmt"
	"github.com/miniyus/keyword-search-backend/database"
	"github.com/miniyus/keyword-search-backend/entity"
	"github.com/miniyus/keyword-search-backend/pkg/worker"
	"gorm.io/gorm"
	"log"
)

type WriteableField string

const (
	UserId WriteableField = "UserId"
)

func convJobHistory(job *worker.Job, err error) entity.JobHistory {
	jh := entity.JobHistory{}
	var errMsg *string
	if err != nil {
		errMessage := err.Error()
		errMsg = &errMessage
	}

	jh.JobId = job.JobId
	jh.Status = job.Status
	jh.UUID = job.UUID
	jh.WorkerName = job.WorkerName
	jh.Error = errMsg

	return jh
}

func parseMeta(jh entity.JobHistory, m map[string]interface{}) entity.JobHistory {
	if len(m) != 0 {
		for key, val := range m {
			switch key {
			case string(UserId):
				if v, ok := val.(uint); ok {
					jh.UserId = &v
				}
				break
			default:
				break
			}
		}
	}

	return jh
}

func createJobHistory(db *gorm.DB, j *worker.Job) error {
	log.Print("create job history")
	if j == nil {
		return errors.New("job is nil")
	}

	if db == nil {
		db = database.GetDB()
	}

	jh := convJobHistory(j, nil)
	jh = parseMeta(jh, j.Meta)
	tx := db.Create(&jh)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func updateJobHistory(db *gorm.DB, j *worker.Job, err error) error {
	if j == nil {
		return errors.New("job is nil")
	}

	if db == nil {
		db = database.GetDB()
	}

	jh := convJobHistory(j, err)

	var find entity.JobHistory

	tx := db.Where("uuid = ?", jh.UUID).First(&find)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("can not find job(%s)", jh.JobId))
	}

	if find.JobId == jh.JobId {
		jh.ID = find.ID
		jh.CreatedAt = find.CreatedAt
		jh.UserId = find.UserId
		db.Save(&jh)
	}

	return errors.New(fmt.Sprintf("can not find job(%s) from db", jh.JobId))
}

func RecordHistory(dispatcher worker.Dispatcher, db *gorm.DB) {
	jobMeta := make(map[string]interface{})

	dispatcher.OnDispatch(func(j *worker.Job) error {
		j.Meta = jobMeta
		return createJobHistory(db, j)
	})

	dispatcher.BeforeJob(func(j *worker.Job) error {
		return updateJobHistory(db, j, nil)
	})

	dispatcher.AfterJob(func(j *worker.Job, err error) error {
		return updateJobHistory(db, j, err)
	})

}

func AddMetaOnDispatch(dispatcher worker.Dispatcher, db *gorm.DB, meta map[WriteableField]interface{}) {
	jobMeta := make(map[string]interface{})

	for key, val := range meta {
		jobMeta[string(key)] = val
	}

	dispatcher.OnDispatch(func(j *worker.Job) error {
		j.Meta = jobMeta
		return createJobHistory(db, j)
	})
}

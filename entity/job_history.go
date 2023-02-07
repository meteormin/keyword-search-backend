package entity

import (
	"github.com/google/uuid"
	"github.com/miniyus/gofiber/pkg/worker"
	"gorm.io/gorm"
)

type JobHistory struct {
	gorm.Model
	UserId     *uint            `json:"user_id"`
	UUID       uuid.UUID        `json:"uuid" gorm:"column:uuid;type:varchar(100);uniqueIndex"`
	WorkerName string           `json:"worker_name" gorm:"column:worker_name;type:varchar(50);"`
	JobId      string           `json:"job_id" gorm:"column:job_id;type:varchar(50)"`
	Status     worker.JobStatus `json:"status" gorm:"column:status;type:varchar(10)"`
	Error      *string          `json:"error" gorm:"column:error;type:varchar(255)"`
	User       User             `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

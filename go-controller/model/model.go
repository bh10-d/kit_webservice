package model

import (
	"gorm.io/gorm"
	"github.com/lib/pq"
	"github.com/google/uuid"
	"time"
)

type Runner struct {
	ID   string `gorm:"primaryKey" json:"id"`
	HostName string `json:"hostname"`
	IP   string `json:"ip"`
	Tags string `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    // DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // nếu muốn soft delete
}

type Job struct {
	ID             uint   `gorm:"primaryKey"`
	RunnerID       string
	MsgID          string  // Stores the base job ID (same for all runners in one request)
	Status         string
	RequestPayload string
	ResponsePayload string
	Timeout        bool
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    // DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // nếu muốn soft delete
}

type Scripts struct {
	ScriptID   string `gorm:"primaryKey" json:"script_id"`
	FileName   string `json:"file_name"`
	Description string `json:"description"`
	Param      pq.StringArray  `gorm:"type:text[]" json:"param"`
	Status	 bool `json:"status"`
	Tag 	  pq.StringArray  `gorm:"type:text[]" json:"tag"` // Mảng các tag
	Runner        pq.StringArray  `gorm:"type:text[]" json:"runner"`  // Mảng các VM
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    // DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // nếu muốn soft delete
}

type Logs struct {
	MsgID    string `gorm:"primaryKey" json:"msg_id"`
	RunnerID string `gorm:"primaryKey" json:"runner_id"`
	Logs     string `json:"logs"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    // DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // nếu muốn soft delete
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Runner{}, &Job{}, &Logs{}, &Scripts{})
}

func (s *Scripts) BeforeCreate(tx *gorm.DB) (err error) {
    if s.ScriptID == "" {
        s.ScriptID = uuid.New().String()
    }
    return
}

// GetJobsByBaseID returns all jobs with the same base job ID
func GetJobsByBaseID(db *gorm.DB, baseJobID string) ([]Job, error) {
    var jobs []Job
    err := db.Where("msg_id = ?", baseJobID).Find(&jobs).Error
    return jobs, err
}

// GetJobGroupSummary returns summary of job execution for a base job ID
func GetJobGroupSummary(db *gorm.DB, baseJobID string) (map[string]interface{}, error) {
    var jobs []Job
    err := db.Where("msg_id = ?", baseJobID).Find(&jobs).Error
    if err != nil {
        return nil, err
    }
    
    total := len(jobs)
    success := 0
    failed := 0
    timeout := 0
    
    for _, job := range jobs {
        switch job.Status {
        case "success", "200":
            success++
        case "timeout":
            timeout++
        default:
            failed++
        }
    }
    
    return map[string]interface{}{
        "base_job_id": baseJobID,
        "total_runners": total,
        "success_count": success,
        "failed_count": failed,
        "timeout_count": timeout,
        "jobs": jobs,
    }, nil
}

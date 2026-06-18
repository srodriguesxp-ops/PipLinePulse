package entities

import "time"

type PipelineStatus string
type JobStatus string

const (
	StatusPendig   PipelineStatus = "pending"
	StatusRunning  PipelineStatus = "running"
	StatusSuccess  PipelineStatus = "success"
	StatusFailed   PipelineStatus = "failed"
	StatusCanceled PipelineStatus = "canceled"
)

const (
	JobPending  JobStatus = "pending"
	JobRunning  JobStatus = "running"
	JobSuccess  JobStatus = "success"
	JobFailed   JobStatus = "failed"
	JobCanceled JobStatus = "canceled"
)

type Pipelinerun struct {
	ID         int64          `json:"id"`
	RunID      int64          `json:"run_id"`     //Id Workflow github
	Repository string         `json:"repository"` //puxar repositorio do workflow
	Branch     string         `json:"branch"`
	CommitSHA  string         `json:"commit_sha"`
	Status     PipelineStatus `json:"status"`
	StartedAt  time.Time      `json:"started_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Jobs       []Job          `json:"jobs"`
}

type Job struct {
	ID            int64      `json:"id"`
	PipelinerunID int64      `json:"pipelinerun_id"` //FK
	Name          string     `json:"name"`
	Status        JobStatus  `json:"status"`
	StartedAt     time.Time  `json:"started_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at"` //Pode ser nulo se o job ainda estiver em execução
}

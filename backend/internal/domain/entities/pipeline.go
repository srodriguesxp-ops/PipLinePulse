package entities

import "time"

type PipelineStatus string
type JobStatus string

const (
	StatusPending   PipelineStatus = "pending"
	StatusRunning   PipelineStatus = "running"
	StatusSuccess   PipelineStatus = "success"
	StatusFailed    PipelineStatus = "failed"
	StatusCancelled PipelineStatus = "cancelled"
)

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobSuccess   JobStatus = "success"
	JobFailed    JobStatus = "failed"
	JobCancelled JobStatus = "cancelled"
)

type PipelineRun struct {
	ID         int64          `json:"id"`
	RunID      int64          `json:"run_id"`
	Repository string         `json:"repository"`
	Branch     string         `json:"branch"`
	CommitSHA  string         `json:"commit_sha"`
	Status     PipelineStatus `json:"status"`
	StartedAt  time.Time      `json:"started_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	Jobs       []Job          `json:"jobs"`
}

type Job struct {
	ID            int64      `json:"id"`
	PipelineRunID int64      `json:"pipelinerun_id"`
	Name          string     `json:"name"`
	Status        JobStatus  `json:"status"`
	StartedAt     time.Time  `json:"started_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

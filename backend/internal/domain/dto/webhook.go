package dto

import "time"

type WorkflowRunEvent struct {
	Action      string      `json:"action"`
	WorkflowRun WorkflowRun `json:"workflow_run"`
	Repository  Repository  `json:"repository"`
	Sender      Sender      `json:"sender"`
}

type WorkflowRun struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	HeadBranch string    `json:"head_branch"`
	HeadSHA    string    `json:"head_sha"`
	Status     string    `json:"status"`
	Conclusion *string   `json:"conclusion"`
	HTMLURL    string    `json:"html_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RunAttempt int       `json:"run_attempt"`
	Jobs       []Job     `json:"jobs,omitempty"`
}

type Job struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  *string    `json:"conclusion"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Steps       []Step     `json:"steps"`
	HTMLURL     string     `json:"html_url"`
}

type Step struct {
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  *string    `json:"conclusion"`
	Number      int        `json:"number"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	Private  bool   `json:"private"`
}

type Sender struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

package github

import (
	"github.com/srodriguesxp-ops/PipLinePulse/backend/internal/domain/dto"
	"github.com/srodriguesxp-ops/PipLinePulse/backend/internal/domain/entities"
)

// ParseWorkflowRun converte o payload bruto recebido do webhook do GitHub Actions
func ParseWorkflowRun(event dto.WorkflowRunEvent) entities.PipelineRun {
	run := event.WorkflowRun

	pipeline := entities.PipelineRun{
		RunID:      run.ID,
		Repository: event.Repository.FullName,
		Branch:     run.HeadBranch,
		CommitSHA:  run.HeadSHA,
		Status:     mapPipelineStatus(run.Status, run.Conclusion),
		StartedAt:  run.CreatedAt,
		UpdatedAt:  run.UpdatedAt,
		Jobs:       parseJobs(run.ID, run.Jobs),
	}

	return pipeline
}

// mapPipelineStatus traduz o par (status, conclusion) do GitHub Actions
func mapPipelineStatus(status string, conclusion *string) entities.PipelineStatus {
	switch status {
	case "queued":
		return entities.StatusPending
	case "in_progress":
		return entities.StatusRunning
	case "completed":
		if conclusion == nil {
			// Caso defensivo: GitHub não deveria mandar "completed" sem conclusion,
			return entities.StatusFailed
		}
		switch *conclusion {
		case "success":
			return entities.StatusSuccess
		case "cancelled":
			return entities.StatusCancelled
		default:
			// "failure", "timed_out", "action_required", etc. caem aqui.
			return entities.StatusFailed
		}
	default:
		return entities.StatusPending
	}
}

// mapJobStatus traduz o status individual de um job dentro do workflow run.
func mapJobStatus(status string, conclusion *string) entities.JobStatus {
	switch status {
	case "queued":
		return entities.JobPending
	case "in_progress":
		return entities.JobRunning
	case "completed":
		if conclusion == nil {
			return entities.JobFailed
		}
		switch *conclusion {
		case "success":
			return entities.JobSuccess
		case "cancelled":
			return entities.JobCancelled
		default:
			return entities.JobFailed
		}
	default:
		return entities.JobPending
	}
}

// parseJobs converte a lista de jobs do DTO na lista de jobs da entity,
// associando cada um ao PipelineRun via pipelineRunID.
func parseJobs(pipelineRunID int64, dtoJobs []dto.Job) []entities.Job {
	jobs := make([]entities.Job, 0, len(dtoJobs))

	for _, j := range dtoJobs {
		job := entities.Job{
			PipelineRunID: pipelineRunID,
			Name:          j.Name,
			Status:        mapJobStatus(j.Status, j.Conclusion),
			CompletedAt:   j.CompletedAt,
		}

		if j.StartedAt != nil {
			job.StartedAt = *j.StartedAt
		}

		jobs = append(jobs, job)
	}

	return jobs
}

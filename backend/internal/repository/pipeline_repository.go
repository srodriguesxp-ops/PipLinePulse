package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/srodriguesxp-ops/PipLinePulse/backend/internal/domain/entities"
)

// PipelineRepository encapsula todas as operações de banco de dados
// relacionadas a pipeline_runs e jobs.
type PipelineRepository struct {
	db *pgxpool.Pool
}

// NewPipelineRepository cria uma nova instância do repositório
// recebendo o pool de conexões por injeção de dependência.
func NewPipelineRepository(db *pgxpool.Pool) *PipelineRepository {
	return &PipelineRepository{db: db}
}

// Save persiste um PipelineRun e seus jobs no banco de dados dentro de uma transaction.
func (r *PipelineRepository) Save(ctx context.Context, pipeline entities.PipelineRun) (entities.PipelineRun, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entities.PipelineRun{}, fmt.Errorf("repository.Save: begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		INSERT INTO pipeline_runs (run_id, repository, branch, commit_sha, status, started_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		pipeline.RunID,
		pipeline.Repository,
		pipeline.Branch,
		pipeline.CommitSHA,
		pipeline.Status,
		pipeline.StartedAt,
		pipeline.UpdatedAt,
	).Scan(&pipeline.ID)
	if err != nil {
		return entities.PipelineRun{}, fmt.Errorf("repository.Save: insert pipeline_run: %w", err)
	}

	for i, job := range pipeline.Jobs {
		err = tx.QueryRow(ctx, `
			INSERT INTO jobs (pipeline_run_id, name, status, started_at, updated_at, completed_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`,
			pipeline.ID,
			job.Name,
			job.Status,
			job.StartedAt,
			job.UpdatedAt,
			job.CompletedAt,
		).Scan(&pipeline.Jobs[i].ID)
		if err != nil {
			return entities.PipelineRun{}, fmt.Errorf("repository.Save: insert job %q: %w", job.Name, err)
		}
		pipeline.Jobs[i].PipelineRunID = pipeline.ID
	}

	if err = tx.Commit(ctx); err != nil {
		return entities.PipelineRun{}, fmt.Errorf("repository.Save: commit transaction: %w", err)
	}

	return pipeline, nil
}

// FindByRunID busca um pipeline_run pelo ID de execução do GitHub Actions.
func (r *PipelineRepository) FindByRunID(ctx context.Context, runID int64) (entities.PipelineRun, error) {
	var pipeline entities.PipelineRun

	err := r.db.QueryRow(ctx, `
		SELECT id, run_id, repository, branch, commit_sha, status, started_at, updated_at
		FROM pipeline_runs
		WHERE run_id = $1`,
		runID,
	).Scan(
		&pipeline.ID,
		&pipeline.RunID,
		&pipeline.Repository,
		&pipeline.Branch,
		&pipeline.CommitSHA,
		&pipeline.Status,
		&pipeline.StartedAt,
		&pipeline.UpdatedAt,
	)
	if err != nil {
		return entities.PipelineRun{}, fmt.Errorf("repository.FindByRunID: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, pipeline_run_id, name, status, started_at, updated_at, completed_at
		FROM jobs
		WHERE pipeline_run_id = $1
		ORDER BY id ASC`,
		pipeline.ID,
	)
	if err != nil {
		return entities.PipelineRun{}, fmt.Errorf("repository.FindByRunID: query jobs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var job entities.Job
		err = rows.Scan(
			&job.ID,
			&job.PipelineRunID,
			&job.Name,
			&job.Status,
			&job.StartedAt,
			&job.UpdatedAt,
			&job.CompletedAt,
		)
		if err != nil {
			return entities.PipelineRun{}, fmt.Errorf("repository.FindByRunID: scan job: %w", err)
		}
		pipeline.Jobs = append(pipeline.Jobs, job)
	}

	return pipeline, nil
}

// UpdateStatus atualiza o status de um pipeline_run existente.
func (r *PipelineRepository) UpdateStatus(ctx context.Context, runID int64, status entities.PipelineStatus) error {
	_, err := r.db.Exec(ctx, `
		UPDATE pipeline_runs
		SET status = $1, updated_at = NOW()
		WHERE run_id = $2`,
		status,
		runID,
	)
	if err != nil {
		return fmt.Errorf("repository.UpdateStatus: %w", err)
	}
	return nil
}

// ListByRepository retorna os pipelines mais recentes de um repositório.
func (r *PipelineRepository) ListByRepository(ctx context.Context, repository string) ([]entities.PipelineRun, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, run_id, repository, branch, commit_sha, status, started_at, updated_at
		FROM pipeline_runs
		WHERE repository = $1
		ORDER BY started_at DESC
		LIMIT 50`,
		repository,
	)
	if err != nil {
		return nil, fmt.Errorf("repository.ListByRepository: %w", err)
	}
	defer rows.Close()

	var pipelines []entities.PipelineRun
	for rows.Next() {
		var p entities.PipelineRun
		err = rows.Scan(
			&p.ID,
			&p.RunID,
			&p.Repository,
			&p.Branch,
			&p.CommitSHA,
			&p.Status,
			&p.StartedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository.ListByRepository: scan: %w", err)
		}
		pipelines = append(pipelines, p)
	}

	return pipelines, nil
}

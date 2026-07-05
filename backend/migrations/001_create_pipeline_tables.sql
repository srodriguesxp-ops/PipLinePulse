--Tabela Pipeline_runs

CREATE TABLE IF NOT EXISTS pipeline_runs (
    id          BIGSERIAL    PRIMARY KEY,
    run_id      BIGINT       NOT NULL     UNIQUE,
    repository  VARCHAR(255) NOT NULL,
    branch      VARCHAR(255) NOT NULL,
    commit_sha  CHAR(40)     NOT NULL,
    status      VARCHAR(20)  NOT NULL     DEFAULT 'pending',
                CHECK (status in ('pending', 'running', 'success', 'failed', 'canceled')),
    created_at  TIMESTAMP    NOT NULL     DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL     DEFAULT NOW()
)


--Tabela Jobs

CREATE TABLE IF NOT EXISTS jobs (

    id              BIGSERIAL PRIMARY KEY,
    pipeline_run_id BIGINT NOT NULL REFERENCES pipeline_runs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
                    CHECK (status in ('pending', 'running', 'success', 'failed', 'canceled')),
    started_at      TIMESTAMP    NOT NULL     DEFAULT NOW(),
    updated_at      TIMESTAMP    NOT NULL     DEFAULT NOW()
    completed_at    TIMESTAMP    NOT NULL     --Null enquanto o job estiver em execução, preenchido quando o job for concluído
)


--Indices

CREATE INDEX IF NOT EXISTS idx_jobs_pipelin_run_id 
    ON jobs(pipeline_run_id);

CREATE INDEX IF NOT EXISTS idx_pipeline_runs_repository
    ON pipeline_runs (repository);

CREATE INDEX IF NOT EXISTS idx_pipeline_runs_started_at
    ON pipeline_runs (started_at DESC);
CREATE TABLE IF NOT EXISTS schema_revisions
(
    version        VARCHAR(255) NOT NULL,
    description    VARCHAR(255) NOT NULL DEFAULT '',
    applied        BIGINT       NOT NULL DEFAULT 0,
    total          BIGINT       NOT NULL DEFAULT 0,
    executed_at    TIMESTAMP    NULL,
    execution_time BIGINT       NOT NULL DEFAULT 0,
    error          LONGTEXT     NULL,
    hash           VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (version)
) CHARSET utf8mb4
  COLLATE utf8mb4_bin

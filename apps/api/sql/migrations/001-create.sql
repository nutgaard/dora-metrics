DO
$$
    BEGIN
        CREATE DOMAIN ksuid as VARCHAR(27);
        CREATE DOMAIN url as varchar CHECK (value ~ '^https?:\/\/');
        CREATE DOMAIN environment as varchar(15) CHECK (value ~ '^..+$');
        CREATE DOMAIN status as varchar(15) CHECK (value ~ '^open|closed|merged|cancelled$');
        CREATE DOMAIN event_type as varchar(15) CHECK (value ~ '^pullrequest|incident$');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

CREATE TABLE IF NOT EXISTS deployment
(
    id             ksuid PRIMARY KEY NOT NULL,
    created_at     TIMESTAMP         NOT NULL DEFAULT NOW(),
    started_at     TIMESTAMP         NOT NULL,
    finished_at    TIMESTAMP         NOT NULL,
    repository_url url               NOT NULL,
    environment    environment       NOT NULL,
    metadata       JSONB             NOT NULL DEFAULT '{}'::jsonb
);


CREATE TABLE IF NOT EXISTS event
(
    id                   ksuid PRIMARY KEY NOT NULL,
    event_type           event_type        NOT NULL,
    created_at           TIMESTAMP         NOT NULL DEFAULT NOW(),
    status               status            NOT NULL,
    repository_url       url               NOT NULL,
    environment          environment       NOT NULL,
    metadata             JSONB             NOT NULL DEFAULT '{}'::jsonb,
    opened_at            TIMESTAMP,
    closed_at            TIMESTAMP,

    deployment_reference ksuid,
    CONSTRAINT fk_event_reference FOREIGN KEY (deployment_reference) REFERENCES deployment (id)
);

CREATE INDEX IF NOT EXISTS event_metadata_idx on event USING gin (metadata jsonb_ops);
CREATE INDEX IF NOT EXISTS deployment_metadata_idx on deployment USING gin (metadata jsonb_ops);
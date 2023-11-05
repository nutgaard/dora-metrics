DO
$$
    BEGIN
        CREATE DOMAIN ksuid as VARCHAR(27);
        CREATE DOMAIN url as varchar CHECK (value ~ '^https?:\/\/');
        CREATE DOMAIN environment as varchar(5) CHECK (value ~ '^(prod|stage|test|dev|ldev)$');
    EXCEPTION
        WHEN duplicate_object THEN null;
    END
$$;

CREATE TABLE IF NOT EXISTS deployment
(
    id             ksuid PRIMARY KEY       NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW() NOT NULL,
    started_at     TIMESTAMP               NOT NULL,
    finished_at    TIMESTAMP               NOT NULL,

    repository_url url                     NOT NULL,
    application    VARCHAR                 NOT NULL,
    environment    environment             NOT NULL,

    department     VARCHAR,
    team           VARCHAR,
    product        VARCHAR,
    version        VARCHAR
);
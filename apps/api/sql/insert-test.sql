insert into deployment (id, started_at, finished_at, repository_url, application, environment, department, team,
                        product, version)
values ('random-ksuid',
        NOW(),
        NOW() + interval '10 minute',
        'https://github.com/code-obos/boligjakt-frontend',
        'my application',
        'stage',
        null, null, null, null);
# dora-metrics-reporter

Report events in repo to dora-metrics

## Inputs

### API-Key

**Required** API-Key for your organization/user

## Outputs

## Setup

```yaml
on:
  - commit
  - pr

jobs:
  report-to-dora-metrics:
    runs-on: ubuntu-latest
    name: Report to dora-metrics
    steps:
      - name: Report
        uses: nutgaard/dora-metrics/apps/reporter@v1
        env:
          api-key: ${secrets.dora-api-key}
```
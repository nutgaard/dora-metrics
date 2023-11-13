# DORA Metrics

Application to capture and visualize DevOps Research and Assessment (DORA) metrics.

## The four principles

- **Deployment frequency** — How often an organization successfully releases to production
- **Lead time for changes** — The amount of time it takes a commit to get into production
- **Change failure rate** — The percentage of deployments causing a failure in production
- **Time to restore service** — How long it takes an organization to recover from a failure in production


## Captured data

**Pull requests:**
Pull-requests are captured at creation, and updated when they are part of a deployment.

**Incidents:**
Incidents are captured at creation, and updated when they are part of a deployment.

**Deployments:**
A deployment should consist of one or more PRs which should be delivered to production.
It can optionally also contain references to incidents, which should be solved within the given deployment.



## How it is measured

**Deployment frequence:** `average(deployments.groupBy('date').count())`

**Lead time for changes:** `average(deployments.map(timeBetweenFirstCommitAndDeploy))`

**Change failure rate:** `deployments.filter(failures).length / deployments.length`

**Time to restore service:** `average(incidents.map(timeBetweenResolutionAndCreation))`
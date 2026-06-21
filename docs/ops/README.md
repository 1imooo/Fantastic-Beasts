# Ops

Operator procedures for **fantastic-beasts**.

| Script | Purpose |
|--------|---------|
| `./cmd/apply-dashboards` | Push Grafana/Kibana dashboards to mono |
| `./cmd/apply-alerts` | Push alert rules to mono |

## Sitemap

See [sitemap.md](sitemap.md). The Go server serves `/sitemap.xml` dynamically — no build step required.

## Deploy

Release pipeline: `release/production.sh` (live) / `release/test.sh` (test). CI runs `release/steps/*.sh` — build, test, scan, deploy.

Host ports and domain wiring: [project.yml](../project.yml).

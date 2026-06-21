# Telemetry

Observability for **fantastic-beasts** on mono.

| Signal | Doc | Source | Destination |
|--------|-----|--------|-------------|
| **Container logs** | [logging.md](logging.md) | Go server → stdout | Filebeat → Elasticsearch → Kibana |
| **Container health** | — | Docker / Watchtower | mono host |

This is a Go HTTP server — there are no custom application metrics. Dashboards focus on container log volume and access patterns from ingress or proxy logs where available.

## Platform artifacts

| Kind | Repo path |
|------|-----------|
| Dashboards | [dashboards/](../../dashboards/) |
| Alerts | [alerts/](../alerts/) |

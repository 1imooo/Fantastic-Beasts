# Alerts

Alert rule definitions for **fantastic-beasts**. Source files in [alerts/](../../alerts/). Applied via `./cmd/apply-alerts`.

## Suggested rules

| Rule | Condition | When |
|------|-----------|------|
| Site down | No Apache access logs for 15m | Container or reverse proxy issue |
| 5xx spike | HTTP 5xx rate elevated in access logs | Application or upstream errors |
| Container restarts | Frequent restarts on mono | Deploy or runtime instability |

Configure notification channels on mono (Slack, email, PagerDuty) in Grafana/Kibana UI.

See [logging.md](../telemetry/logging.md) for log fields to match in queries.

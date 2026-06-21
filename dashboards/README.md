# Dashboards (source of truth)

UI exports for **fantastic-beasts**. Not deployed by app CI — import on **mono** separately.

| File | Platform |
|------|----------|
| `grafana.json` | Grafana / Mimir |
| `kibana.ndjson` | Kibana / Elasticsearch |

## Apply (from laptop)

```bash
./cmd/apply-dashboards --on-mono
./cmd/apply-alerts --on-mono
```

On the mono host:

```bash
./cmd/apply-dashboards
./cmd/apply-alerts
```

Credentials (`GRAFANA_ADMIN_*`, optional `ELASTIC_*`) load from `.env.test` / `.env.production` via `./cmd/apply-env-file-hashicorp`.

| Choice | `DASHBOARD_ENV` / `ALERTS_ENV` |
|--------|--------------------------------|
| test (default) | `test` |
| live | `production` |

## Where to look after deploy

| Signal | Where |
|--------|--------|
| Apache access / errors | `docker logs <container>` or Kibana Discover on **`logstash-*`**, filter `container.name: *fantastic-beasts*` |
| Saved searches | **Kibana** `fantastic-beasts-access` dashboard |
| Request overview | **Grafana** dashboard (`fantastic-beasts-test` / `fantastic-beasts-production`) |

Data view name: **`fantastic-beasts-access`** (Apache access logs).

## Grafana

Dashboard UID/title include the environment (`fantastic-beasts-test` / `fantastic-beasts-production`). Datasource UID `mimir`.

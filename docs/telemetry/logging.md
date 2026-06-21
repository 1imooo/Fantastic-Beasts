# Logging

The Go server writes errors and startup messages to **stdout**. Container logs are the primary telemetry for **fantastic-beasts**.

## Log fields

Filter in Kibana on `logstash-*` or the container name:

| Field | Use |
|-------|-----|
| `container.name` | `fantastic-beasts-live` / `fantastic-beasts-test` |
| `message` | Log line (errors, startup) |
| `@timestamp` | Event time |

If a reverse proxy sits in front of the container, access logs may appear in ingress or proxy logs rather than the app container.

## Useful KQL filters

| Goal | KQL |
|------|-----|
| All traffic | `container.name: *fantastic-beasts*` |
| Errors | `message: *error*` or `message: *Internal Server Error*` |
| Beast detail pages | `message: */beasts/*` |
| Search usage | `message: *search=*` |
| Legacy PHP redirects | `message: *beast-details.php*` |

## Grafana

HTTP request rates can be derived from log-based metrics if configured on mono. The bundled Grafana dashboard uses container log volume as a proxy when Prometheus metrics are unavailable.

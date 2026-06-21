# Fantastic Beasts

Go web application for browsing and searching magical creatures from the Fantastic Beasts universe.

## What it does

1. **Browse** beasts with classification badges and descriptions.
2. **Search** using semantic TF-IDF matching over names and descriptions.
3. **Sort** alphabetically or by Ministry of Magic classification (1–5).
4. **Detail pages** per beast with full descriptions.

## Layout

```
src/              Go server, templates, search
assets/           beasts.json (beast data)
dashboards/       Grafana + Kibana exports
alerts/           Alert rule definitions
docs/             Project docs and deploy manifest
```

## Quick start

```bash
make build
make run
# http://localhost:8080
```

## Ops

| Task | Doc |
|------|-----|
| Sitemap (dynamic) | [ops/sitemap.md](ops/sitemap.md) |
| Operator scripts | [ops/README.md](ops/README.md) |

## Telemetry

[telemetry/README.md](telemetry/README.md)

## Deploy

[project.yml](project.yml)

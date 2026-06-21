# Fantastic Beasts

A Go web application for exploring magical creatures from the Fantastic Beasts universe.

## Features

- Browse beasts with Ministry classifications and descriptions
- Semantic search (TF-IDF + cosine similarity)
- Sort alphabetically or by classification (1–5)
- Individual beast detail pages
- Dynamic robots.txt and sitemap.xml
- Docker deployment

## Quick start

```bash
make build
make run
# http://localhost:8080
```

Local development without Docker:

```bash
make run-local
# http://localhost:8080
```

## Project structure

```
Fantastic-Beasts/
├── src/                   Go server, templates (nottinfra stack)
│   ├── main.go
│   ├── search.go
│   └── templates/
├── assets/                beasts.json
├── dashboards/            Grafana + Kibana exports
├── alerts/                Alert rule definitions
├── docs/                  Project docs and deploy manifest
├── Makefile
└── Dockerfile
```

## Development

### Adding beasts

Edit `assets/beasts.json`:

```json
{
  "name": "Beast Name",
  "classification": 3,
  "description": "Detailed description..."
}
```

### Sitemap and robots.txt

Both are served dynamically at `/sitemap.xml` and `/robots.txt` from `src/templates/`. No build step is required after editing beasts.

## Author

**1imo**

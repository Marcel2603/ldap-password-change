# Configuration

The service is configured in layers, applied in this priority order (highest wins):

1. `cmd/config/app.default.yml` – bundled defaults
2. `app.yml` – user-supplied overrides (optional)
3. Environment variables – highest priority, useful for containers

See the sub-pages for the full reference of each configuration section.

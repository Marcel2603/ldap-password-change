# Customisation

The service supports full visual customisation without rebuilding the binary. All assets are resolved
at runtime based on the `ui` configuration block.

## Configuration

```yaml
ui:
  backgroundImage: ""   # Replaces the animated gradient background
  customCss: ""         # Injected as an additional stylesheet (after Bootstrap)
  favicon: ""           # Replaces the browser tab icon
  icon: ""              # Replaces the logo image above the form
```

## Asset Resolution

Values are resolved in this order of precedence:

| Input value              | Resolved to                           |
|--------------------------|---------------------------------------|
| `""`                     | Bundled default asset                 |
| `"logo.png"`             | `/custom/logo.png` (auto-prefixed)    |
| `"custom/logo.png"`      | `/custom/logo.png` (slash prepended)  |
| `"/custom/logo.png"`     | `/custom/logo.png` (used as-is)       |
| `"https://example.com/logo.png"` | External URL (used as-is)     |

This means you can simply write `logo.png` in your config — no need to type the full `/custom/` prefix.

## Serving Local Files

The service exposes a `/custom/*` static file endpoint automatically backed by the `./custom/` directory
on the host filesystem.

### Docker Compose Example

```yaml
# docker-compose.yml
services:
  ldap-password-change:
    image: ghcr.io/marcel2603/github.com/Marcel2603/ldap-password-change/ldap-password-change:latest
    ports:
      - "3000:3000"
    volumes:
      - ./app.yml:/app/app.yml
      - ./custom:/app/custom   # <-- mount your assets here
```

```yaml
# app.yml
ui:
  backgroundImage: "bg.jpg"           # served from ./custom/bg.jpg
  customCss: "corporate-theme.css"    # served from ./custom/corporate-theme.css
  favicon: "favicon.ico"              # served from ./custom/favicon.ico
  icon: "logo.png"                    # served from ./custom/logo.png
```

## Theme Selector

The built-in theme switcher in the top-right corner lets users pick:

- ☀️ **Light** – forced light mode
- 🌙 **Dark** – forced dark mode
- 💻 **System** – follows the OS `prefers-color-scheme` media query

The choice is stored in `localStorage` and persisted across page reloads.

Your custom CSS is loaded **after** the built-in styles, so you can override any CSS custom property:

```css
/* custom/corporate-theme.css */
:root {
  --app-bg: linear-gradient(135deg, #003366, #0066cc);
}

[data-bs-theme="dark"] {
  --glass-bg: rgba(0, 20, 50, 0.6);
}
```

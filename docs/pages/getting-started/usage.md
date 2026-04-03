# Usage

## Changing a Password

1. Open the service URL in your browser (default: `http://localhost:3000`).
2. Enter your **Username** — must match the configured `validation.username` pattern.
3. Enter your **Current Password**.
4. Enter and confirm your **New Password** — must satisfy the `validation.password` policy.
5. Click **Change Password**.

The service will search the LDAP directory for your user, verify your current credentials via a bind,
and then apply the new password using an LDAP Modify operation.

## Form Validation

All validation runs client-side (HTML5 + Alpine.js) before the request is submitted:

| Field              | Rule                                            |
|--------------------|--------------------------------------------------|
| Username           | Must match `validation.username` regex pattern  |
| Current Password   | Must not be empty                               |
| New Password       | Must match `validation.password` regex + differ from current |
| Confirm Password   | Must match New Password exactly                 |

Server-side validation applies the same regex patterns as an additional security layer.

## Theme Switching

The UI supports three appearance modes, switchable via the dropdown in the top-right corner:

| Mode   | Behaviour                                          |
|--------|----------------------------------------------------|
| ☀️ Light  | Always light theme                               |
| 🌙 Dark   | Always dark theme                                |
| 💻 System | Follows the OS `prefers-color-scheme` setting    |

The selected mode is persisted in `localStorage` across page reloads.

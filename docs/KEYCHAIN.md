# Keychain storage (macOS, Windows, Linux)

By default, the CLI tries to store your token in the OS keychain and falls back to a local file if a keychain isn't available.

- macOS: uses the system Keychain. Works out of the box on a normal desktop session.
- Windows: uses Windows Credential Manager. Works out of the box.
- Linux: uses the Secret Service API. Works on desktop sessions running a keyring daemon (e.g., GNOME Keyring, KWallet). On headless servers/WSL, a keychain may not be available—CLI will fall back to `~/.projet-iac/token.json` with 0700/0600 perms.

## Force behavior

- `--keychain on` or `KEYCHAIN=on`: try keychain; if unavailable, fall back to file.
- `--keychain off` or `KEYCHAIN=off`: disable keychain; always use file.
- `--keychain auto` or `KEYCHAIN=auto` (default): use keychain if available, else file.

## Linux setup tips

If you want keychain on Linux without a full desktop:

- Install a keyring and dbus:
  - Debian/Ubuntu: `sudo apt install gnome-keyring dbus`
  - Fedora: `sudo dnf install gnome-keyring dbus`
- Ensure a DBus session and unlocked keyring are available. In non-graphical environments this can be tricky; for servers, the secure file fallback is acceptable for many school/demo scenarios.

## Security notes

- Keychain backends encrypt secrets at rest and integrate with OS policies (screen lock, login, etc.).
- File fallback uses restrictive permissions (0700 dir, 0600 file). Keep your account protected and disk encrypted to protect secrets at rest.

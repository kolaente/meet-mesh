---
# meet-mesh-us08
title: Add storage config and avatar directory initialization
status: completed
type: task
priority: normal
created_at: 2026-02-11T21:00:00Z
updated_at: 2026-02-11T19:27:17Z
parent: meet-mesh-us07
---

# Add Storage Config and Avatar Directory Initialization

**Goal:** Add a `storage` section to the config with an `avatars_path` field, and ensure the directory is created on startup.

**Architecture:** Extend the existing `Config` struct in `api/config.go` with a new `StorageConfig` struct. Add directory creation logic in `api/cmd/main.go` at startup, right after database init. Default path is `./data/avatars/`.

---

## Files

- Modify: `api/config.go`
- Modify: `api/cmd/main.go`
- Modify: `config.example.yaml`

---

## Step 1: Add StorageConfig to config.go

Open `api/config.go`. Add a new struct after the existing `SMTPConfig`:

```go
type StorageConfig struct {
	AvatarsPath string `yaml:"avatars_path"`
}
```

Add it as a field in the `Config` struct:

```go
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	OIDC     OIDCConfig     `yaml:"oidc"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	Storage  StorageConfig  `yaml:"storage"`
}
```

Add a `SetDefaults` method to the `Config` struct:

```go
func (c *Config) SetDefaults() {
	if c.Storage.AvatarsPath == "" {
		c.Storage.AvatarsPath = "./data/avatars"
	}
}
```

At the end of the `LoadConfig` function, before `return &cfg, nil`, add:

```go
cfg.SetDefaults()
```

---

## Step 2: Create avatar directory on startup in main.go

Open `api/cmd/main.go`. Add `"os"` to the imports.

After the database initialization block (after `db, err := api.InitDatabase(...)` and its error check), add:

```go
// Ensure avatar storage directory exists
if err := os.MkdirAll(cfg.Storage.AvatarsPath, 0755); err != nil {
    log.Fatalf("Failed to create avatar storage directory: %v", err)
}
```

---

## Step 3: Update config.example.yaml

Add the storage section at the end of the file:

```yaml
storage:
  avatars_path: ./data/avatars
```

---

## Step 4: Verify

Run:

```bash
cd api && go build ./cmd
```

Expected: Compiles without errors.

---

## Step 5: Commit

```bash
git add api/config.go api/cmd/main.go config.example.yaml
git commit -m "feat(api): add storage config with avatars_path and directory initialization"
```

## Summary of Changes

- Added StorageConfig struct with AvatarsPath field to api/config.go
- Added Storage field to Config struct
- Added SetDefaults method to Config with default avatars_path "./data/avatars"
- Added os.MkdirAll call in main.go to create avatar directory on startup
- Added storage section to config.example.yaml

Build passes with go build ./cmd

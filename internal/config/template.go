package config

const DefaultConfigTemplate = `# dsync config (starter template)
#
# Source of truth documentation:
# - dsync_go_documentation/spec/CONFIG.md
# - dsync_go_documentation/spec/CLI.md
# - dsync_go_documentation/spec/BEHAVIORS.md

[global]
excludes = [
  ".DS_Store",
  ".git/",
  "node_modules/",
  ".dsync-partial/",
]

# Endpoints are named directory roots.
# Note: dsync uses rsync "contents semantics" (trailing /): roots represent directory contents.

[endpoints.example_local]
type = "local"
path = "/Users/you/photos"

[endpoints.example_remote]
type = "ssh"
host = "photo-box" # ssh-config Host alias
path = "/srv/photos"

# Links are 1:1 mappings between one local and one remote endpoint.
[links.photos]
local = "example_local"
remote = "example_remote"
mirror = true
partial_only = false
paths = []
excludes = ["*.tmp"]
`

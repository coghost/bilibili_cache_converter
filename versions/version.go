package versions

// Version is automatically set via: go build -ldflags="-X 'bilibili/versions.Version=1.0.0'"
// No need to manually update this.
// Note: main.Version will not work — use versions.Version instead.
var Version = ""

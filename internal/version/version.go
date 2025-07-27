// Package version provides version information for DictCLI.
package version

import (
	"fmt"
	"runtime"
)

// Build information populated via ldflags during compilation.
var (
	Version   = "dev"      // Version number (e.g., "1.0.0")
	GitCommit = "unknown"  // Git commit hash
	BuildTime = "unknown"  // Build timestamp
	GoVersion = runtime.Version() // Go version used for compilation
)

// Info returns formatted version information.
func Info() string {
	return fmt.Sprintf("DictCLI %s", Version)
}

// FullInfo returns detailed version information including build details.
func FullInfo() string {
	return fmt.Sprintf(`DictCLI %s
Git Commit: %s
Build Time: %s
Go Version: %s
Platform:   %s/%s`,
		Version,
		GitCommit,
		BuildTime,
		GoVersion,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
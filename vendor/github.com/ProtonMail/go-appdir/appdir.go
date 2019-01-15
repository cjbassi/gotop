// Get application directories such as config and cache.
package appdir

// Dirs requests application directories paths.
type Dirs interface {
	// Get the user-specific config directory.
	UserConfig() string
	// Get the user-specific cache directory.
	UserCache() string
	// Get the user-specific logs directory.
	UserLogs() string
}

// New creates a new App with the provided name.
func New(name string) Dirs {
	return &dirs{name: name}
}

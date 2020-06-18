package widgets

import (
	"fmt"
	"strings"
)

// makeName creates a prometheus metric name in the gotop space
// This function doesn't have to be very efficient because it's only
// called at init time, and only a few dozen times... and it isn't
// (very efficient).
func makeName(parts ...interface{}) string {
	args := make([]string, len(parts)+1)
	args[0] = "gotop"
	for i, v := range parts {
		args[i+1] = fmt.Sprintf("%v", v)
	}
	return strings.Join(args, "_")
}

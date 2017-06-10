// Package goqueryaws provides aws sdk wrapper functions that use tags
// to query AWS resources.
package goqueryaws

import (
	"os"
)

// GetProfile is a hack to try to get the profile of the current session,
// since I can't find a way to do that using aws sdk.
func GetProfile() string {
	if os.Getenv("AWS_PROFILE") == "" {
		// return default when a) AWS_PROFILE is not set, or b) AWS_PROFILE is
		// set to empty string
		return "default"
	} else {
		return os.Getenv("AWS_PROFILE")
	}
}

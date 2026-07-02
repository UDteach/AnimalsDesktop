//go:build windows && !animalsdesktop_nonetwork

package main

import "testing"

func TestDefaultBuildEnablesNetworkUpdates(t *testing.T) {
	if !networkUpdatesEnabled {
		t.Fatalf("default Windows build should keep update network access enabled")
	}
}

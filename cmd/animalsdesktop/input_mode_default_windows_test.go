//go:build windows && !animalsdesktop_nonetwork

package main

import "testing"

func TestDefaultWindowsBuildKeepsInputMonitoring(t *testing.T) {
	if !inputMonitoringEnabled {
		t.Fatalf("default Windows build should keep input monitoring")
	}
}

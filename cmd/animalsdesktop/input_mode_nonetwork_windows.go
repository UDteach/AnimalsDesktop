//go:build windows && animalsdesktop_nonetwork

package main

const inputMonitoringEnabled = false

func (a *petApp) installInputMonitoring() {}

func (a *petApp) cleanupInputMonitoring() {}

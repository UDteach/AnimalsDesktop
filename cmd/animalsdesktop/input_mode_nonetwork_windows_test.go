//go:build windows && animalsdesktop_nonetwork

package main

import "testing"

func TestNoNetworkBuildDisablesInputMonitoring(t *testing.T) {
	if inputMonitoringEnabled {
		t.Fatalf("security-check build should disable input monitoring")
	}
	if got := normalizeBehaviorMode(int(modeKeyboard)); got != modeRandom {
		t.Fatalf("normalizeBehaviorMode(modeKeyboard) = %d, want modeRandom", got)
	}
}

func TestNoNetworkBuildIgnoresInputReactionPaths(t *testing.T) {
	a := &petApp{
		mode:         modeKeyboard,
		wheelEnabled: true,
		pets: []desktopPet{{
			state:      stateIdle,
			stateTicks: 12,
		}},
	}

	a.onTyping()
	if got := a.pets[0].state; got != stateIdle {
		t.Fatalf("onTyping() changed pet state to %d in security-check build", got)
	}
	a.onMouseClick(0, 0)
	if got := len(a.reactions); got != 0 {
		t.Fatalf("onMouseClick() created %d reactions in security-check build", got)
	}
	if a.handleMenuCommand(menuModeKeyboard) {
		t.Fatalf("security-check build should reject the keyboard mode command")
	}
	if a.handleMenuCommand(menuWheelToggle) {
		t.Fatalf("security-check build should reject the typing-wheel command")
	}
}

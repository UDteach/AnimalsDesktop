package main

import (
	"os"
	"strings"
	"testing"
)

func TestSecurityEditionSourceSeparatesGlobalInputHooks(t *testing.T) {
	defaultSource := readSourceFile(t, "input_mode_default_windows.go")
	for _, want := range []string{
		"//go:build windows && !animalsdesktop_nonetwork",
		"SetWindowsHookExW",
		"UnhookWindowsHookEx",
		"CallNextHookEx",
	} {
		if !strings.Contains(defaultSource, want) {
			t.Fatalf("default input source does not contain %q", want)
		}
	}

	securitySource := readSourceFile(t, "input_mode_nonetwork_windows.go")
	if !strings.Contains(securitySource, "//go:build windows && animalsdesktop_nonetwork") {
		t.Fatalf("security-check input source is missing its no-network build tag")
	}
	for _, forbidden := range []string{
		"SetWindowsHookExW",
		"UnhookWindowsHookEx",
		"CallNextHookEx",
		"GetCursorPos",
	} {
		if strings.Contains(securitySource, forbidden) {
			t.Fatalf("security-check input source contains %q", forbidden)
		}
	}

	mainSource := readSourceFile(t, "main_windows.go")
	for _, forbidden := range []string{
		`NewProc("SetWindowsHookExW")`,
		`NewProc("UnhookWindowsHookEx")`,
		`NewProc("CallNextHookEx")`,
	} {
		if strings.Contains(mainSource, forbidden) {
			t.Fatalf("shared Windows source contains default-only hook binding %q", forbidden)
		}
	}
}

func readSourceFile(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(data)
}

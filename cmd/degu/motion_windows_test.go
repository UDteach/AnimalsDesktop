//go:build windows

package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHorizontalMotionFramesUseFullWalkSequence(t *testing.T) {
	states := []behaviorState{
		stateWalk,
		stateForage,
		stateCarry,
	}

	for _, state := range states {
		seen := map[int]bool{}
		for frame := 0; frame < walkFrames*2; frame++ {
			got := currentFrame(state, frame)
			if got < walkStart || got >= walkStart+walkFrames {
				t.Fatalf("currentFrame(%v, %d) = %d, want full walk block frame", state, frame, got)
			}
			seen[got] = true
		}
		if len(seen) != walkFrames {
			t.Fatalf("currentFrame(%v) used %d walk frames, want %d", state, len(seen), walkFrames)
		}
	}
}

func TestScurryUsesDedicatedFastMotionFrames(t *testing.T) {
	for frame := 0; frame < scurryFrames*2; frame++ {
		got := currentFrame(stateScurry, frame)
		if got < scurryStart || got >= scurryStart+scurryFrames {
			t.Fatalf("currentFrame(stateScurry, %d) = %d, want dedicated scurry frame", frame, got)
		}
	}
	if got := currentFrame(stateScurry, scurryFrames); got != scurryStart {
		t.Fatalf("scurry loop frame = %d, want %d", got, scurryStart)
	}
}

func TestWheelUsesDedicatedRunFrames(t *testing.T) {
	for frame := 0; frame < wheelRunFrames*2; frame++ {
		got := currentFrame(stateWheel, frame)
		if got < wheelRunStart || got >= wheelRunStart+wheelRunFrames {
			t.Fatalf("currentFrame(stateWheel, %d) = %d, want dedicated wheelrun frame", frame, got)
		}
	}
	if got := currentFrame(stateWheel, wheelRunFrames); got != wheelRunStart {
		t.Fatalf("wheelrun loop frame = %d, want %d", got, wheelRunStart)
	}
}

func TestFrameFromSeqHandlesEmptyAndBadDivisor(t *testing.T) {
	if got := frameFromSeq(nil, 12, 2); got != idleStart {
		t.Fatalf("frameFromSeq(nil) = %d, want %d", got, idleStart)
	}
	seq := []int{7, 9}
	if got := frameFromSeq(seq, 3, 0); got != 9 {
		t.Fatalf("frameFromSeq with zero divisor = %d, want 9", got)
	}
}

func TestFrameFromSeqClampedHoldsFinalFrame(t *testing.T) {
	seq := []int{7, 9, 11}
	if got := frameFromSeqClamped(seq, 999, 2); got != 11 {
		t.Fatalf("frameFromSeqClamped past end = %d, want 11", got)
	}
	if got := frameFromSeqClamped(seq, 3, 0); got != 11 {
		t.Fatalf("frameFromSeqClamped with zero divisor = %d, want 11", got)
	}
}

func TestTypingStartsAndExtendsWheelOnlyInKeyboardMode(t *testing.T) {
	a := &petApp{
		mode:         modeKeyboard,
		wheelEnabled: true,
		wheelX:       400,
		sceneW:       1200,
		speed:        3,
		pets: []deguPet{
			{state: stateWalk, stateTicks: 12, item: noItem},
			{state: stateWalk, stateTicks: 12, item: noItem},
		},
	}

	a.onTyping()
	if got := a.pets[0].state; got != stateWheel {
		t.Fatalf("first pet state = %v, want stateWheel", got)
	}
	if got := a.pets[0].stateTicks; got != wheelKeyHold {
		t.Fatalf("wheel hold ticks = %d, want %d", got, wheelKeyHold)
	}
	if got := a.pets[0].moveSpeed; got != 0 {
		t.Fatalf("wheel pet moveSpeed = %d, want 0", got)
	}
	wantX := clamp(a.wheelX-wheelSize/2, 0, max(0, a.sceneW-spriteW))
	if got := a.pets[0].x; got != wantX {
		t.Fatalf("wheel pet x = %d, want %d", got, wantX)
	}
	if got := a.pets[1].state; got != stateScurry {
		t.Fatalf("second pet state = %v, want stateScurry", got)
	}

	a.pets[0].frame = 7
	a.pets[0].stateTicks = 3
	a.onTyping()
	if got := a.pets[0].frame; got != 7 {
		t.Fatalf("wheel frame reset while extending: got %d, want 7", got)
	}
	if got := a.pets[0].stateTicks; got != wheelKeyHold {
		t.Fatalf("extended wheel hold ticks = %d, want %d", got, wheelKeyHold)
	}
}

func TestTypingDoesNotStartWheelInRandomMode(t *testing.T) {
	a := &petApp{
		mode:         modeRandom,
		wheelEnabled: true,
		wheelX:       400,
		sceneW:       1200,
		pets: []deguPet{
			{state: stateWalk, stateTicks: 12, item: noItem},
		},
	}

	a.onTyping()
	if got := a.pets[0].state; got == stateWheel {
		t.Fatalf("typing in random mode started wheel state")
	}
}

func TestForageItemsStayHidden(t *testing.T) {
	a := &petApp{
		sceneW: 1200,
		speed:  3,
		forage: []forageItem{
			{x: 100, kind: 0, owner: noItem, active: true},
			{x: 200, kind: 1, owner: reservedItem, active: true},
		},
	}

	a.ensureForageItems()
	for i, item := range a.forage {
		if item.active || item.owner != noItem {
			t.Fatalf("forage item %d = %+v, want inactive and unowned", i, item)
		}
	}

	p := deguPet{state: stateWalk, item: noItem, dir: 1}
	if a.maybeAssignForageTarget(&p) {
		t.Fatalf("maybeAssignForageTarget returned true while forage is hidden")
	}
}

func TestRuntimeCatalogIsReleaseScopedToChinchilla(t *testing.T) {
	if got := len(variants); got != 1 {
		t.Fatalf("runtime variants = %d, want 1", got)
	}
	if got := variants[0].ID; got != "chinchilla_standard_gray" {
		t.Fatalf("runtime variant = %q, want chinchilla_standard_gray", got)
	}
	for _, variant := range variants {
		if variant.SpeciesID == "degu" {
			t.Fatalf("runtime variants include degu: %+v", variant)
		}
	}
}

func TestSettingsRoundTripPersistsCoreOptions(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)

	a := &petApp{
		variant:       4,
		coatMode:      coatSelected,
		selectedCoats: [maxPetCount]int{1, 3, 5, 7, 9, 0, 2, 4, 6, 8},
		petNames:      [maxPetCount]string{"モカ", "Sora", "  Nagi  ", "", "", "", "", "", "", ""},
		nameLabels:    true,
		speed:         5,
		mode:          modeKeyboard,
		petCount:      10,
		wheelEnabled:  false,
		bidirectional: false,
		lang:          langEnglish,
		settingsX:     220,
		settingsY:     180,
	}
	if err := a.saveSettings(); err != nil {
		t.Fatalf("saveSettings() error = %v", err)
	}

	path := filepath.Join(configRoot, settingsDirName, settingsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("settings file was not written: %v", err)
	}
	var saved appSettings
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("settings json is invalid: %v", err)
	}
	if saved.Version != 1 || saved.PetCount != 10 || saved.Mode != int(modeKeyboard) {
		t.Fatalf("saved settings = %+v, want version 1 petCount 10 keyboard mode", saved)
	}
	if !saved.NameLabels {
		t.Fatalf("saved NameLabels = false, want true")
	}
	if got := saved.PetNames[0]; got != "モカ" {
		t.Fatalf("saved pet name 0 = %q, want モカ", got)
	}
	if got := saved.PetNames[2]; got != "Nagi" {
		t.Fatalf("saved pet name 2 = %q, want sanitized Nagi", got)
	}

	b := &petApp{
		variant:       0,
		coatMode:      coatRandom,
		selectedCoats: [maxPetCount]int{0, 1, 2, 4, 8, 6, 3, 7, 5, 9},
		speed:         3,
		mode:          modeRandom,
		petCount:      2,
		wheelEnabled:  true,
		bidirectional: true,
		lang:          langJapanese,
		settingsX:     120,
		settingsY:     120,
	}
	if err := b.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if b.variant != clamp(a.variant, 0, len(variants)-1) || b.coatMode != a.coatMode || b.speed != a.speed || b.mode != a.mode || b.petCount != a.petCount {
		t.Fatalf("loaded scalar settings = variant:%d coat:%d speed:%d mode:%d count:%d", b.variant, b.coatMode, b.speed, b.mode, b.petCount)
	}
	if b.wheelEnabled != a.wheelEnabled || b.bidirectional != a.bidirectional || b.lang != a.lang {
		t.Fatalf("loaded flags = wheel:%v bidirectional:%v lang:%d", b.wheelEnabled, b.bidirectional, b.lang)
	}
	if b.nameLabels != a.nameLabels {
		t.Fatalf("loaded nameLabels = %v, want %v", b.nameLabels, a.nameLabels)
	}
	for i := 0; i < maxPetCount; i++ {
		want := clamp(a.selectedCoats[i], 0, len(variants)-1)
		if b.selectedCoats[i] != want {
			t.Fatalf("selectedCoats[%d] = %d, want %d", i, b.selectedCoats[i], want)
		}
	}
	if b.petNames[0] != "モカ" || b.petNames[1] != "Sora" || b.petNames[2] != "Nagi" {
		t.Fatalf("loaded pet names = %#v", b.petNames[:3])
	}
}

func TestPetVariantRectsFitTenPetsInSettingsWindow(t *testing.T) {
	seen := map[[4]int]bool{}
	for i := 0; i < maxPetCount; i++ {
		numberRect, buttonRect := settingsPetVariantRects(i)
		if buttonRect.Right > 708 || buttonRect.Bottom > 502 {
			t.Fatalf("pet variant button %d rect %+v overflows selected-coats panel", i, buttonRect)
		}
		if numberRect.Left < 238 || buttonRect.Left <= numberRect.Right {
			t.Fatalf("pet variant %d number/button rects overlap or escape: number=%+v button=%+v", i, numberRect, buttonRect)
		}
		key := [4]int{int(buttonRect.Left), int(buttonRect.Top), int(buttonRect.Right), int(buttonRect.Bottom)}
		if seen[key] {
			t.Fatalf("pet variant button %d duplicates another rect: %+v", i, buttonRect)
		}
		seen[key] = true
	}
}

func TestPetNameRectsFitTenPetsWithCoatPicker(t *testing.T) {
	for i := 0; i < maxPetCount; i++ {
		numberRect, nameRect := settingsPetNameRects(i)
		_, coatRect := settingsPetVariantRects(i)
		if nameRect.Right >= coatRect.Left {
			t.Fatalf("pet %d name rect overlaps coat rect: name=%+v coat=%+v", i, nameRect, coatRect)
		}
		if numberRect.Left < 238 || nameRect.Left <= numberRect.Right || coatRect.Right > 708 || nameRect.Bottom > 502 {
			t.Fatalf("pet %d name/coat row escapes panel: number=%+v name=%+v coat=%+v", i, numberRect, nameRect, coatRect)
		}
	}
}

func TestUpdateVersionComparison(t *testing.T) {
	tests := []struct {
		latest  string
		current string
		want    bool
	}{
		{"v1.2.0", "v1.1.9", true},
		{"v1.2.0", "1.2.0", false},
		{"v1.2.0", "v1.3.0", false},
		{"v2.0.0", "dev", true},
		{"v2.0.0", "pages-abc123", true},
		{"not-semver", "v1.0.0", false},
	}
	for _, tt := range tests {
		if got := isNewerVersion(tt.latest, tt.current); got != tt.want {
			t.Fatalf("isNewerVersion(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
		}
	}
}

func TestSelectUpdateAssetFindsWindowsZip(t *testing.T) {
	rel := &githubRelease{Assets: []githubReleaseAsset{
		{Name: "notes.txt", BrowserDownloadURL: "https://example.test/notes.txt"},
		{Name: "AnimalsDesktop-windows-amd64.zip", BrowserDownloadURL: "https://example.test/app.zip"},
		{Name: "AnimalsDesktop-windows-386.zip", BrowserDownloadURL: "https://example.test/app-x86.zip"},
	}}
	asset := selectUpdateAsset(rel, "amd64")
	if asset == nil || asset.BrowserDownloadURL != "https://example.test/app.zip" {
		t.Fatalf("selectUpdateAsset(amd64) = %+v", asset)
	}
	asset = selectUpdateAsset(rel, "386")
	if asset == nil || asset.BrowserDownloadURL != "https://example.test/app-x86.zip" {
		t.Fatalf("selectUpdateAsset(386) = %+v", asset)
	}
}

func TestVerifyDownloadedAssetChecksSizeAndSHA256Digest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "update.zip")
	data := []byte("trusted update bytes")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write update: %v", err)
	}
	sum := sha256.Sum256(data)
	asset := githubReleaseAsset{
		Size:   int64(len(data)),
		Digest: fmt.Sprintf("sha256:%x", sum[:]),
	}
	if err := verifyDownloadedAsset(path, asset); err != nil {
		t.Fatalf("verifyDownloadedAsset() error = %v", err)
	}
	asset.Size++
	if err := verifyDownloadedAsset(path, asset); err == nil {
		t.Fatalf("verifyDownloadedAsset accepted a size mismatch")
	}
	asset.Size = int64(len(data))
	asset.Digest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	if err := verifyDownloadedAsset(path, asset); err == nil {
		t.Fatalf("verifyDownloadedAsset accepted a digest mismatch")
	}
}

func TestParseUpdateApplyArgsRequiresSafeCleanupDir(t *testing.T) {
	cleanupDir := filepath.Join(os.TempDir(), updateTempPrefix+"unit-test")
	opts, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(cleanupDir, "AnimalsDesktop.exe"),
		"--target", filepath.Join(t.TempDir(), "AnimalsDesktop.exe"),
		"--parent-pid", "1234",
		"--cleanup-dir", cleanupDir,
	})
	if err != nil {
		t.Fatalf("parseUpdateApplyArgs() error = %v", err)
	}
	if opts.ParentPID != 1234 || opts.CleanupDir != cleanupDir {
		t.Fatalf("parseUpdateApplyArgs() = %+v", opts)
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", "a.exe",
		"--target", "b.exe",
		"--cleanup-dir", t.TempDir(),
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a non-update cleanup dir")
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(t.TempDir(), "AnimalsDesktop.exe"),
		"--target", filepath.Join(t.TempDir(), "AnimalsDesktop.exe"),
		"--cleanup-dir", cleanupDir,
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a source outside cleanup dir")
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(cleanupDir, "payload", "AnimalsDesktop.exe"),
		"--target", filepath.Join(cleanupDir, "installed", "AnimalsDesktop.exe"),
		"--cleanup-dir", cleanupDir,
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a target inside cleanup dir")
	}
	if _, err := parseUpdateApplyArgs([]string{
		"--source", filepath.Join(cleanupDir, "payload", "AnimalsDesktop.exe"),
		"--target", filepath.Join(t.TempDir(), "renamed.exe"),
		"--cleanup-dir", cleanupDir,
	}); err == nil {
		t.Fatalf("parseUpdateApplyArgs accepted a renamed target")
	}
}

func TestExtractUpdateExeUsesFixedPayloadPath(t *testing.T) {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "update.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}
	zw := zip.NewWriter(zipFile)
	w, err := zw.Create("release/nested/AnimalsDesktop.exe")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write([]byte("exe bytes")); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := zipFile.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	exePath, err := extractUpdateExe(zipPath, tmpDir)
	if err != nil {
		t.Fatalf("extractUpdateExe() error = %v", err)
	}
	want := filepath.Join(tmpDir, "payload", "AnimalsDesktop.exe")
	if exePath != want {
		t.Fatalf("extractUpdateExe() = %q, want %q", exePath, want)
	}
	data, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatalf("read extracted exe: %v", err)
	}
	if string(data) != "exe bytes" {
		t.Fatalf("extracted data = %q", data)
	}
}

func TestUpdaterCommandsInvokeAppExeDirectly(t *testing.T) {
	cleanupDir := filepath.Join(os.TempDir(), updateTempPrefix+"command-test")
	sourceExe := filepath.Join(cleanupDir, "payload", "AnimalsDesktop.exe")
	targetExe := filepath.Join(t.TempDir(), "AnimalsDesktop.exe")
	helperExe := filepath.Join(cleanupDir, "helper", "AnimalsDesktop.exe")

	applyCmd := newUpdaterHelperCommand(helperExe, cleanupDir, sourceExe, targetExe, 1234)
	assertCommandAvoidsPowerShell(t, applyCmd)
	if !strings.EqualFold(applyCmd.Path, helperExe) {
		t.Fatalf("apply command path = %q, want %q", applyCmd.Path, helperExe)
	}
	assertArgsContainInOrder(t, applyCmd.Args,
		updaterApplyArg,
		"--source", sourceExe,
		"--target", targetExe,
		"--parent-pid", "1234",
		"--cleanup-dir", cleanupDir,
	)
	if applyCmd.SysProcAttr == nil || !applyCmd.SysProcAttr.HideWindow {
		t.Fatalf("apply command should hide its helper window")
	}

	cleanupCmd := newUpdaterCleanupCommand(targetExe, cleanupDir)
	assertCommandAvoidsPowerShell(t, cleanupCmd)
	if !strings.EqualFold(cleanupCmd.Path, targetExe) {
		t.Fatalf("cleanup command path = %q, want %q", cleanupCmd.Path, targetExe)
	}
	assertArgsContainInOrder(t, cleanupCmd.Args, updaterCleanupArg, cleanupDir)
	if cleanupCmd.SysProcAttr == nil || !cleanupCmd.SysProcAttr.HideWindow {
		t.Fatalf("cleanup command should hide its helper window")
	}
}

func TestReleaseWorkflowPowerShellBlocksParseAfterGitHubSubstitution(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	scriptByStep := map[string]string{
		"Generate assets":       extractWorkflowRunBlock(t, workflowPath, "Generate assets"),
		"Build":                 extractWorkflowRunBlock(t, workflowPath, "Build"),
		"Package":               extractWorkflowRunBlock(t, workflowPath, "Package"),
		"Prepare release notes": extractWorkflowRunBlock(t, workflowPath, "Prepare release notes"),
	}
	for stepName, script := range scriptByStep {
		t.Run(stepName, func(t *testing.T) {
			script = strings.ReplaceAll(script, "${{ github.ref_name }}", "v0.2.0")
			assertPowerShellParses(t, script)
		})
	}
}

func TestReleaseWorkflowPackageIncludesSecurityManifestAndHashes(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	script := extractWorkflowRunBlock(t, workflowPath, "Package")
	for _, want := range []string{
		"SECURITY.txt",
		"SHA256SUMS.txt",
		"AnimalsDesktop-windows-amd64.zip/AnimalsDesktop.exe",
		"AnimalsDesktop-windows-386.zip/AnimalsDesktop.exe",
		"Microsoft Security Intelligence",
		"McAfee Dispute Detection & Allowlisting",
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("Package script does not contain %q", want)
		}
	}
}

func TestReleaseWorkflowKeepsV020PreviewReleaseExplicit(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	workflow := string(data)
	for _, want := range []string{
		`$version -eq "v0.2.0"`,
		"go run ./cmd/validatemotion -runtime-only",
		"go run ./cmd/validatemotion -runtime-only -require-accepted",
		"body_path: dist/RELEASE_NOTES.md",
		"github.ref_name == 'v0.2.0'",
		"docs/releases/$version.md",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("release workflow does not contain %q", want)
		}
	}
}

func TestUpdateCleanupDirOnlyAcceptsUpdateTempDirs(t *testing.T) {
	cleanupDir := filepath.Join(os.TempDir(), updateTempPrefix+"cleanup-test")
	if got := updateCleanupDir([]string{updaterCleanupArg, cleanupDir}); got != cleanupDir {
		t.Fatalf("updateCleanupDir() = %q, want %q", got, cleanupDir)
	}
	if got := updateCleanupDir([]string{updaterCleanupArg, t.TempDir()}); got != "" {
		t.Fatalf("updateCleanupDir accepted non-update dir %q", got)
	}
}

func assertCommandAvoidsPowerShell(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
	joined := strings.ToLower(strings.Join(cmd.Args, " "))
	for _, forbidden := range []string{"powershell", "pwsh", ".ps1", "executionpolicy", "bypass"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("command unexpectedly contains %q: %q", forbidden, joined)
		}
	}
}

func assertArgsContainInOrder(t *testing.T, args []string, want ...string) {
	t.Helper()
	next := 0
	for _, arg := range args {
		if next < len(want) && arg == want[next] {
			next++
		}
	}
	if next != len(want) {
		t.Fatalf("args %q did not contain %q in order", args, want)
	}
}

func extractWorkflowRunBlock(t *testing.T, workflowPath, stepName string) string {
	t.Helper()
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	stepIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "- name: "+stepName {
			stepIndex = i
			break
		}
	}
	if stepIndex < 0 {
		t.Fatalf("workflow step %q was not found", stepName)
	}
	runIndex := -1
	for i := stepIndex + 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "- name: ") {
			break
		}
		if trimmed == "run: |" {
			runIndex = i
			break
		}
	}
	if runIndex < 0 {
		t.Fatalf("workflow step %q has no run block", stepName)
	}
	contentIndent := -1
	for i := runIndex + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		contentIndent = leadingSpaces(lines[i])
		break
	}
	if contentIndent < 0 {
		t.Fatalf("workflow step %q has an empty run block", stepName)
	}
	var out []string
	for i := runIndex + 1; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r")
		if strings.TrimSpace(line) != "" && leadingSpaces(line) < contentIndent {
			break
		}
		if len(line) >= contentIndent {
			line = line[contentIndent:]
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func leadingSpaces(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}

func assertPowerShellParses(t *testing.T, script string) {
	t.Helper()
	powershell, err := exec.LookPath("powershell.exe")
	if err != nil {
		powershell, err = exec.LookPath("powershell")
	}
	if err != nil {
		t.Skipf("PowerShell was not found: %v", err)
	}
	scriptPath := filepath.Join(t.TempDir(), "script-under-test.ps1")
	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		t.Fatalf("write script under test: %v", err)
	}
	parser := `
$script = Get-Content -LiteralPath $args[0] -Raw
$tokens = $null
$errors = $null
[System.Management.Automation.Language.Parser]::ParseInput($script, [ref]$tokens, [ref]$errors) | Out-Null
if ($errors.Count -gt 0) {
  $errors | ForEach-Object { $_.Message }
  exit 1
}
`
	cmd := exec.Command(powershell, "-NoProfile", "-NonInteractive", "-Command", parser, scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("PowerShell parser rejected script: %v\n%s", err, out)
	}
}

func TestTurnStateUsesGeneratedTurnFrames(t *testing.T) {
	if got := currentFrame(stateTurn, 0); got != turnStart {
		t.Fatalf("turn frame 0 = %d, want %d", got, turnStart)
	}
	if got := currentFrame(stateTurn, turnTicks-1); got != turnStart+turnFrames-1 {
		t.Fatalf("turn final active frame = %d, want %d", got, turnStart+turnFrames-1)
	}
	if got := currentFrame(stateTurn, turnTicks+10); got != turnStart+turnFrames-1 {
		t.Fatalf("turn frame after duration = %d, want held final frame %d", got, turnStart+turnFrames-1)
	}
}

func TestTurnDrawDirectionMirrorsOnlyLeftToRightTurns(t *testing.T) {
	if got := turnDrawDirection(1, -1); got != 1 {
		t.Fatalf("right-to-left turn draw direction = %d, want 1", got)
	}
	if got := turnDrawDirection(-1, 1); got != -1 {
		t.Fatalf("left-to-right turn draw direction = %d, want -1", got)
	}
}

func TestSetBidirectionalOffNormalizesPets(t *testing.T) {
	a := &petApp{
		bidirectional: true,
		speed:         3,
		pets: []deguPet{
			{state: stateTurn, dir: -1, nextDir: -1, item: noItem},
			{state: stateWalk, dir: -1, nextDir: -1, item: noItem},
		},
	}

	a.setBidirectional(false)
	if a.bidirectional {
		t.Fatalf("bidirectional stayed enabled")
	}
	for i, pet := range a.pets {
		if pet.dir != 1 || pet.nextDir != 1 {
			t.Fatalf("pet %d direction = (%d,%d), want (1,1)", i, pet.dir, pet.nextDir)
		}
		if pet.state == stateTurn {
			t.Fatalf("pet %d remained in stateTurn", i)
		}
	}
}

func TestFixedCoatModeRefreshesAllPets(t *testing.T) {
	a := &petApp{
		variant: 99,
		pets: []deguPet{
			{variant: 0},
			{variant: 0},
			{variant: 0},
		},
	}

	a.setCoatMode(coatFixed)

	want := len(variants) - 1
	for i, pet := range a.pets {
		if pet.variant != want {
			t.Fatalf("pet %d variant = %d, want fixed variant %d", i, pet.variant, want)
		}
	}
}

func TestSelectedCoatModeUsesPerPetChoices(t *testing.T) {
	a := &petApp{
		selectedCoats: [maxPetCount]int{0, 3, 5, 7, 9},
		pets: []deguPet{
			{variant: 0},
			{variant: 0},
			{variant: 0},
		},
	}

	a.setCoatMode(coatSelected)

	for i := range []int{0, 1, 2} {
		want := clamp(a.selectedCoats[i], 0, len(variants)-1)
		if got := a.pets[i].variant; got != want {
			t.Fatalf("pet %d variant = %d, want %d", i, got, want)
		}
	}
	a.setSelectedVariant(1, 8)
	want := len(variants) - 1
	if got := a.pets[1].variant; got != want {
		t.Fatalf("selected variant update = %d, want %d", got, want)
	}
}

func TestRandomCoatModeAssignsValidVariants(t *testing.T) {
	a := &petApp{coatMode: coatRandom}
	for i := 0; i < 100; i++ {
		got := a.variantForIndex(i)
		if got < 0 || got >= len(variants) {
			t.Fatalf("random variant = %d, want 0..%d", got, len(variants)-1)
		}
	}
}

func TestPetAtScenePointFindsTopmostPet(t *testing.T) {
	a := &petApp{
		sceneW: 800,
		pets: []deguPet{
			{x: 100, laneOffset: 0, state: stateWalk},
			{x: 110, laneOffset: 0, state: stateIdle},
		},
	}

	got := a.petAtScenePoint(132, sceneH-spriteH+24)
	if got != 1 {
		t.Fatalf("petAtScenePoint overlap = %d, want topmost pet 1", got)
	}
	if got := a.petAtScenePoint(4, 4); got != -1 {
		t.Fatalf("petAtScenePoint outside = %d, want -1", got)
	}
}

func TestShowPetReactionRefreshesExistingBubble(t *testing.T) {
	a := &petApp{
		pets: []deguPet{{state: stateWalk}},
		reactions: []petReaction{
			{pet: 0, kind: 1, ticks: 3},
		},
	}

	a.showPetReaction(0)
	if len(a.reactions) != 1 {
		t.Fatalf("reaction count = %d, want 1 refreshed reaction", len(a.reactions))
	}
	if a.reactions[0].ticks != reactionTicks {
		t.Fatalf("reaction ticks = %d, want %d", a.reactions[0].ticks, reactionTicks)
	}
}

func TestTickReactionsDropsExpiredAndInvalid(t *testing.T) {
	a := &petApp{
		pets: []deguPet{{state: stateWalk}},
		reactions: []petReaction{
			{pet: 0, ticks: 1},
			{pet: 3, ticks: 5},
		},
	}

	a.tickReactions()
	if len(a.reactions) != 0 {
		t.Fatalf("remaining reactions = %d, want 0", len(a.reactions))
	}
}

func TestSpriteCacheLoadsVariantOnDemand(t *testing.T) {
	cache := newSpriteCache()
	if got := len(cache.loaded); got != 0 {
		t.Fatalf("new cache loaded variants = %d, want 0", got)
	}
	frame := cache.frame(variants[0], 0, 0)
	if frame == nil || frame.Bounds().Dx() != frameW || frame.Bounds().Dy() != frameH {
		t.Fatalf("loaded frame bounds = %v", frame.Bounds())
	}
	if got := len(cache.loaded); got != 1 {
		t.Fatalf("cache loaded variants = %d, want 1", got)
	}
	_ = cache.frame(variants[0], motionSets+99, frameCount+99)
	if got := len(cache.loaded); got != 1 {
		t.Fatalf("cache reloaded same variant; loaded = %d", got)
	}
}

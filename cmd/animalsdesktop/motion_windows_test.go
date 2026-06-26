//go:build windows

package main

import (
	"encoding/json"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/lxn/win"
)

func TestHorizontalMotionFramesUseStableRightFacingWalkSequence(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	states := []behaviorState{
		stateWalk,
		stateScurry,
		stateWheel,
		stateForage,
		stateCarry,
	}

	for _, state := range states {
		for frame := 0; frame < 64; frame++ {
			got := currentFrame(state, frame)
			if !allowed[got] {
				t.Fatalf("currentFrame(%v, %d) = %d, want stable right-facing walk frame", state, frame, got)
			}
		}
	}
}

func TestScurryUsesStableFastWalkFrames(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	for frame := 0; frame < scurryFrames*2; frame++ {
		got := currentFrame(stateScurry, frame)
		if !allowed[got] {
			t.Fatalf("currentFrame(stateScurry, %d) = %d, want stable fast walk frame", frame, got)
		}
	}
}

func TestWheelUsesStableWalkFramesForAnimalRuntime(t *testing.T) {
	allowed := map[int]bool{
		walkStart:     true,
		walkStart + 1: true,
		walkStart + 3: true,
	}
	for frame := 0; frame < wheelRunFrames*2; frame++ {
		got := currentFrame(stateWheel, frame)
		if !allowed[got] {
			t.Fatalf("currentFrame(stateWheel, %d) = %d, want stable walk frame", frame, got)
		}
	}
}

func TestWeakNibbleVariantsUseStableActionFallback(t *testing.T) {
	for _, variant := range []coatVariant{
		{ID: "sugar_glider_gray"},
		{ID: "rabbit_chestnut_agouti"},
	} {
		for frame := 0; frame < 32; frame++ {
			got := currentFrameForVariant(stateNibble, frame, variant)
			if got < hopStart || got >= hopStart+4 {
				t.Fatalf("%s nibble frame %d = %d, want stable action fallback", variant.ID, frame, got)
			}
			got = currentFrameForVariant(stateGroom, frame, variant)
			if got < groomStart || got >= groomStart+groomFrames {
				t.Fatalf("%s groom frame %d = %d, want groom fallback", variant.ID, frame, got)
			}
		}
	}

	hamster := coatVariant{ID: "hamster_golden_syrian"}
	for frame := 0; frame < 12; frame++ {
		got := currentFrameForVariant(stateNibble, frame, hamster)
		if got < nibbleStart || got >= nibbleStart+3 {
			t.Fatalf("hamster nibble frame %d = %d, want original nibble frames", frame, got)
		}
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
		pets: []desktopPet{
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
		pets: []desktopPet{
			{state: stateWalk, stateTicks: 12, item: noItem},
		},
	}

	a.onTyping()
	if got := a.pets[0].state; got == stateWheel {
		t.Fatalf("typing in random mode started wheel state")
	}
}

func TestRuntimeCatalogIsReleaseScopedToPreviewAnimals(t *testing.T) {
	wantIDs := []string{
		"chinchilla_standard_gray",
		"hamster_golden_syrian",
		"djungarian_hamster",
		"campbell_hamster",
		"macaroni_mouse_tan",
		"sugar_glider_gray",
		"rabbit_chestnut_agouti",
		"holland_lop_broken_orange",
		"netherland_dwarf_chestnut",
		"himalayan_rabbit",
		"gecko_gray_brown",
		"guinea_pig_tricolor",
		"fancy_rat_hooded",
		"albino_chipmunk",
		"richardsons_ground_squirrel",
		"yorkshire_terrier_longcoat",
	}
	if got := len(variants); got != len(wantIDs) {
		t.Fatalf("runtime variants = %d, want %d", got, len(wantIDs))
	}
	for i, variant := range variants {
		if variant.ID != wantIDs[i] {
			t.Fatalf("runtime variant[%d] = %q, want %q", i, variant.ID, wantIDs[i])
		}
		if variant.SpeciesID == "degu" {
			t.Fatalf("runtime variants include degu: %+v", variant)
		}
	}
}

func TestSettingsRoundTripPersistsCoreOptions(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)

	a := &petApp{
		variant:        4,
		coatMode:       coatSelected,
		selectedCoats:  [maxPetCount]int{1, 3, 5, 7, 9, 0, 2, 4, 6, 8},
		petSizes:       [maxPetCount]int{80, 90, 100, 110, 120, 70, 80, 90, 100, 110},
		petNames:       [maxPetCount]string{"モカ", "Sora", "  Nagi  ", "", "", "", "", "", "", ""},
		nameLabels:     true,
		speed:          5,
		mode:           modeKeyboard,
		petCount:       10,
		wheelEnabled:   false,
		bidirectional:  false,
		positionMode:   positionScreenBottom,
		overlayOffsetY: 24,
		displayIndex:   0,
		displayScope:   displayScopeSingle,
		displaySpanEnd: 0,
		walkRangeStart: 15,
		walkRangeEnd:   85,
		lang:           langEnglish,
		settingsX:      220,
		settingsY:      180,
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
	if saved.Version != settingsVersion || saved.PetCount != 10 || saved.Mode != int(modeKeyboard) {
		t.Fatalf("saved settings = %+v, want version %d petCount 10 keyboard mode", saved, settingsVersion)
	}
	if !saved.NameLabels {
		t.Fatalf("saved NameLabels = false, want true")
	}
	if saved.PositionMode == nil || *saved.PositionMode != int(positionScreenBottom) {
		t.Fatalf("saved PositionMode = %#v, want screen bottom", saved.PositionMode)
	}
	if saved.VerticalOffset == nil || *saved.VerticalOffset != 24 {
		t.Fatalf("saved VerticalOffset = %#v, want 24", saved.VerticalOffset)
	}
	if saved.DisplayScope == nil || *saved.DisplayScope != int(displayScopeSingle) || saved.DisplayIndex == nil || *saved.DisplayIndex != 0 || saved.DisplaySpanEnd == nil || *saved.DisplaySpanEnd != 0 {
		t.Fatalf("saved display selection = scope:%#v index:%#v span:%#v", saved.DisplayScope, saved.DisplayIndex, saved.DisplaySpanEnd)
	}
	if saved.WalkRangeStart == nil || *saved.WalkRangeStart != 15 || saved.WalkRangeEnd == nil || *saved.WalkRangeEnd != 85 {
		t.Fatalf("saved walk range = start:%#v end:%#v", saved.WalkRangeStart, saved.WalkRangeEnd)
	}
	if got := saved.PetNames[0]; got != "モカ" {
		t.Fatalf("saved pet name 0 = %q, want モカ", got)
	}
	if got := saved.PetNames[2]; got != "Nagi" {
		t.Fatalf("saved pet name 2 = %q, want sanitized Nagi", got)
	}
	if len(saved.PetSizes) != maxPetCount || saved.PetSizes[0] != 80 || saved.PetSizes[4] != 120 || saved.PetSizes[5] != 70 {
		t.Fatalf("saved pet sizes = %#v", saved.PetSizes)
	}

	b := &petApp{
		variant:        0,
		coatMode:       coatSelected,
		selectedCoats:  defaultSelectedCoats(),
		petSizes:       defaultPetSizes(),
		speed:          3,
		mode:           modeRandom,
		petCount:       2,
		wheelEnabled:   true,
		bidirectional:  true,
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: defaultOverlayOffsetY,
		displayScope:   displayScopeSingle,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
		lang:           langJapanese,
		settingsX:      120,
		settingsY:      120,
	}
	if err := b.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if b.variant != 4 || b.coatMode != a.coatMode || b.speed != a.speed || b.mode != a.mode || b.petCount != a.petCount {
		t.Fatalf("loaded scalar settings = variant:%d coat:%d speed:%d mode:%d count:%d", b.variant, b.coatMode, b.speed, b.mode, b.petCount)
	}
	if b.wheelEnabled != a.wheelEnabled || b.bidirectional != a.bidirectional || b.lang != a.lang {
		t.Fatalf("loaded flags = wheel:%v bidirectional:%v lang:%d", b.wheelEnabled, b.bidirectional, b.lang)
	}
	if b.nameLabels != a.nameLabels {
		t.Fatalf("loaded nameLabels = %v, want %v", b.nameLabels, a.nameLabels)
	}
	if b.positionMode != a.positionMode || b.overlayOffsetY != a.overlayOffsetY || b.displayScope != a.displayScope || b.displayIndex != a.displayIndex || b.displaySpanEnd != a.displaySpanEnd {
		t.Fatalf("loaded display settings = mode:%d offset:%d scope:%d index:%d span:%d", b.positionMode, b.overlayOffsetY, b.displayScope, b.displayIndex, b.displaySpanEnd)
	}
	if b.walkRangeStart != a.walkRangeStart || b.walkRangeEnd != a.walkRangeEnd {
		t.Fatalf("loaded walk range = %d-%d, want %d-%d", b.walkRangeStart, b.walkRangeEnd, a.walkRangeStart, a.walkRangeEnd)
	}
	wantCoats := [maxPetCount]int{1, 3, 5, 7, 9, 0, 2, 4, 6, 8}
	for i := 0; i < maxPetCount; i++ {
		if b.selectedCoats[i] != wantCoats[i] {
			t.Fatalf("selectedCoats[%d] = %d, want %d", i, b.selectedCoats[i], wantCoats[i])
		}
	}
	if b.petNames[0] != "モカ" || b.petNames[1] != "Sora" || b.petNames[2] != "Nagi" {
		t.Fatalf("loaded pet names = %#v", b.petNames[:3])
	}
	for i, want := range a.petSizes {
		if got := b.petSizes[i]; got != want {
			t.Fatalf("loaded petSizes[%d] = %d, want %d", i, got, want)
		}
	}
}

func TestVersionOneSettingsKeepOptionsButResetOldAnimalSelection(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("APPDATA", configRoot)
	dir := filepath.Join(configRoot, settingsDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir settings dir: %v", err)
	}
	path := filepath.Join(dir, settingsFileName)
	data := []byte(`{
  "version": 1,
  "variant": 8,
  "coatMode": 2,
  "selectedCoats": [8, 7, 6, 5, 4, 3, 2, 1, 0, 9],
  "speed": 5,
  "mode": 0,
  "petCount": 4,
  "wheelEnabled": false,
  "bidirectional": false,
  "language": 1,
  "nameLabels": true,
  "petNames": ["モカ"]
}`)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}

	a := &petApp{
		variant:       0,
		coatMode:      coatSelected,
		selectedCoats: defaultSelectedCoats(),
		speed:         3,
		mode:          modeRandom,
		petCount:      5,
		wheelEnabled:  true,
		bidirectional: true,
		lang:          langJapanese,
	}
	if err := a.loadSettings(); err != nil {
		t.Fatalf("loadSettings() error = %v", err)
	}
	if a.speed != 5 || a.mode != modeKeyboard || a.petCount != 4 || a.wheelEnabled || a.bidirectional || a.lang != langEnglish {
		t.Fatalf("loaded preserved settings = speed:%d mode:%d count:%d wheel:%v bidi:%v lang:%d", a.speed, a.mode, a.petCount, a.wheelEnabled, a.bidirectional, a.lang)
	}
	if a.variant != 0 || a.coatMode != coatSelected || a.selectedCoats != defaultSelectedCoats() {
		t.Fatalf("old animal selection was not reset: variant:%d coat:%d selected:%v", a.variant, a.coatMode, a.selectedCoats)
	}
	if !a.nameLabels || a.petNames[0] != "モカ" {
		t.Fatalf("loaded name settings = labels:%v names:%v", a.nameLabels, a.petNames[:1])
	}
	for i, size := range a.petSizes {
		if size != defaultPetSizePercent {
			t.Fatalf("legacy petSizes[%d] = %d, want default %d", i, size, defaultPetSizePercent)
		}
	}
}

func TestNormalizeWalkRangeKeepsMinimumSpan(t *testing.T) {
	start, end := normalizeWalkRange(48, 52)
	if end-start != minWalkRangeSpan {
		t.Fatalf("normalizeWalkRange narrow span = %d-%d, want %d point span", start, end, minWalkRangeSpan)
	}
	start, end = normalizeWalkRange(95, 10)
	if start != 10 || end != 95 {
		t.Fatalf("normalizeWalkRange reversed = %d-%d, want 10-95", start, end)
	}
}

func TestOverlayRectForAppliesScreenBottomOffsetAndWalkRange(t *testing.T) {
	a := &petApp{
		positionMode:   positionScreenBottom,
		overlayOffsetY: 24,
		walkRangeStart: 25,
		walkRangeEnd:   75,
	}
	work := win.RECT{Left: 0, Top: 0, Right: 1000, Bottom: 760}
	screen := win.RECT{Left: 0, Top: 0, Right: 1000, Bottom: 800}
	got := a.overlayRectFor(work, screen)
	if got.Left != 250 || got.Right != 750 {
		t.Fatalf("overlay x range = %d-%d, want 250-750", got.Left, got.Right)
	}
	if got.Top != int32(800-sceneH) || got.Bottom != 800 {
		t.Fatalf("overlay y range = %d-%d, want clamped to screen bottom", got.Top, got.Bottom)
	}
}

func TestOverlayRectForTaskbarOffsetStaysInsideScreen(t *testing.T) {
	a := &petApp{
		positionMode:   positionTaskbarEdge,
		overlayOffsetY: -20,
		walkRangeStart: defaultWalkRangeStart,
		walkRangeEnd:   defaultWalkRangeEnd,
	}
	work := win.RECT{Left: 100, Top: 50, Right: 900, Bottom: 700}
	screen := win.RECT{Left: 0, Top: 0, Right: 1000, Bottom: 800}
	got := a.overlayRectFor(work, screen)
	if got.Left != work.Left || got.Right != work.Right {
		t.Fatalf("overlay x = %d-%d, want work area %d-%d", got.Left, got.Right, work.Left, work.Right)
	}
	wantTop := int32(700 - sceneH - 20)
	if got.Top != wantTop || got.Bottom != wantTop+sceneH {
		t.Fatalf("overlay y = %d-%d, want %d-%d", got.Top, got.Bottom, wantTop, wantTop+sceneH)
	}
}

func TestPetScenePositionsDistributeFivePetsAcrossTwoDisplays(t *testing.T) {
	positions := petScenePositions(3840, 5, []sceneSegment{
		{Left: 0, Right: 1920},
		{Left: 1920, Right: 3840},
	})
	if len(positions) != 5 {
		t.Fatalf("positions = %d, want 5", len(positions))
	}
	mainCount := 0
	subCount := 0
	for i, x := range positions {
		switch {
		case x >= 0 && x+spriteW <= 1920:
			mainCount++
		case x >= 1920 && x+spriteW <= 3840:
			subCount++
		default:
			t.Fatalf("position[%d] = %d escapes monitor segments", i, x)
		}
	}
	if mainCount != 3 || subCount != 2 {
		t.Fatalf("pet distribution = main:%d sub:%d, want 3 and 2", mainCount, subCount)
	}
}

func TestPetScenePositionsAvoidMonitorGaps(t *testing.T) {
	positions := petScenePositions(3800, 4, []sceneSegment{
		{Left: 0, Right: 1600},
		{Left: 2200, Right: 3800},
	})
	for i, x := range positions {
		onLeft := x >= 0 && x+spriteW <= 1600
		onRight := x >= 2200 && x+spriteW <= 3800
		if !onLeft && !onRight {
			t.Fatalf("position[%d] = %d falls in the monitor gap or offscreen", i, x)
		}
	}
}

func TestSetPetCountPlacesAllPetsInsideCurrentScene(t *testing.T) {
	a := &petApp{
		sceneW:        3840,
		speed:         3,
		coatMode:      coatFixed,
		bidirectional: true,
		petCount:      2,
	}

	a.setPetCount(5)
	if len(a.pets) != 5 {
		t.Fatalf("pets = %d, want 5", len(a.pets))
	}
	subCount := 0
	for i, pet := range a.pets {
		if pet.x < 0 || pet.x+spriteW > a.sceneW {
			t.Fatalf("pet %d x = %d escapes scene width %d", i, pet.x, a.sceneW)
		}
		if pet.x >= 1920 {
			subCount++
		}
	}
	if subCount < 2 {
		t.Fatalf("sub-display pets = %d, want at least 2 after choosing 5 pets", subCount)
	}
}

func TestResetPositionDistributesPetsAcrossDetectedMultiMonitorSpan(t *testing.T) {
	areas := monitorAreasByPosition()
	if len(areas) < 2 {
		t.Skip("multi-monitor placement check requires at least two detected displays")
	}
	a := &petApp{
		speed:          3,
		coatMode:       coatFixed,
		bidirectional:  true,
		petCount:       5,
		displayScope:   displayScopeSpan,
		displayIndex:   0,
		displaySpanEnd: len(areas) - 1,
		positionMode:   positionTaskbarEdge,
		walkRangeEnd:   100,
	}
	a.resetPosition()
	overlay := a.overlayRect()
	segments := a.sceneSegmentsForOverlay(overlay)
	if len(segments) < 2 {
		t.Fatalf("detected display span produced %d visible segments, want at least 2", len(segments))
	}
	seen := make([]int, len(segments))
	for _, pet := range a.pets {
		for i, segment := range segments {
			if pet.x >= segment.Left && pet.x+spriteW <= segment.Right {
				seen[i]++
				break
			}
		}
	}
	if len(a.pets) >= len(segments) {
		for i, count := range seen {
			if count == 0 {
				t.Fatalf("segment %d received no pets; distribution=%v segments=%+v", i, seen, segments)
			}
		}
	}
}

func TestWalkRangeSummaryDescribesMultiDisplaySegments(t *testing.T) {
	a := &petApp{lang: langJapanese}
	segments := []sceneSegment{
		{Left: 0, Right: 1920},
		{Left: 1920, Right: 3840},
	}

	cases := []struct {
		name       string
		start, end int
		want       string
	}{
		{name: "all selected displays", start: 0, end: 100, want: "選択した画面ぜんぶ"},
		{name: "first display only", start: 0, end: 50, want: "画面1だけ"},
		{name: "second display only", start: 50, end: 100, want: "画面2だけ"},
		{name: "partial displays", start: 25, end: 75, want: "画面1-2の一部"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := a.walkRangeSummaryForSegments(tt.start, tt.end, segments); got != tt.want {
				t.Fatalf("summary = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDisplaySpanDefaultsWalkingRangeToAllDisplays(t *testing.T) {
	if len(monitorAreas()) < 2 {
		t.Skip("multi-monitor scope reset requires at least two detected displays")
	}
	a := &petApp{
		displayScope:   displayScopeSingle,
		displayIndex:   0,
		displaySpanEnd: 0,
		positionMode:   positionTaskbarEdge,
		walkRangeStart: 25,
		walkRangeEnd:   75,
		petCount:       2,
		speed:          3,
		coatMode:       coatFixed,
	}
	a.resetPosition()
	a.walkRangeStart = 25
	a.walkRangeEnd = 75
	a.setDisplayScope(displayScopeSpan)
	if a.walkRangeStart != 0 || a.walkRangeEnd != 100 {
		t.Fatalf("span walk range = %d-%d, want reset to 0-100", a.walkRangeStart, a.walkRangeEnd)
	}

	a.setWalkRange(50, 100)
	a.adjustDisplaySpan(1)
	if a.walkRangeStart != 50 || a.walkRangeEnd != 100 {
		t.Fatalf("adjusting span changed fine range = %d-%d, want 50-100", a.walkRangeStart, a.walkRangeEnd)
	}
}

func TestPetVariantRectsFitTenPetsInSettingsWindow(t *testing.T) {
	seen := map[[4]int]bool{}
	for i := 0; i < maxPetCount; i++ {
		numberRect, buttonRect := settingsPetVariantRects(i)
		sizeRect := settingsPetSizeRect(i)
		if buttonRect.Right > 708 || buttonRect.Bottom > 562 || sizeRect.Right > 708 || sizeRect.Bottom > 562 {
			t.Fatalf("pet variant button %d rect %+v overflows selected-coats panel", i, buttonRect)
		}
		if numberRect.Left < 238 || buttonRect.Left <= numberRect.Right || sizeRect.Left <= buttonRect.Right {
			t.Fatalf("pet variant %d row rects overlap or escape: number=%+v button=%+v size=%+v", i, numberRect, buttonRect, sizeRect)
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
		sizeRect := settingsPetSizeRect(i)
		if nameRect.Right >= coatRect.Left {
			t.Fatalf("pet %d name rect overlaps coat rect: name=%+v coat=%+v", i, nameRect, coatRect)
		}
		if coatRect.Right >= sizeRect.Left {
			t.Fatalf("pet %d coat rect overlaps size rect: coat=%+v size=%+v", i, coatRect, sizeRect)
		}
		if numberRect.Left < 238 || nameRect.Left <= numberRect.Right || sizeRect.Right > 708 || nameRect.Bottom > 562 || sizeRect.Bottom > 562 {
			t.Fatalf("pet %d name/coat/size row escapes panel: number=%+v name=%+v coat=%+v size=%+v", i, numberRect, nameRect, coatRect, sizeRect)
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
		pets: []desktopPet{
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

func TestResetPetAtEdgeReentersFromOppositeSideWithMatchingDirection(t *testing.T) {
	a := &petApp{
		sceneW:        500,
		speed:         3,
		coatMode:      coatSelected,
		selectedCoats: defaultSelectedCoats(),
		forage: []forageItem{
			{owner: 0, active: true},
			{owner: reservedItem, active: true},
		},
	}

	right := desktopPet{dir: 1, nextDir: 1, item: 0, carryKind: 2, state: stateCarry}
	a.resetPetAtEdge(0, &right, 1)
	if right.x > -spriteW || right.dir != 1 || right.nextDir != 1 {
		t.Fatalf("right-moving reset = x:%d dir:%d next:%d, want off-left and direction +1", right.x, right.dir, right.nextDir)
	}
	if right.item != noItem || right.carryKind != noItem || right.state != stateWalk {
		t.Fatalf("right-moving reset state = item:%d carry:%d state:%v, want cleared walk", right.item, right.carryKind, right.state)
	}
	if a.forage[0].owner != noItem {
		t.Fatalf("owned forage was not released: owner=%d", a.forage[0].owner)
	}

	left := desktopPet{dir: -1, nextDir: -1, item: noItem, state: stateWalk}
	a.resetPetAtEdge(1, &left, -1)
	if left.x < a.sceneW || left.dir != -1 || left.nextDir != -1 {
		t.Fatalf("left-moving reset = x:%d dir:%d next:%d, want off-right and direction -1", left.x, left.dir, left.nextDir)
	}
}

func TestForagePropsDisabledClearsPropsAndStopsAssignment(t *testing.T) {
	if foragePropsEnabled {
		t.Fatalf("foragePropsEnabled = true, want false for preview release")
	}
	a := &petApp{
		sceneW: 500,
		speed:  3,
		forage: []forageItem{
			{x: 100, kind: 2, owner: 0, active: true},
			{x: 160, kind: 1, owner: reservedItem, active: true},
		},
		pets: []desktopPet{
			{state: stateCarry, item: 0, carryKind: 2, dir: 1, nextDir: 1},
			{state: stateForage, item: 1, carryKind: noItem, dir: -1, nextDir: -1},
		},
	}

	a.clearForageItems()

	for i, item := range a.forage {
		if item.active || item.owner != noItem {
			t.Fatalf("forage %d = active:%v owner:%d, want cleared", i, item.active, item.owner)
		}
	}
	for i, pet := range a.pets {
		if pet.item != noItem || pet.carryKind != noItem || pet.state != stateWalk {
			t.Fatalf("pet %d = item:%d carry:%d state:%v, want cleared walk", i, pet.item, pet.carryKind, pet.state)
		}
	}

	a.forage = []forageItem{{x: 140, kind: 2, owner: noItem, active: true}}
	p := desktopPet{state: stateWalk, item: noItem, carryKind: noItem, x: 40, dir: 1}
	if a.maybeAssignForageTarget(&p) {
		t.Fatalf("maybeAssignForageTarget assigned forage while props are disabled")
	}
	if p.item != noItem || p.state != stateWalk {
		t.Fatalf("pet after disabled assignment = item:%d state:%v, want unchanged", p.item, p.state)
	}
}

func TestTickPetMovesByDirectionAndWrapsPastEdges(t *testing.T) {
	a := &petApp{
		sceneW:        240,
		speed:         3,
		coatMode:      coatSelected,
		selectedCoats: defaultSelectedCoats(),
	}

	right := desktopPet{x: 20, dir: 1, nextDir: 1, state: stateWalk, moveSpeed: 4, stateTicks: 10, item: noItem, carryKind: noItem}
	a.tickPet(0, &right)
	if right.x != 24 {
		t.Fatalf("right-moving tick x = %d, want 24", right.x)
	}

	left := desktopPet{x: 20, dir: -1, nextDir: -1, state: stateWalk, moveSpeed: 4, stateTicks: 10, item: noItem, carryKind: noItem}
	a.tickPet(0, &left)
	if left.x != 16 {
		t.Fatalf("left-moving tick x = %d, want 16", left.x)
	}

	right.x = a.sceneW + 9
	a.tickPet(0, &right)
	if right.x > -spriteW || right.dir != 1 {
		t.Fatalf("right edge wrap = x:%d dir:%d, want off-left dir +1", right.x, right.dir)
	}

	left.x = -spriteW - 9
	a.tickPet(0, &left)
	if left.x < a.sceneW || left.dir != -1 {
		t.Fatalf("left edge wrap = x:%d dir:%d, want off-right dir -1", left.x, left.dir)
	}
}

func TestDrawFacingImageMirrorsNegativeDirection(t *testing.T) {
	red := color.RGBA{R: 255, A: 255}
	blue := color.RGBA{B: 255, A: 255}
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	src.SetRGBA(0, 0, red)
	src.SetRGBA(1, 0, blue)

	dst := image.NewRGBA(src.Bounds())
	drawFacingImage(dst, src, dst.Bounds(), 1)
	if got := dst.RGBAAt(0, 0); got != red {
		t.Fatalf("positive left pixel = %#v, want %#v", got, red)
	}
	if got := dst.RGBAAt(1, 0); got != blue {
		t.Fatalf("positive right pixel = %#v, want %#v", got, blue)
	}

	dst = image.NewRGBA(src.Bounds())
	drawFacingImage(dst, src, dst.Bounds(), -1)
	if got := dst.RGBAAt(0, 0); got != blue {
		t.Fatalf("negative left pixel = %#v, want %#v", got, blue)
	}
	if got := dst.RGBAAt(1, 0); got != red {
		t.Fatalf("negative right pixel = %#v, want %#v", got, red)
	}
}

func TestDrawDirectionForVariantCompensatesLeftFacingSource(t *testing.T) {
	sugarGlider := coatVariant{ID: "sugar_glider_gray"}
	hamster := coatVariant{ID: "hamster_golden_syrian"}

	if got := drawDirectionForVariant(1, sugarGlider); got != -1 {
		t.Fatalf("sugar glider right-moving draw direction = %d, want -1", got)
	}
	if got := drawDirectionForVariant(-1, sugarGlider); got != 1 {
		t.Fatalf("sugar glider left-moving draw direction = %d, want 1", got)
	}
	if got := drawDirectionForVariant(1, hamster); got != 1 {
		t.Fatalf("hamster right-moving draw direction = %d, want 1", got)
	}
	if got := drawDirectionForVariant(-1, hamster); got != -1 {
		t.Fatalf("hamster left-moving draw direction = %d, want -1", got)
	}
}

func TestFixedCoatModeRefreshesAllPets(t *testing.T) {
	a := &petApp{
		variant: 99,
		pets: []desktopPet{
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
		pets: []desktopPet{
			{variant: 0},
			{variant: 0},
			{variant: 0},
		},
	}

	a.setCoatMode(coatSelected)

	want := []int{0, 3, 5}
	for i := range []int{0, 1, 2} {
		if got := a.pets[i].variant; got != want[i] {
			t.Fatalf("pet %d variant = %d, want %d", i, got, want[i])
		}
	}
	a.setSelectedVariant(1, 8)
	if got, want := a.pets[1].variant, 8; got != want {
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
		sceneW:   800,
		petSizes: defaultPetSizes(),
		pets: []desktopPet{
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

func TestPetSizeAffectsHitTestingAndBounds(t *testing.T) {
	a := &petApp{
		sceneW:   800,
		petSizes: defaultPetSizes(),
		pets: []desktopPet{
			{x: 100, laneOffset: 0, state: stateWalk},
		},
	}
	a.setPetSize(0, 120)
	w, h := a.petSpriteSize(0)
	if w != 115 || h != 76 {
		t.Fatalf("petSpriteSize(120%%) = %dx%d, want 115x76", w, h)
	}
	if got := a.petAtScenePoint(100+w-8, sceneH-h+8); got != 0 {
		t.Fatalf("petAtScenePoint on enlarged pet = %d, want 0", got)
	}
	a.pets[0].x = 760
	a.setPetSize(0, 120)
	if a.pets[0].x+w > a.sceneW {
		t.Fatalf("setPetSize did not clamp x: x=%d w=%d scene=%d", a.pets[0].x, w, a.sceneW)
	}
}

func TestNormalizePetSizePercent(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{0, defaultPetSizePercent},
		{64, minPetSizePercent},
		{86, 90},
		{119, maxPetSizePercent},
		{200, maxPetSizePercent},
	}
	for _, tt := range tests {
		if got := normalizePetSizePercent(tt.in); got != tt.want {
			t.Fatalf("normalizePetSizePercent(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestShowPetReactionRefreshesExistingBubble(t *testing.T) {
	a := &petApp{
		pets: []desktopPet{{state: stateWalk}},
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
		pets: []desktopPet{{state: stateWalk}},
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

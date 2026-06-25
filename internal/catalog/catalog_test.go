package catalog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogInvariants(t *testing.T) {
	if len(Variants) != 100 {
		t.Fatalf("variants = %d, want exactly 100", len(Variants))
	}
	if got := len(SeedVariants()); got != 89 {
		t.Fatalf("seed variants = %d, want 89", got)
	}
	if got := len(DeguVariantIDs()); got != 11 {
		t.Fatalf("degu variants = %d, want 11", got)
	}

	species := map[string]Species{}
	for _, entry := range SpeciesList {
		if entry.ID == "" || entry.Label == "" || entry.Profile == "" {
			t.Fatalf("incomplete species entry: %+v", entry)
		}
		if _, exists := species[entry.ID]; exists {
			t.Fatalf("duplicate species id %q", entry.ID)
		}
		species[entry.ID] = entry
	}

	ids := map[string]bool{}
	spriteBases := map[string]bool{}
	for _, variant := range Variants {
		if variant.ID == "" || variant.SpeciesID == "" || variant.LabelEN == "" || variant.SpriteBase == "" {
			t.Fatalf("incomplete variant: %+v", variant)
		}
		if variant.BreedOrMorph == "" || variant.Color == "" || variant.PopularityTier == 0 || variant.MotionProfile == "" || variant.SourceStatus == "" {
			t.Fatalf("variant %q missing catalog metadata: %+v", variant.ID, variant)
		}
		if _, ok := species[variant.SpeciesID]; !ok {
			t.Fatalf("variant %q references unknown species %q", variant.ID, variant.SpeciesID)
		}
		if ids[variant.ID] {
			t.Fatalf("duplicate variant id %q", variant.ID)
		}
		ids[variant.ID] = true
		if spriteBases[variant.SpriteBase] {
			t.Fatalf("duplicate sprite base %q", variant.SpriteBase)
		}
		spriteBases[variant.SpriteBase] = true
		if variant.SeedStage && variant.SourcePath == "" && variant.Shape == "" {
			t.Fatalf("seed variant %q has neither source path nor procedural shape", variant.ID)
		}
		if variant.SeedStage && variant.SourceStatus != SourceStatusPrototypeOnly && variant.SourceStatus != SourceStatusImageGenQueued && variant.SourceStatus != SourceStatusMotionDraft && variant.SourceStatus != SourceStatusMotionAccepted {
			t.Fatalf("seed variant %q source status = %q, want prototype/imagegen/motion status", variant.ID, variant.SourceStatus)
		}
		if (variant.SourceStatus == SourceStatusMotionDraft || variant.SourceStatus == SourceStatusMotionAccepted) && variant.MotionSourcePath == "" {
			t.Fatalf("motion source variant %q has no motion source path", variant.ID)
		}
		if variant.MotionSourcePath != "" && variant.SourceStatus != SourceStatusMotionDraft && variant.SourceStatus != SourceStatusMotionAccepted {
			t.Fatalf("variant %q has motion source path with source status %q", variant.ID, variant.SourceStatus)
		}
		if !variant.SeedStage && variant.SourceStatus != SourceStatusDeguMotion {
			t.Fatalf("non-seed variant %q source status = %q, want degu motion source", variant.ID, variant.SourceStatus)
		}
	}
}

func TestRuntimeVariantsAreReleaseScoped(t *testing.T) {
	runtime := RuntimeVariants()
	wantIDs := []string{
		"chinchilla_standard_gray",
		"hamster_golden_syrian",
		"macaroni_mouse_tan",
		"sugar_glider_gray",
		"rabbit_chestnut_agouti",
	}
	if got := len(runtime); got != len(wantIDs) {
		t.Fatalf("runtime variants = %d, want %d release-scoped variants", got, len(wantIDs))
	}
	for i, variant := range runtime {
		if variant.ID != wantIDs[i] {
			t.Fatalf("runtime variant[%d] = %q, want %q", i, variant.ID, wantIDs[i])
		}
		if variant.SpeciesID == "degu" {
			t.Fatalf("runtime variants must not include degu: %+v", variant)
		}
		if variant.SourceStatus != SourceStatusMotionAccepted {
			t.Fatalf("runtime variant %q source status = %q, want accepted", variant.ID, variant.SourceStatus)
		}
	}
}

func TestRuntimeSpritesMatchAcceptedMotionSources(t *testing.T) {
	const runtimeMotionSets = 10
	for _, variant := range RuntimeVariants() {
		sourcePaths := expectedRuntimeMotionSources(t, variant.MotionSourcePath, runtimeMotionSets)
		for set := 0; set < runtimeMotionSets; set++ {
			runtimePath := repoPath("assets", "sprites", fmt.Sprintf("%s_set%02d.png", variant.SpriteBase, set))
			runtimeData, err := os.ReadFile(runtimePath)
			if err != nil {
				t.Fatalf("read runtime sprite %s: %v", runtimePath, err)
			}
			sourceData, err := os.ReadFile(sourcePaths[set])
			if err != nil {
				t.Fatalf("read motion source %s: %v", sourcePaths[set], err)
			}
			if !bytes.Equal(runtimeData, sourceData) {
				t.Fatalf("runtime sprite %s does not match accepted source %s", runtimePath, sourcePaths[set])
			}
		}
	}
}

func expectedRuntimeMotionSources(t *testing.T, set00Path string, sets int) []string {
	t.Helper()
	if set00Path == "" {
		t.Fatalf("runtime variant has no motion source path")
	}
	set00 := repoPath(filepath.FromSlash(set00Path))
	out := make([]string, sets)
	for set := 0; set < sets; set++ {
		candidate := set00
		if strings.Contains(set00, "set00") {
			candidate = strings.Replace(set00, "set00", fmt.Sprintf("set%02d", set), 1)
		}
		if _, err := os.Stat(candidate); err != nil {
			if set > 0 && os.IsNotExist(err) {
				candidate = set00
			} else {
				t.Fatalf("stat motion source %s: %v", candidate, err)
			}
		}
		out[set] = candidate
	}
	return out
}

func repoPath(parts ...string) string {
	all := append([]string{"..", ".."}, parts...)
	return filepath.Join(all...)
}

func TestWheelCapabilityByProfile(t *testing.T) {
	tests := []struct {
		species string
		want    bool
	}{
		{"degu", true},
		{"hamster", true},
		{"mouse", true},
		{"corn_snake", false},
		{"whites_tree_frog", false},
		{"gecko", false},
		{"tortoise", false},
		{"dog", false},
		{"capybara", false},
	}
	for _, tt := range tests {
		if got := WheelCapableSpecies(tt.species); got != tt.want {
			t.Fatalf("WheelCapableSpecies(%q) = %v, want %v", tt.species, got, tt.want)
		}
	}
}

func TestRequestedPopularVariantsArePresent(t *testing.T) {
	required := []string{
		"french_bulldog_fawn",
		"labrador_yellow",
		"golden_retriever_golden",
		"maine_coon_brown_tabby",
		"ragdoll_seal_bicolor",
		"holland_lop_broken_orange",
		"fancy_rat_hooded",
		"bearded_dragon_citrus",
		"corn_snake_amelanistic",
		"whites_tree_frog_green",
	}
	seen := map[string]bool{}
	for _, variant := range Variants {
		seen[variant.ID] = true
	}
	for _, id := range required {
		if !seen[id] {
			t.Fatalf("missing requested variant %q", id)
		}
	}
}

func TestMotionProfilesCoverCatalog(t *testing.T) {
	known := map[string]bool{
		MotionProfileDegu:               true,
		MotionProfileSmallRodentScurry:  true,
		MotionProfileRabbitHop:          true,
		MotionProfileDogTrot:            true,
		MotionProfileCatStalk:           true,
		MotionProfileGeckoCrawl:         true,
		MotionProfileTortoisePlod:       true,
		MotionProfileFerretSlink:        true,
		MotionProfileGuineaPigWaddle:    true,
		MotionProfileHedgehogShuffle:    true,
		MotionProfileSquirrelBound:      true,
		MotionProfileFoxTrot:            true,
		MotionProfileRedPandaAmble:      true,
		MotionProfileOtterSlide:         true,
		MotionProfileSugarGliderSkitter: true,
		MotionProfileCapybaraLumber:     true,
		MotionProfileSnakeSlither:       true,
		MotionProfileDragonPlod:         true,
		MotionProfileFrogHop:            true,
	}
	for _, variant := range Variants {
		if !known[MotionProfileForVariant(variant)] {
			t.Fatalf("variant %q has unknown motion profile %q", variant.ID, MotionProfileForVariant(variant))
		}
	}
}

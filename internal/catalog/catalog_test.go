package catalog

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogInvariants(t *testing.T) {
	if len(Variants) != 137 {
		t.Fatalf("variants = %d, want exactly 137", len(Variants))
	}
	if got := len(SeedVariants()); got != 126 {
		t.Fatalf("seed variants = %d, want 126", got)
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
		"chipmunk_striped",
		"gecko_leopard",
		"whites_tree_frog_blue",
		"budgerigar_green_yellow",
		"cockatiel_normal_gray",
		"java_sparrow_normal",
		"parrotlet_green",
		"parrotlet_blue_green",
		"lovebird_peach_faced",
		"ragdoll_seal_bicolor",
		"scottish_fold_silver_tabby",
		"french_bulldog_fawn",
		"maine_coon_brown_tabby",
		"domestic_shorthair_calico",
		"british_shorthair_blue",
		"toy_poodle_apricot",
		"munchkin_brown_tabby",
		"roborovski_hamster",
		"guinea_pig_russian_smoke_white",
		"quokka",
		"miniature_schnauzer_salt_pepper",
		"japanese_giant_salamander",
		"white_wagtail",
		"domestic_shorthair_tabby_white_stocky",
		"lionhead_rabbit_brown_white",
		"shoebill_stork",
		"leucistic_sugar_glider",
		"african_dormouse",
		"netherland_dwarf_himalayan",
		"american_flying_squirrel",
		"longhair_hamster_black_white",
		"djungarian_hamster_yellow",
		"djungarian_hamster_pearl_white",
		"fancy_rat_blue_hooded",
		"fancy_rat_chocolate_self",
		"fancy_rat_cream_agouti",
		"rabbit_gray",
		"african_fat_tailed_gecko",
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
			runtimeImg, err := readPNG(runtimePath)
			if err != nil {
				t.Fatalf("read runtime sprite %s: %v", runtimePath, err)
			}
			sourceImg, err := readPNG(sourcePaths[set])
			if err != nil {
				t.Fatalf("read motion source %s: %v", sourcePaths[set], err)
			}
			if !imagesEqual(runtimeImg, sourceImg) {
				t.Fatalf("runtime sprite %s does not match accepted source %s", runtimePath, sourcePaths[set])
			}
		}
	}
}

func readPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func imagesEqual(a image.Image, b image.Image) bool {
	ab := a.Bounds()
	bb := b.Bounds()
	if ab != bb {
		return false
	}
	for y := ab.Min.Y; y < ab.Max.Y; y++ {
		for x := ab.Min.X; x < ab.Max.X; x++ {
			ar, ag, ablu, aa := a.At(x, y).RGBA()
			br, bg, bblu, ba := b.At(x, y).RGBA()
			if ar != br || ag != bg || ablu != bblu || aa != ba {
				return false
			}
		}
	}
	return true
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

func TestWheelCapabilityIsLimitedToChinchillaAndHamster(t *testing.T) {
	wantBySpecies := map[string]bool{
		"chinchilla":     true,
		"hamster":        true,
		"degu":           false,
		"macaroni_mouse": false,
		"mouse":          false,
		"sugar_glider":   false,
		"rabbit":         false,
		"gecko":          false,
	}
	for species, want := range wantBySpecies {
		if got := WheelCapableSpecies(species); got != want {
			t.Fatalf("WheelCapableSpecies(%q) = %v, want %v", species, got, want)
		}
	}
	if WheelCapableMotionProfile(MotionProfileSmallRodentScurry) {
		t.Fatalf("small-rodent motion profile must not imply wheel capability")
	}
	wantByVariant := map[string]bool{
		"chinchilla_standard_gray": true,
		"hamster_golden_syrian":    true,
		"macaroni_mouse_tan":       false,
		"sugar_glider_gray":        false,
		"rabbit_chestnut_agouti":   false,
		"wild_agouti":              false,
	}
	for id, want := range wantByVariant {
		variant, ok := VariantByID(id)
		if !ok {
			t.Fatalf("missing test variant %q", id)
		}
		if got := WheelCapableVariant(variant); got != want {
			t.Fatalf("WheelCapableVariant(%q) = %v, want %v", id, got, want)
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
		MotionProfileBirdHop:            true,
	}
	for _, variant := range Variants {
		if !known[MotionProfileForVariant(variant)] {
			t.Fatalf("variant %q has unknown motion profile %q", variant.ID, MotionProfileForVariant(variant))
		}
	}
}

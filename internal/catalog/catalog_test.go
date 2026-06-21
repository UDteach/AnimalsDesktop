package catalog

import "testing"

func TestCatalogInvariants(t *testing.T) {
	if len(Variants) < 64 {
		t.Fatalf("variants = %d, want at least 64", len(Variants))
	}
	if got := len(SeedVariants()); got < 53 {
		t.Fatalf("seed variants = %d, want at least 53", got)
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
	}
}

func TestWheelCapabilityByProfile(t *testing.T) {
	tests := []struct {
		species string
		want    bool
	}{
		{"degu", true},
		{"hamster", true},
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

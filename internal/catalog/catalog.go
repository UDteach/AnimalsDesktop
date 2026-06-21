package catalog

type Species struct {
	ID      string
	Label   string
	Profile string
}

type Variant struct {
	ID         string
	SpeciesID  string
	LabelEN    string
	LabelJA    string
	SpriteBase string
	SeedStage  bool
	SourcePath string
}

var SpeciesList = []Species{
	{ID: "degu", Label: "Degu", Profile: "degu"},
	{ID: "chinchilla", Label: "Chinchilla", Profile: "small-mammal"},
	{ID: "macaroni_mouse", Label: "Macaroni mouse / fat-tailed gerbil", Profile: "small-mammal"},
	{ID: "rabbit", Label: "Rabbit", Profile: "hopper"},
	{ID: "dog", Label: "Small dog", Profile: "companion"},
	{ID: "cat", Label: "Cat", Profile: "companion"},
	{ID: "gecko", Label: "Gecko", Profile: "low-crawler"},
	{ID: "hamster", Label: "Hamster", Profile: "small-mammal"},
}

var Variants = []Variant{
	{ID: "wild_agouti", SpeciesID: "degu", LabelEN: "Degu - wild agouti", LabelJA: "Degu - wild agouti", SpriteBase: "degu_wild_agouti"},
	{ID: "black", SpeciesID: "degu", LabelEN: "Degu - black", LabelJA: "Degu - black", SpriteBase: "degu_black"},
	{ID: "blue", SpeciesID: "degu", LabelEN: "Degu - blue (slate gray)", LabelJA: "Degu - blue (slate gray)", SpriteBase: "degu_blue"},
	{ID: "gray", SpeciesID: "degu", LabelEN: "Degu - gray", LabelJA: "Degu - gray", SpriteBase: "degu_gray"},
	{ID: "white_cream", SpeciesID: "degu", LabelEN: "Degu - white / cream", LabelJA: "Degu - white / cream", SpriteBase: "degu_white_cream"},
	{ID: "sand_champagne", SpeciesID: "degu", LabelEN: "Degu - sand / champagne", LabelJA: "Degu - sand / champagne", SpriteBase: "degu_sand_champagne"},
	{ID: "chocolate", SpeciesID: "degu", LabelEN: "Degu - chocolate", LabelJA: "Degu - chocolate", SpriteBase: "degu_chocolate"},
	{ID: "black_pied", SpeciesID: "degu", LabelEN: "Degu - black pied", LabelJA: "Degu - black pied", SpriteBase: "degu_black_pied"},
	{ID: "agouti_pied", SpeciesID: "degu", LabelEN: "Degu - agouti pied", LabelJA: "Degu - agouti pied", SpriteBase: "degu_agouti_pied"},
	{ID: "blue_pied", SpeciesID: "degu", LabelEN: "Degu - blue pied (slate gray)", LabelJA: "Degu - blue pied (slate gray)", SpriteBase: "degu_blue_pied"},
	{ID: "cream_pied", SpeciesID: "degu", LabelEN: "Degu - cream pied", LabelJA: "Degu - cream pied", SpriteBase: "degu_cream_pied"},
	{
		ID:         "chinchilla_standard_gray",
		SpeciesID:  "chinchilla",
		LabelEN:    "Chinchilla - standard gray",
		LabelJA:    "Chinchilla - standard gray",
		SpriteBase: "chinchilla_standard_gray",
		SeedStage:  true,
		SourcePath: "docs/art-source/chinchilla/chinchilla-standard-gray-source-truth-transparent.png",
	},
	{
		ID:         "macaroni_mouse_tan",
		SpeciesID:  "macaroni_mouse",
		LabelEN:    "Macaroni mouse - tan",
		LabelJA:    "Macaroni mouse - tan",
		SpriteBase: "macaroni_mouse_tan",
		SeedStage:  true,
		SourcePath: "docs/art-source/macaroni-mouse/macaroni-mouse-source-truth.png",
	},
	{
		ID:         "rabbit_chestnut_agouti",
		SpeciesID:  "rabbit",
		LabelEN:    "Rabbit - chestnut agouti",
		LabelJA:    "Rabbit - chestnut agouti",
		SpriteBase: "rabbit_chestnut_agouti",
		SeedStage:  true,
		SourcePath: "docs/source-truth/rabbit-source-truth.png",
	},
	{
		ID:         "dog_cream_tan",
		SpeciesID:  "dog",
		LabelEN:    "Small dog - cream and tan",
		LabelJA:    "Small dog - cream and tan",
		SpriteBase: "dog_cream_tan",
		SeedStage:  true,
		SourcePath: "docs/art-intake/dog/dog-source-truth-transparent.png",
	},
	{
		ID:         "cat_brown_tabby",
		SpeciesID:  "cat",
		LabelEN:    "Cat - brown tabby",
		LabelJA:    "Cat - brown tabby",
		SpriteBase: "cat_brown_tabby",
		SeedStage:  true,
		SourcePath: "docs/art-source/cat/cat-kijitora-source-truth-transparent.png",
	},
	{
		ID:         "gecko_gray_brown",
		SpeciesID:  "gecko",
		LabelEN:    "Gecko - gray brown",
		LabelJA:    "Gecko - gray brown",
		SpriteBase: "gecko_gray_brown",
		SeedStage:  true,
		SourcePath: "docs/art-source/gecko/gecko-source-truth.png",
	},
	{
		ID:         "hamster_golden_syrian",
		SpeciesID:  "hamster",
		LabelEN:    "Hamster - golden Syrian",
		LabelJA:    "Hamster - golden Syrian",
		SpriteBase: "hamster_golden_syrian",
		SeedStage:  true,
		SourcePath: "docs/art-source/hamster/hamster-source-truth.png",
	},
}

func DeguVariantIDs() []string {
	ids := make([]string, 0, len(Variants))
	for _, variant := range Variants {
		if variant.SpeciesID == "degu" {
			ids = append(ids, variant.ID)
		}
	}
	return ids
}

func SeedVariants() []Variant {
	out := make([]Variant, 0)
	for _, variant := range Variants {
		if variant.SeedStage {
			out = append(out, variant)
		}
	}
	return out
}

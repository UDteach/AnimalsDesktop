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
	Shape      string
	TintHex    string
	AccentHex  string
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
	{ID: "ferret", Label: "Ferret", Profile: "small-mammal"},
	{ID: "guinea_pig", Label: "Guinea pig", Profile: "small-mammal"},
	{ID: "hedgehog", Label: "Hedgehog", Profile: "small-mammal"},
	{ID: "squirrel", Label: "Squirrel", Profile: "small-mammal"},
	{ID: "fox", Label: "Fox", Profile: "companion"},
	{ID: "red_panda", Label: "Red panda", Profile: "companion"},
	{ID: "otter", Label: "Otter", Profile: "low-crawler"},
	{ID: "sugar_glider", Label: "Sugar glider", Profile: "small-mammal"},
	{ID: "capybara", Label: "Capybara", Profile: "large-mammal"},
	{ID: "tortoise", Label: "Tortoise", Profile: "low-crawler"},
}

const (
	srcChinchilla    = "docs/art-source/chinchilla/chinchilla-standard-gray-source-truth-transparent.png"
	srcMacaroniMouse = "docs/art-source/macaroni-mouse/macaroni-mouse-source-truth.png"
	srcRabbit        = "docs/source-truth/rabbit-source-truth.png"
	srcDog           = "docs/art-intake/dog/dog-source-truth-transparent.png"
	srcCat           = "docs/art-source/cat/cat-kijitora-source-truth-transparent.png"
	srcGecko         = "docs/art-source/gecko/gecko-source-truth.png"
	srcHamster       = "docs/art-source/hamster/hamster-source-truth.png"
)

var Variants = []Variant{
	deguVariant("wild_agouti", "wild agouti", "degu_wild_agouti"),
	deguVariant("black", "black", "degu_black"),
	deguVariant("blue", "blue (slate gray)", "degu_blue"),
	deguVariant("gray", "gray", "degu_gray"),
	deguVariant("white_cream", "white / cream", "degu_white_cream"),
	deguVariant("sand_champagne", "sand / champagne", "degu_sand_champagne"),
	deguVariant("chocolate", "chocolate", "degu_chocolate"),
	deguVariant("black_pied", "black pied", "degu_black_pied"),
	deguVariant("agouti_pied", "agouti pied", "degu_agouti_pied"),
	deguVariant("blue_pied", "blue pied (slate gray)", "degu_blue_pied"),
	deguVariant("cream_pied", "cream pied", "degu_cream_pied"),

	sourceVariant("chinchilla_standard_gray", "chinchilla", "Chinchilla - standard gray", "chinchilla_standard_gray", srcChinchilla, "", ""),
	sourceVariant("chinchilla_beige", "chinchilla", "Chinchilla - beige", "chinchilla_beige", srcChinchilla, "c8b184", ""),
	sourceVariant("chinchilla_ebony", "chinchilla", "Chinchilla - ebony", "chinchilla_ebony", srcChinchilla, "3e3d39", ""),
	sourceVariant("chinchilla_white_mosaic", "chinchilla", "Chinchilla - white mosaic", "chinchilla_white_mosaic", srcChinchilla, "e6e1d3", "8f9690"),

	sourceVariant("macaroni_mouse_tan", "macaroni_mouse", "Macaroni mouse - tan", "macaroni_mouse_tan", srcMacaroniMouse, "", ""),
	sourceVariant("macaroni_mouse_gray", "macaroni_mouse", "Macaroni mouse - gray", "macaroni_mouse_gray", srcMacaroniMouse, "8f8f84", ""),
	sourceVariant("macaroni_mouse_cream", "macaroni_mouse", "Macaroni mouse - cream", "macaroni_mouse_cream", srcMacaroniMouse, "dcc698", ""),

	sourceVariant("rabbit_chestnut_agouti", "rabbit", "Rabbit - chestnut agouti", "rabbit_chestnut_agouti", srcRabbit, "", ""),
	sourceVariant("rabbit_black", "rabbit", "Rabbit - black", "rabbit_black", srcRabbit, "2f2c29", ""),
	sourceVariant("rabbit_white", "rabbit", "Rabbit - white", "rabbit_white", srcRabbit, "ebe7dc", ""),
	sourceVariant("rabbit_blue_gray", "rabbit", "Rabbit - blue gray", "rabbit_blue_gray", srcRabbit, "777f86", ""),
	sourceVariant("rabbit_fawn", "rabbit", "Rabbit - fawn", "rabbit_fawn", srcRabbit, "bf8754", ""),

	sourceVariant("dog_cream_tan", "dog", "Small dog - cream and tan", "dog_cream_tan", srcDog, "", ""),
	sourceVariant("dog_black_tan", "dog", "Small dog - black and tan", "dog_black_tan", srcDog, "3a3026", "b8874c"),
	sourceVariant("dog_white", "dog", "Small dog - white", "dog_white", srcDog, "e8e1d0", ""),
	sourceVariant("dog_sable", "dog", "Small dog - sable", "dog_sable", srcDog, "9a6839", ""),
	sourceVariant("dog_gray", "dog", "Small dog - gray", "dog_gray", srcDog, "8d8d88", ""),

	sourceVariant("cat_brown_tabby", "cat", "Cat - brown tabby", "cat_brown_tabby", srcCat, "", ""),
	sourceVariant("cat_black", "cat", "Cat - black", "cat_black", srcCat, "252321", ""),
	sourceVariant("cat_white", "cat", "Cat - white", "cat_white", srcCat, "eee8dd", ""),
	sourceVariant("cat_orange_tabby", "cat", "Cat - orange tabby", "cat_orange_tabby", srcCat, "c87830", ""),
	sourceVariant("cat_gray", "cat", "Cat - gray", "cat_gray", srcCat, "818783", ""),
	sourceVariant("cat_cream", "cat", "Cat - cream", "cat_cream", srcCat, "d7b983", ""),

	sourceVariant("gecko_gray_brown", "gecko", "Gecko - gray brown", "gecko_gray_brown", srcGecko, "", ""),
	sourceVariant("gecko_leopard", "gecko", "Gecko - leopard", "gecko_leopard", srcGecko, "c99b4a", "4d3a22"),
	sourceVariant("gecko_tangerine", "gecko", "Gecko - tangerine", "gecko_tangerine", srcGecko, "de8a34", ""),
	sourceVariant("gecko_blizzard", "gecko", "Gecko - blizzard", "gecko_blizzard", srcGecko, "d8d2c5", ""),
	sourceVariant("gecko_albino", "gecko", "Gecko - albino", "gecko_albino", srcGecko, "e0cfa8", ""),

	sourceVariant("hamster_golden_syrian", "hamster", "Hamster - golden Syrian", "hamster_golden_syrian", srcHamster, "", ""),
	sourceVariant("hamster_cream", "hamster", "Hamster - cream", "hamster_cream", srcHamster, "e2c18b", ""),
	sourceVariant("hamster_black_banded", "hamster", "Hamster - black banded", "hamster_black_banded", srcHamster, "3f342d", "ddd2bd"),
	sourceVariant("hamster_white", "hamster", "Hamster - white", "hamster_white", srcHamster, "eee7d8", ""),
	sourceVariant("hamster_cinnamon", "hamster", "Hamster - cinnamon", "hamster_cinnamon", srcHamster, "b56b38", ""),

	shapeVariant("ferret_sable", "ferret", "Ferret - sable", "ferret", "8b6746", "ece0c8"),
	shapeVariant("ferret_albino", "ferret", "Ferret - albino", "ferret", "eadcc7", "c58b78"),
	shapeVariant("ferret_champagne", "ferret", "Ferret - champagne", "ferret", "c7a476", "f1e1c1"),
	shapeVariant("guinea_pig_tricolor", "guinea_pig", "Guinea pig - tricolor", "guinea_pig", "b46c32", "f0e6d2"),
	shapeVariant("guinea_pig_cream", "guinea_pig", "Guinea pig - cream", "guinea_pig", "d9b879", "f1dfb4"),
	shapeVariant("guinea_pig_black", "guinea_pig", "Guinea pig - black", "guinea_pig", "2b2926", "d8cdb9"),
	shapeVariant("hedgehog_salt_pepper", "hedgehog", "Hedgehog - salt and pepper", "hedgehog", "81786a", "e6dcc6"),
	shapeVariant("hedgehog_cinnamon", "hedgehog", "Hedgehog - cinnamon", "hedgehog", "b06e3e", "ebdcc2"),
	shapeVariant("squirrel_red", "squirrel", "Squirrel - red", "squirrel", "b45f2f", "e9c193"),
	shapeVariant("squirrel_gray", "squirrel", "Squirrel - gray", "squirrel", "85857c", "d8d4c8"),
	shapeVariant("fox_red", "fox", "Fox - red", "fox", "c46832", "eee2c9"),
	shapeVariant("fox_silver", "fox", "Fox - silver", "fox", "4b4d4c", "d6d3c9"),
	shapeVariant("red_panda_classic", "red_panda", "Red panda - classic", "red_panda", "b85f31", "efe1c5"),
	shapeVariant("red_panda_dark", "red_panda", "Red panda - dark", "red_panda", "6f3f2f", "e9d9bc"),
	shapeVariant("otter_brown", "otter", "Otter - brown", "otter", "6b4a31", "d4bb91"),
	shapeVariant("sugar_glider_gray", "sugar_glider", "Sugar glider - gray", "sugar_glider", "8c8d88", "f1eee4"),
	shapeVariant("capybara_brown", "capybara", "Capybara - brown", "capybara", "8a603d", "c79d6a"),
	shapeVariant("capybara_sand", "capybara", "Capybara - sand", "capybara", "b98a58", "e2c69c"),
	shapeVariant("tortoise_olive", "tortoise", "Tortoise - olive", "tortoise", "6f7146", "3f3d2a"),
	shapeVariant("tortoise_dark_shell", "tortoise", "Tortoise - dark shell", "tortoise", "45452f", "8b8052"),
}

func deguVariant(id string, label string, spriteBase string) Variant {
	return Variant{
		ID:         id,
		SpeciesID:  "degu",
		LabelEN:    "Degu - " + label,
		LabelJA:    "Degu - " + label,
		SpriteBase: spriteBase,
	}
}

func sourceVariant(id string, speciesID string, label string, spriteBase string, sourcePath string, tintHex string, accentHex string) Variant {
	return Variant{
		ID:         id,
		SpeciesID:  speciesID,
		LabelEN:    label,
		LabelJA:    label,
		SpriteBase: spriteBase,
		SeedStage:  true,
		SourcePath: sourcePath,
		TintHex:    tintHex,
		AccentHex:  accentHex,
	}
}

func shapeVariant(id string, speciesID string, label string, shape string, tintHex string, accentHex string) Variant {
	return Variant{
		ID:         id,
		SpeciesID:  speciesID,
		LabelEN:    label,
		LabelJA:    label,
		SpriteBase: id,
		SeedStage:  true,
		Shape:      shape,
		TintHex:    tintHex,
		AccentHex:  accentHex,
	}
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

func SpeciesByID(id string) (Species, bool) {
	for _, species := range SpeciesList {
		if species.ID == id {
			return species, true
		}
	}
	return Species{}, false
}

func ProfileForSpecies(id string) string {
	if species, ok := SpeciesByID(id); ok {
		return species.Profile
	}
	return ""
}

func WheelCapableSpecies(speciesID string) bool {
	switch ProfileForSpecies(speciesID) {
	case "degu", "small-mammal":
		return true
	default:
		return false
	}
}

func WheelCapableVariant(variant Variant) bool {
	return WheelCapableSpecies(variant.SpeciesID)
}

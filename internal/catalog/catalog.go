package catalog

type Species struct {
	ID      string
	Label   string
	Profile string
}

type Variant struct {
	ID             string
	SpeciesID      string
	BreedOrMorph   string
	Color          string
	PopularityTier int
	MotionProfile  string
	SourceStatus   string
	LabelEN        string
	LabelJA        string
	SpriteBase     string
	SeedStage      bool
	SourcePath     string
	Shape          string
	TintHex        string
	AccentHex      string
}

const (
	SourceStatusDeguMotion     = "degu_motion_source"
	SourceStatusPrototypeOnly  = "prototype_only"
	SourceStatusImageGenQueued = "imagegen_queued"

	MotionProfileDegu               = "degu"
	MotionProfileSmallRodentScurry  = "small-rodent-scurry"
	MotionProfileRabbitHop          = "rabbit-hop"
	MotionProfileDogTrot            = "dog-trot"
	MotionProfileCatStalk           = "cat-stalk"
	MotionProfileGeckoCrawl         = "gecko-crawl"
	MotionProfileTortoisePlod       = "tortoise-plod"
	MotionProfileFerretSlink        = "ferret-slink"
	MotionProfileGuineaPigWaddle    = "guinea-pig-waddle"
	MotionProfileHedgehogShuffle    = "hedgehog-shuffle"
	MotionProfileSquirrelBound      = "squirrel-bound"
	MotionProfileFoxTrot            = "fox-trot"
	MotionProfileRedPandaAmble      = "red-panda-amble"
	MotionProfileOtterSlide         = "otter-slide"
	MotionProfileSugarGliderSkitter = "sugar-glider-skitter"
	MotionProfileCapybaraLumber     = "capybara-lumber"
	MotionProfileSnakeSlither       = "snake-slither"
	MotionProfileDragonPlod         = "dragon-plod"
	MotionProfileFrogHop            = "frog-hop"
)

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
	{ID: "rat", Label: "Fancy rat", Profile: "small-mammal"},
	{ID: "mouse", Label: "Fancy mouse", Profile: "small-mammal"},
	{ID: "gerbil", Label: "Mongolian gerbil", Profile: "small-mammal"},
	{ID: "prairie_dog", Label: "Prairie dog", Profile: "small-mammal"},
	{ID: "chipmunk", Label: "Chipmunk", Profile: "small-mammal"},
	{ID: "bearded_dragon", Label: "Bearded dragon", Profile: "low-crawler"},
	{ID: "crested_gecko", Label: "Crested gecko", Profile: "low-crawler"},
	{ID: "corn_snake", Label: "Corn snake", Profile: "low-crawler"},
	{ID: "whites_tree_frog", Label: "White's tree frog", Profile: "hopper"},
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

	sourceVariantMeta("french_bulldog_fawn", "dog", "French Bulldog - fawn", "french_bulldog_fawn", srcDog, "c49a6c", "3a3026", "French Bulldog", "fawn", 1),
	sourceVariantMeta("labrador_yellow", "dog", "Labrador Retriever - yellow", "labrador_yellow", srcDog, "d7b46f", "", "Labrador Retriever", "yellow", 1),
	sourceVariantMeta("labrador_black", "dog", "Labrador Retriever - black", "labrador_black", srcDog, "2b2926", "", "Labrador Retriever", "black", 1),
	sourceVariantMeta("golden_retriever_golden", "dog", "Golden Retriever - golden", "golden_retriever_golden", srcDog, "c8873d", "ead9ad", "Golden Retriever", "golden", 1),
	sourceVariantMeta("german_shepherd_black_tan", "dog", "German Shepherd - black and tan", "german_shepherd_black_tan", srcDog, "3a3026", "b8874c", "German Shepherd Dog", "black and tan", 1),
	sourceVariantMeta("dachshund_red", "dog", "Dachshund - red", "dachshund_red", srcDog, "a8572b", "", "Dachshund", "red", 1),
	sourceVariantMeta("poodle_white", "dog", "Poodle - white", "poodle_white", srcDog, "eee7d8", "", "Poodle", "white", 1),
	sourceVariantMeta("beagle_tricolor", "dog", "Beagle - tricolor", "beagle_tricolor", srcDog, "7b4d2e", "f1e6d0", "Beagle", "tricolor", 1),
	sourceVariantMeta("bulldog_white_brindle", "dog", "Bulldog - white brindle", "bulldog_white_brindle", srcDog, "e7dfcf", "75513a", "Bulldog", "white brindle", 1),
	sourceVariantMeta("shiba_inu_red", "dog", "Shiba Inu - red", "shiba_inu_red", srcDog, "b7612f", "efe0c5", "Shiba Inu", "red", 2),
	sourceVariantMeta("pomeranian_orange", "dog", "Pomeranian - orange", "pomeranian_orange", srcDog, "d58233", "f2d4a2", "Pomeranian", "orange", 2),
	sourceVariantMeta("corgi_sable", "dog", "Corgi - sable", "corgi_sable", srcDog, "9b6236", "ead6b0", "Pembroke Welsh Corgi", "sable", 2),

	sourceVariantMeta("maine_coon_brown_tabby", "cat", "Maine Coon - brown tabby", "maine_coon_brown_tabby", srcCat, "6e5239", "a88457", "Maine Coon", "brown tabby", 1),
	sourceVariantMeta("ragdoll_seal_bicolor", "cat", "Ragdoll - seal bicolor", "ragdoll_seal_bicolor", srcCat, "e7ddca", "594333", "Ragdoll", "seal bicolor", 1),
	sourceVariantMeta("persian_white", "cat", "Persian - white", "persian_white", srcCat, "eee8dd", "", "Persian", "white", 1),
	sourceVariantMeta("british_shorthair_blue", "cat", "British Shorthair - blue", "british_shorthair_blue", srcCat, "737d84", "", "British Shorthair", "blue", 1),
	sourceVariantMeta("siamese_seal_point", "cat", "Siamese - seal point", "siamese_seal_point", srcCat, "d8c8aa", "4a3428", "Siamese", "seal point", 1),
	sourceVariantMeta("sphynx_pink", "cat", "Sphynx - pink", "sphynx_pink", srcCat, "d7a894", "", "Sphynx", "pink", 2),
	sourceVariantMeta("scottish_fold_silver_tabby", "cat", "Scottish Fold - silver tabby", "scottish_fold_silver_tabby", srcCat, "a7aaa5", "5e625f", "Scottish Fold", "silver tabby", 2),
	sourceVariantMeta("bengal_rosetted", "cat", "Bengal - rosetted", "bengal_rosetted", srcCat, "c69048", "3f2b20", "Bengal", "rosetted", 2),
	sourceVariantMeta("domestic_shorthair_calico", "cat", "Domestic Shorthair - calico", "domestic_shorthair_calico", srcCat, "e8dfcd", "bd6d2e", "Domestic Shorthair", "calico", 2),
	sourceVariantMeta("domestic_shorthair_tuxedo", "cat", "Domestic Shorthair - tuxedo", "domestic_shorthair_tuxedo", srcCat, "242322", "eee8dd", "Domestic Shorthair", "tuxedo", 2),

	sourceVariantMeta("holland_lop_broken_orange", "rabbit", "Holland Lop - broken orange", "holland_lop_broken_orange", srcRabbit, "e6d8c3", "c97833", "Holland Lop", "broken orange", 1),
	sourceVariantMeta("netherland_dwarf_chestnut", "rabbit", "Netherland Dwarf - chestnut", "netherland_dwarf_chestnut", srcRabbit, "8c633d", "d3b17a", "Netherland Dwarf", "chestnut", 1),
	sourceVariantMeta("mini_rex_black_otter", "rabbit", "Mini Rex - black otter", "mini_rex_black_otter", srcRabbit, "2f2c29", "b89865", "Mini Rex", "black otter", 1),
	sourceVariantMeta("lionhead_tort", "rabbit", "Lionhead - tort", "lionhead_tort", srcRabbit, "a66a3a", "4f3326", "Lionhead", "tort", 2),
	sourceVariantMeta("dutch_black_white", "rabbit", "Dutch - black and white", "dutch_black_white", srcRabbit, "2c2a28", "eee7dc", "Dutch", "black and white", 2),

	shapeVariantMeta("fancy_rat_hooded", "rat", "Fancy rat - hooded", "small_rodent", "efe5d4", "3b332c", "Fancy rat", "hooded", 2),
	shapeVariantMeta("fancy_mouse_white", "mouse", "Fancy mouse - white", "small_rodent", "eee7d8", "c99a8c", "Fancy mouse", "white", 2),
	shapeVariantMeta("mongolian_gerbil_agouti", "gerbil", "Mongolian gerbil - agouti", "small_rodent", "9d7448", "e5c997", "Mongolian gerbil", "agouti", 2),
	shapeVariantMeta("prairie_dog_tan", "prairie_dog", "Prairie dog - tan", "prairie_dog", "b98958", "e0c08e", "Prairie dog", "tan", 3),
	shapeVariantMeta("chipmunk_striped", "chipmunk", "Chipmunk - striped", "chipmunk", "a06a3a", "2f241e", "Chipmunk", "striped", 3),

	shapeVariantMeta("bearded_dragon_citrus", "bearded_dragon", "Bearded dragon - citrus", "dragon", "d99b37", "f0d268", "Bearded dragon", "citrus", 2),
	sourceVariantMeta("crested_gecko_harlequin", "crested_gecko", "Crested gecko - harlequin", "crested_gecko_harlequin", srcGecko, "b77b3d", "eee0bd", "Crested gecko", "harlequin", 2),
	shapeVariantMeta("corn_snake_amelanistic", "corn_snake", "Corn snake - amelanistic", "snake", "d66b36", "f0d096", "Corn snake", "amelanistic", 2),
	shapeVariantMeta("whites_tree_frog_green", "whites_tree_frog", "White's tree frog - green", "frog", "77a95a", "d7e9b8", "White's tree frog", "green", 2),
}

var runtimeVariantIDs = []string{
	"chinchilla_standard_gray",
}

func RuntimeVariants() []Variant {
	out := make([]Variant, 0, len(runtimeVariantIDs))
	for _, id := range runtimeVariantIDs {
		if variant, ok := VariantByID(id); ok {
			out = append(out, variant)
		}
	}
	return out
}

func VariantByID(id string) (Variant, bool) {
	for _, variant := range Variants {
		if variant.ID == id {
			return variant, true
		}
	}
	return Variant{}, false
}

func deguVariant(id string, label string, spriteBase string) Variant {
	return Variant{
		ID:             id,
		SpeciesID:      "degu",
		BreedOrMorph:   "Degu",
		Color:          label,
		PopularityTier: 1,
		MotionProfile:  MotionProfileDegu,
		SourceStatus:   SourceStatusDeguMotion,
		LabelEN:        "Degu - " + label,
		LabelJA:        "Degu - " + label,
		SpriteBase:     spriteBase,
	}
}

func sourceVariant(id string, speciesID string, label string, spriteBase string, sourcePath string, tintHex string, accentHex string) Variant {
	return Variant{
		ID:             id,
		SpeciesID:      speciesID,
		BreedOrMorph:   label,
		Color:          label,
		PopularityTier: 3,
		MotionProfile:  DefaultMotionProfileForSpecies(speciesID),
		SourceStatus:   SourceStatusPrototypeOnly,
		LabelEN:        label,
		LabelJA:        label,
		SpriteBase:     spriteBase,
		SeedStage:      true,
		SourcePath:     sourcePath,
		TintHex:        tintHex,
		AccentHex:      accentHex,
	}
}

func sourceVariantMeta(id string, speciesID string, label string, spriteBase string, sourcePath string, tintHex string, accentHex string, breedOrMorph string, color string, popularityTier int) Variant {
	variant := sourceVariant(id, speciesID, label, spriteBase, sourcePath, tintHex, accentHex)
	variant.BreedOrMorph = breedOrMorph
	variant.Color = color
	variant.PopularityTier = popularityTier
	return variant
}

func shapeVariant(id string, speciesID string, label string, shape string, tintHex string, accentHex string) Variant {
	return Variant{
		ID:             id,
		SpeciesID:      speciesID,
		BreedOrMorph:   label,
		Color:          label,
		PopularityTier: 3,
		MotionProfile:  DefaultMotionProfileForSpecies(speciesID),
		SourceStatus:   SourceStatusPrototypeOnly,
		LabelEN:        label,
		LabelJA:        label,
		SpriteBase:     id,
		SeedStage:      true,
		Shape:          shape,
		TintHex:        tintHex,
		AccentHex:      accentHex,
	}
}

func shapeVariantMeta(id string, speciesID string, label string, shape string, tintHex string, accentHex string, breedOrMorph string, color string, popularityTier int) Variant {
	variant := shapeVariant(id, speciesID, label, shape, tintHex, accentHex)
	variant.BreedOrMorph = breedOrMorph
	variant.Color = color
	variant.PopularityTier = popularityTier
	return variant
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

func DefaultMotionProfileForSpecies(speciesID string) string {
	switch speciesID {
	case "degu":
		return MotionProfileDegu
	case "rabbit":
		return MotionProfileRabbitHop
	case "dog":
		return MotionProfileDogTrot
	case "cat":
		return MotionProfileCatStalk
	case "gecko", "crested_gecko":
		return MotionProfileGeckoCrawl
	case "tortoise":
		return MotionProfileTortoisePlod
	case "ferret":
		return MotionProfileFerretSlink
	case "guinea_pig":
		return MotionProfileGuineaPigWaddle
	case "hedgehog":
		return MotionProfileHedgehogShuffle
	case "squirrel", "chipmunk":
		return MotionProfileSquirrelBound
	case "fox":
		return MotionProfileFoxTrot
	case "red_panda":
		return MotionProfileRedPandaAmble
	case "otter":
		return MotionProfileOtterSlide
	case "sugar_glider":
		return MotionProfileSugarGliderSkitter
	case "capybara":
		return MotionProfileCapybaraLumber
	case "corn_snake":
		return MotionProfileSnakeSlither
	case "bearded_dragon":
		return MotionProfileDragonPlod
	case "whites_tree_frog":
		return MotionProfileFrogHop
	default:
		return MotionProfileSmallRodentScurry
	}
}

func MotionProfileForVariant(variant Variant) string {
	if variant.MotionProfile != "" {
		return variant.MotionProfile
	}
	return DefaultMotionProfileForSpecies(variant.SpeciesID)
}

func WheelCapableSpecies(speciesID string) bool {
	return WheelCapableMotionProfile(DefaultMotionProfileForSpecies(speciesID))
}

func WheelCapableMotionProfile(profile string) bool {
	switch profile {
	case MotionProfileDegu, MotionProfileSmallRodentScurry:
		return true
	default:
		return false
	}
}

func WheelCapableVariant(variant Variant) bool {
	return WheelCapableMotionProfile(MotionProfileForVariant(variant))
}

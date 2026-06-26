#!/usr/bin/env python3
"""Build public-page animal icons and the hero preview from runtime sheets.

This is intentionally page-only. It does not modify runtime sprite orientation
or motion sheets. Use PAGE_RIGHT_FACING_VARIANTS for user-facing page assets
that should face right even if the underlying runtime frame faces left.
"""

from __future__ import annotations

from pathlib import Path

from PIL import Image, ImageDraw, ImageFilter


ROOT = Path(__file__).resolve().parents[1]
FRAME_W = 96
FRAME_H = 64
ICON_SIZE = 192
PREVIEW_W = 957
PREVIEW_H = 516

CURRENT_VARIANTS = [
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
]

PAGE_RIGHT_FACING_VARIANTS = {
    "sugar_glider_gray",
}

ICON_FRAMES = {
    "macaroni_mouse_tan": 4,
    "sugar_glider_gray": 4,
    "gecko_gray_brown": 4,
    "yorkshire_terrier_longcoat": 4,
}

PREVIEW_POSES = {
    "chinchilla_standard_gray": 4,
    "hamster_golden_syrian": 6,
    "djungarian_hamster": 8,
    "campbell_hamster": 10,
    "macaroni_mouse_tan": 12,
    "sugar_glider_gray": 4,
    "rabbit_chestnut_agouti": 16,
    "holland_lop_broken_orange": 18,
    "netherland_dwarf_chestnut": 20,
    "himalayan_rabbit": 22,
    "gecko_gray_brown": 6,
    "guinea_pig_tricolor": 10,
    "fancy_rat_hooded": 12,
    "albino_chipmunk": 16,
    "richardsons_ground_squirrel": 20,
    "yorkshire_terrier_longcoat": 14,
}


def frame_for_variant(variant: str, frame: int) -> Image.Image:
    sheet_path = ROOT / "assets" / "sprites" / f"{variant}_set00.png"
    if not sheet_path.exists():
        raise SystemExit(f"missing runtime sheet: {sheet_path}")
    sheet = Image.open(sheet_path).convert("RGBA")
    if sheet.size[1] != FRAME_H or sheet.size[0] < FRAME_W * (frame + 1):
        raise SystemExit(f"unexpected sheet size for {variant}: {sheet.size}")
    sprite = sheet.crop((frame * FRAME_W, 0, (frame + 1) * FRAME_W, FRAME_H))
    if variant in PAGE_RIGHT_FACING_VARIANTS:
        sprite = sprite.transpose(Image.Transpose.FLIP_LEFT_RIGHT)
    return sprite


def alpha_bbox(img: Image.Image) -> tuple[int, int, int, int]:
    bbox = img.getchannel("A").getbbox()
    if bbox is None:
        raise SystemExit("empty sprite frame")
    return bbox


def trim(img: Image.Image, pad: int = 4) -> Image.Image:
    left, top, right, bottom = alpha_bbox(img)
    left = max(0, left - pad)
    top = max(0, top - pad)
    right = min(img.width, right + pad)
    bottom = min(img.height, bottom + pad)
    return img.crop((left, top, right, bottom))


def scaled(sprite: Image.Image, max_w: int, max_h: int) -> Image.Image:
    sprite = trim(sprite)
    scale = min(max_w / sprite.width, max_h / sprite.height)
    w = max(1, round(sprite.width * scale))
    h = max(1, round(sprite.height * scale))
    return sprite.resize((w, h), Image.Resampling.LANCZOS)


def write_icons() -> None:
    out_dir = ROOT / "docs" / "assets" / "animal-icons"
    out_dir.mkdir(parents=True, exist_ok=True)
    for variant in CURRENT_VARIANTS:
        sprite = frame_for_variant(variant, ICON_FRAMES.get(variant, 0))
        animal = scaled(sprite, 150, 118)
        icon = Image.new("RGBA", (ICON_SIZE, ICON_SIZE), (0, 0, 0, 0))
        x = (ICON_SIZE - animal.width) // 2
        y = (ICON_SIZE - animal.height) // 2 + 8
        icon.alpha_composite(animal, (x, y))
        icon.save(out_dir / f"current-{variant}.png")


def draw_soft_shadow(canvas: Image.Image, cx: int, cy: int, w: int, h: int) -> None:
    shadow = Image.new("RGBA", canvas.size, (0, 0, 0, 0))
    d = ImageDraw.Draw(shadow)
    d.ellipse((cx - w // 2, cy - h // 2, cx + w // 2, cy + h // 2), fill=(25, 28, 30, 48))
    shadow = shadow.filter(ImageFilter.GaussianBlur(7))
    canvas.alpha_composite(shadow)


def write_preview() -> None:
    canvas = Image.new("RGBA", (PREVIEW_W, PREVIEW_H), (252, 253, 250, 255))
    d = ImageDraw.Draw(canvas)
    d.rectangle((0, 0, PREVIEW_W, PREVIEW_H), fill=(252, 253, 250, 255))
    d.rectangle((0, 318, PREVIEW_W, 458), fill=(246, 248, 244, 255))
    d.rectangle((0, 458, PREVIEW_W, PREVIEW_H), fill=(236, 240, 244, 255))
    d.rectangle((0, 458, PREVIEW_W, 461), fill=(210, 218, 224, 255))
    for x in range(28, 292, 36):
        d.rounded_rectangle((x, 474, x + 24, 498), radius=6, fill=(222, 228, 233, 255))
    d.rounded_rectangle((720, 475, 918, 498), radius=10, fill=(247, 249, 250, 255))

    positions = [
        (42, 362, 1.08),
        (112, 384, 1.05),
        (178, 344, 1.02),
        (246, 386, 1.02),
        (314, 352, 1.05),
        (382, 388, 1.03),
        (450, 344, 1.05),
        (520, 386, 1.0),
        (586, 352, 1.02),
        (654, 388, 1.02),
        (720, 352, 1.05),
        (790, 386, 1.02),
        (858, 350, 1.0),
        (72, 424, 1.0),
        (640, 424, 1.0),
        (820, 424, 0.98),
    ]
    for variant, (x, base_y, scale_factor) in zip(CURRENT_VARIANTS, positions):
        sprite = frame_for_variant(variant, PREVIEW_POSES.get(variant, 0))
        animal = scaled(sprite, round(90 * scale_factor), round(62 * scale_factor))
        draw_soft_shadow(canvas, x + animal.width // 2, base_y + 4, max(34, animal.width - 18), 10)
        canvas.alpha_composite(animal, (x, base_y - animal.height))

    canvas.convert("RGB").save(ROOT / "docs" / "assets" / "animalsdesktop-preview.png")


def main() -> None:
    write_icons()
    write_preview()
    print(f"wrote {len(CURRENT_VARIANTS)} page icons and preview")


if __name__ == "__main__":
    main()

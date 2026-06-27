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
UPCOMING_W = 256
UPCOMING_H = 160

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

UPCOMING_SILHOUETTES = [
    "leopard_gecko",
    "whites_tree_frog",
    "chipmunk",
    "budgerigar",
    "cockatiel",
    "java_sparrow",
    "lovebird",
    "parrotlet",
    "mixed_cat",
    "toy_poodle",
    "chihuahua",
]


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
        (368, 424, 0.96),
        (500, 424, 0.98),
        (610, 424, 0.98),
        (710, 424, 0.98),
        (820, 424, 0.98),
    ]
    for variant, (x, base_y, scale_factor) in zip(CURRENT_VARIANTS, positions):
        sprite = frame_for_variant(variant, PREVIEW_POSES.get(variant, 0))
        animal = scaled(sprite, round(90 * scale_factor), round(62 * scale_factor))
        draw_soft_shadow(canvas, x + animal.width // 2, base_y + 4, max(34, animal.width - 18), 10)
        canvas.alpha_composite(animal, (x, base_y - animal.height))

    canvas.convert("RGB").save(ROOT / "docs" / "assets" / "animalsdesktop-preview.png")


def _draw_rotated_ellipse(
    layer: Image.Image,
    box: tuple[int, int, int, int],
    angle: float,
    fill: tuple[int, int, int, int],
) -> None:
    w = box[2] - box[0]
    h = box[3] - box[1]
    tmp = Image.new("RGBA", (w + 20, h + 20), (0, 0, 0, 0))
    d = ImageDraw.Draw(tmp)
    d.ellipse((10, 10, 10 + w, 10 + h), fill=fill)
    tmp = tmp.rotate(angle, resample=Image.Resampling.BICUBIC, expand=True)
    x = box[0] - (tmp.width - w) // 2
    y = box[1] - (tmp.height - h) // 2
    layer.alpha_composite(tmp, (x, y))


def _draw_bird(
    d: ImageDraw.ImageDraw,
    *,
    body: tuple[int, int, int, int],
    head: tuple[int, int, int, int],
    beak: tuple[tuple[int, int], tuple[int, int], tuple[int, int]],
    tail: tuple[tuple[int, int], ...],
    legs: tuple[tuple[int, int, int, int], ...] = (),
    crest: tuple[tuple[int, int], ...] | None = None,
    wing: tuple[int, int, int, int] | None = None,
) -> None:
    ink = (28, 33, 30, 235)
    d.polygon(tail, fill=ink)
    d.ellipse(body, fill=ink)
    if wing is not None:
        d.ellipse(wing, fill=(18, 22, 20, 110))
    d.ellipse(head, fill=ink)
    d.polygon(beak, fill=ink)
    if crest is not None:
        d.polygon(crest, fill=ink)
    for line in legs:
        d.line(line, fill=ink, width=5)


def _draw_upcoming_silhouette(kind: str) -> Image.Image:
    scale = 2
    canvas = Image.new("RGBA", (UPCOMING_W * scale, UPCOMING_H * scale), (0, 0, 0, 0))
    d = ImageDraw.Draw(canvas)
    ink = (28, 33, 30, 235)

    if kind == "leopard_gecko":
        d.ellipse((122, 132, 334, 214), fill=ink)
        d.ellipse((304, 104, 390, 166), fill=ink)
        d.polygon(((382, 126), (426, 142), (384, 154)), fill=ink)
        d.line((128, 174, 58, 204), fill=ink, width=30)
        d.line((168, 200, 122, 258), fill=ink, width=14)
        d.line((226, 204, 204, 266), fill=ink, width=14)
        d.line((288, 202, 316, 264), fill=ink, width=14)
        d.line((344, 178, 388, 236), fill=ink, width=14)
    elif kind == "whites_tree_frog":
        d.ellipse((130, 96, 326, 238), fill=ink)
        d.ellipse((258, 62, 378, 160), fill=ink)
        d.ellipse((282, 48, 326, 92), fill=ink)
        d.ellipse((338, 50, 382, 96), fill=ink)
        d.line((168, 198, 102, 252), fill=ink, width=24)
        d.line((292, 204, 352, 266), fill=ink, width=24)
        d.line((300, 152, 392, 186), fill=ink, width=18)
        d.line((178, 156, 112, 188), fill=ink, width=18)
    elif kind == "chipmunk":
        d.ellipse((126, 124, 322, 220), fill=ink)
        d.ellipse((292, 82, 386, 156), fill=ink)
        d.polygon(((380, 112), (424, 126), (382, 140)), fill=ink)
        d.ellipse((312, 76, 340, 102), fill=ink)
        d.line((130, 154, 66, 108, 40, 74), fill=ink, width=22)
        d.line((170, 204, 150, 264), fill=ink, width=12)
        d.line((246, 208, 244, 268), fill=ink, width=12)
        d.line((314, 196, 346, 252), fill=ink, width=12)
    elif kind == "budgerigar":
        _draw_bird(
            d,
            body=(142, 104, 300, 232),
            head=(252, 64, 338, 144),
            beak=((334, 100), (370, 114), (334, 128)),
            tail=((150, 178), (58, 236), (134, 220)),
            legs=((220, 222, 208, 260), (258, 222, 270, 260)),
            wing=(170, 132, 268, 222),
        )
    elif kind == "cockatiel":
        _draw_bird(
            d,
            body=(132, 106, 310, 238),
            head=(254, 58, 352, 152),
            beak=((348, 104), (388, 120), (348, 136)),
            tail=((144, 184), (34, 262), (134, 232)),
            legs=((222, 230, 208, 268), (270, 230, 284, 268)),
            crest=((294, 68), (318, 12), (326, 74), (352, 28), (338, 86)),
            wing=(162, 132, 276, 228),
        )
    elif kind == "java_sparrow":
        _draw_bird(
            d,
            body=(140, 116, 302, 230),
            head=(258, 82, 340, 158),
            beak=((336, 118), (370, 130), (336, 142)),
            tail=((150, 182), (70, 222), (146, 224)),
            legs=((220, 224, 214, 262), (262, 224, 270, 262)),
            wing=(174, 140, 270, 218),
        )
    elif kind == "lovebird":
        _draw_bird(
            d,
            body=(142, 98, 310, 236),
            head=(252, 58, 356, 158),
            beak=((350, 112), (392, 130), (350, 148)),
            tail=((150, 184), (86, 230), (150, 226)),
            legs=((224, 226, 214, 264), (268, 226, 280, 264)),
            wing=(172, 126, 278, 222),
        )
    elif kind == "parrotlet":
        _draw_bird(
            d,
            body=(152, 108, 298, 228),
            head=(254, 72, 336, 150),
            beak=((332, 114), (368, 128), (332, 142)),
            tail=((156, 184), (92, 228), (154, 220)),
            legs=((224, 222, 216, 258), (260, 222, 270, 258)),
            wing=(178, 136, 264, 216),
        )
    elif kind == "mixed_cat":
        d.ellipse((130, 128, 328, 226), fill=ink)
        d.ellipse((300, 88, 398, 168), fill=ink)
        d.polygon(((318, 98), (336, 48), (352, 110)), fill=ink)
        d.polygon(((366, 98), (386, 52), (388, 118)), fill=ink)
        d.polygon(((392, 124), (438, 140), (394, 152)), fill=ink)
        d.line((138, 154, 72, 90, 58, 40), fill=ink, width=22)
        d.line((170, 208, 152, 270), fill=ink, width=15)
        d.line((238, 216, 226, 270), fill=ink, width=15)
        d.line((296, 212, 322, 270), fill=ink, width=15)
        d.line((342, 198, 370, 260), fill=ink, width=15)
    elif kind == "toy_poodle":
        for box in [
            (146, 110, 274, 224),
            (242, 76, 348, 174),
            (286, 42, 360, 112),
            (222, 76, 282, 140),
            (318, 98, 380, 166),
            (130, 126, 190, 190),
            (176, 202, 224, 274),
            (264, 202, 312, 274),
        ]:
            d.ellipse(box, fill=ink)
        d.polygon(((344, 114), (386, 130), (346, 146)), fill=ink)
        d.line((146, 150, 82, 92), fill=ink, width=16)
    elif kind == "chihuahua":
        d.ellipse((142, 132, 318, 224), fill=ink)
        d.ellipse((298, 76, 386, 158), fill=ink)
        d.polygon(((306, 90), (292, 24), (342, 84)), fill=ink)
        d.polygon(((356, 88), (400, 24), (390, 106)), fill=ink)
        d.polygon(((382, 116), (430, 130), (386, 144)), fill=ink)
        d.line((150, 152, 92, 112), fill=ink, width=14)
        d.line((178, 206, 166, 270), fill=ink, width=14)
        d.line((278, 208, 294, 270), fill=ink, width=14)
        d.line((330, 190, 356, 260), fill=ink, width=14)
    else:
        raise ValueError(f"unknown upcoming silhouette: {kind}")

    bbox = canvas.getchannel("A").getbbox()
    if bbox is None:
        raise SystemExit(f"empty silhouette: {kind}")
    return canvas.resize((UPCOMING_W, UPCOMING_H), Image.Resampling.LANCZOS)


def write_upcoming_silhouettes() -> None:
    out_dir = ROOT / "docs" / "assets" / "upcoming-silhouettes"
    out_dir.mkdir(parents=True, exist_ok=True)
    silhouettes: list[Image.Image] = []
    for kind in UPCOMING_SILHOUETTES:
        img = _draw_upcoming_silhouette(kind)
        img.save(out_dir / f"{kind}.png")
        silhouettes.append(img)

    cols = 4
    rows = (len(silhouettes) + cols - 1) // cols
    composite = Image.new("RGBA", (UPCOMING_W * cols, UPCOMING_H * rows), (0, 0, 0, 0))
    for i, img in enumerate(silhouettes):
        x = (i % 4) * UPCOMING_W
        y = (i // 4) * UPCOMING_H
        composite.alpha_composite(img, (x, y))
    composite.save(ROOT / "docs" / "assets" / "animalsdesktop-coming-soon-silhouettes.png")


def main() -> None:
    write_icons()
    write_preview()
    write_upcoming_silhouettes()
    print(
        f"wrote {len(CURRENT_VARIANTS)} page icons, preview, "
        f"and {len(UPCOMING_SILHOUETTES)} upcoming silhouettes"
    )


if __name__ == "__main__":
    main()

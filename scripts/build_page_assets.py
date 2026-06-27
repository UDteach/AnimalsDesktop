#!/usr/bin/env python3
"""Build public-page animal icons and the hero preview from runtime sheets.

This is intentionally page-only. It does not modify runtime sprite orientation
or motion sheets. Use PAGE_RIGHT_FACING_VARIANTS for user-facing page assets
that should face right even if the underlying runtime frame faces left.
"""

from __future__ import annotations

from collections import deque
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
UPCOMING_SOURCE = (
    ROOT
    / "docs"
    / "art-source"
    / "one-frame-method-fullrun-20260627"
    / "page-coming-soon"
    / "coming-soon-fifteen-animals-imagegen-source.png"
)
UPCOMING_SOURCE_COLS = 5
UPCOMING_SOURCE_ROWS = 3
SILHOUETTE_INK = (14, 18, 15, 255)

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
    "gecko_leopard": 4,
    "whites_tree_frog_blue": 19,
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
    "gecko_leopard": 6,
    "whites_tree_frog_blue": 19,
    "guinea_pig_tricolor": 10,
    "fancy_rat_hooded": 12,
    "chipmunk_striped": 16,
    "albino_chipmunk": 16,
    "richardsons_ground_squirrel": 20,
    "yorkshire_terrier_longcoat": 14,
}

# The page-specific ImageGen source sheet was generated in this cell order.
# Keep this separate so the public display order can follow product priority.
UPCOMING_SOURCE_LAYOUT = [
    "leopard_gecko",
    "whites_tree_frog",
    "chipmunk",
    "budgerigar",
    "cockatiel",
    "java_sparrow",
    "lovebird",
    "parrotlet",
    "scottish_fold",
    "mixed_cat",
    "munchkin",
    "ragdoll",
    "minuet",
    "toy_poodle",
    "chihuahua",
]

UPCOMING_SILHOUETTES = [
    "chipmunk",
    "whites_tree_frog",
    "leopard_gecko",
    "budgerigar",
    "cockatiel",
    "java_sparrow",
    "lovebird",
    "parrotlet",
    "scottish_fold",
    "mixed_cat",
    "munchkin",
    "ragdoll",
    "minuet",
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


def _crop_upcoming_source_cell(source: Image.Image, index: int) -> Image.Image:
    col = index % UPCOMING_SOURCE_COLS
    row = index // UPCOMING_SOURCE_COLS
    left = round(source.width * col / UPCOMING_SOURCE_COLS)
    top = round(source.height * row / UPCOMING_SOURCE_ROWS)
    right = round(source.width * (col + 1) / UPCOMING_SOURCE_COLS)
    bottom = round(source.height * (row + 1) / UPCOMING_SOURCE_ROWS)
    return source.crop((left, top, right, bottom))


def _background_like(pixel: tuple[int, int, int, int]) -> bool:
    r, g, b, a = pixel
    if a <= 8:
        return True
    return min(r, g, b) >= 236 and max(r, g, b) - min(r, g, b) <= 24


def _foreground_mask(cell: Image.Image) -> Image.Image:
    rgba = cell.convert("RGBA")
    w, h = rgba.size
    pixels = rgba.load()
    visited = bytearray(w * h)
    queue: deque[tuple[int, int]] = deque()

    def offset(x: int, y: int) -> int:
        return y * w + x

    def enqueue(x: int, y: int) -> None:
        i = offset(x, y)
        if visited[i] or not _background_like(pixels[x, y]):
            return
        visited[i] = 1
        queue.append((x, y))

    for x in range(w):
        enqueue(x, 0)
        enqueue(x, h - 1)
    for y in range(h):
        enqueue(0, y)
        enqueue(w - 1, y)

    while queue:
        x, y = queue.popleft()
        for nx, ny in ((x - 1, y), (x + 1, y), (x, y - 1), (x, y + 1)):
            if 0 <= nx < w and 0 <= ny < h:
                enqueue(nx, ny)

    mask = Image.new("L", (w, h), 0)
    mask_pixels = mask.load()
    for y in range(h):
        for x in range(w):
            r, g, b, a = pixels[x, y]
            if a > 8 and not visited[offset(x, y)]:
                mask_pixels[x, y] = 255

    return _largest_mask_component(mask.filter(ImageFilter.MaxFilter(3)).filter(ImageFilter.MinFilter(3)))


def _largest_mask_component(mask: Image.Image) -> Image.Image:
    mask = mask.convert("L")
    w, h = mask.size
    pixels = mask.load()
    visited = bytearray(w * h)
    best: list[tuple[int, int]] = []

    def offset(x: int, y: int) -> int:
        return y * w + x

    for y in range(h):
        for x in range(w):
            start = offset(x, y)
            if visited[start] or pixels[x, y] == 0:
                continue
            visited[start] = 1
            queue: deque[tuple[int, int]] = deque([(x, y)])
            component: list[tuple[int, int]] = []
            while queue:
                cx, cy = queue.popleft()
                component.append((cx, cy))
                for nx, ny in ((cx - 1, cy), (cx + 1, cy), (cx, cy - 1), (cx, cy + 1)):
                    if not (0 <= nx < w and 0 <= ny < h):
                        continue
                    i = offset(nx, ny)
                    if visited[i] or pixels[nx, ny] == 0:
                        continue
                    visited[i] = 1
                    queue.append((nx, ny))
            if len(component) > len(best):
                best = component

    if not best:
        return mask
    out = Image.new("L", mask.size, 0)
    out_pixels = out.load()
    for x, y in best:
        out_pixels[x, y] = pixels[x, y]
    return out


def _silhouette_from_source_cell(cell: Image.Image, kind: str) -> Image.Image:
    mask = _foreground_mask(cell)
    bbox = mask.getbbox()
    if bbox is None:
        raise SystemExit(f"empty ImageGen source cell for upcoming silhouette: {kind}")

    cropped_mask = mask.crop(bbox)
    silhouette = Image.new("RGBA", cropped_mask.size, SILHOUETTE_INK)
    silhouette.putalpha(cropped_mask)
    scale = min(
        (UPCOMING_W - 28) / silhouette.width,
        (UPCOMING_H - 24) / silhouette.height,
    )
    resized = silhouette.resize(
        (max(1, round(silhouette.width * scale)), max(1, round(silhouette.height * scale))),
        Image.Resampling.LANCZOS,
    )

    canvas = Image.new("RGBA", (UPCOMING_W, UPCOMING_H), (0, 0, 0, 0))
    x = (UPCOMING_W - resized.width) // 2
    y = (UPCOMING_H - resized.height) // 2
    canvas.alpha_composite(resized, (x, y))
    return canvas


def write_upcoming_silhouettes() -> None:
    if not UPCOMING_SOURCE.exists():
        raise SystemExit(f"missing page-specific ImageGen source: {UPCOMING_SOURCE}")
    source = Image.open(UPCOMING_SOURCE).convert("RGBA")
    source_slots = UPCOMING_SOURCE_COLS * UPCOMING_SOURCE_ROWS
    if len(UPCOMING_SOURCE_LAYOUT) > source_slots:
        raise SystemExit(
            f"too many source layout silhouettes ({len(UPCOMING_SOURCE_LAYOUT)}) for "
            f"{UPCOMING_SOURCE_COLS}x{UPCOMING_SOURCE_ROWS} source sheet"
        )

    out_dir = ROOT / "docs" / "assets" / "upcoming-silhouettes"
    out_dir.mkdir(parents=True, exist_ok=True)
    expected = {f"{kind}.png" for kind in UPCOMING_SILHOUETTES}
    for sidecar in out_dir.glob("._*.png"):
        sidecar.unlink()
    for stale in out_dir.glob("*.png"):
        if stale.name not in expected:
            stale.unlink()
    silhouettes: list[Image.Image] = []
    source_indices = {kind: i for i, kind in enumerate(UPCOMING_SOURCE_LAYOUT)}
    for kind in UPCOMING_SILHOUETTES:
        if kind not in source_indices:
            raise SystemExit(f"missing upcoming silhouette source cell for {kind}")
        cell = _crop_upcoming_source_cell(source, source_indices[kind])
        img = _silhouette_from_source_cell(cell, kind)
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
    for sidecar in out_dir.glob("._*.png"):
        sidecar.unlink()


def delete_appledouble_sidecars(root: Path) -> None:
    for sidecar in root.rglob("._*"):
        if sidecar.is_file():
            sidecar.unlink()


def main() -> None:
    write_icons()
    write_preview()
    write_upcoming_silhouettes()
    delete_appledouble_sidecars(ROOT / "docs" / "assets")
    print(
        f"wrote {len(CURRENT_VARIANTS)} page icons, preview, "
        f"and {len(UPCOMING_SILHOUETTES)} upcoming silhouettes"
    )


if __name__ == "__main__":
    main()

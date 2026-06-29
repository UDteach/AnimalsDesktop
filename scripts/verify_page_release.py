#!/usr/bin/env python3
"""Verify release-sensitive fields on the public Pages HTML."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
INDEX = ROOT / "docs" / "index.html"
WINDOWS_AMD64_ASSET = "AnimalsDesktop-windows-amd64.zip"
WINDOWS_386_ASSET = "AnimalsDesktop-windows-386.zip"
MAC_ARM64_ASSET = "AnimalsDesktop-macos-arm64.zip"
MAC_AMD64_ASSET = "AnimalsDesktop-macos-amd64.zip"
CHECKSUM_ASSET = "SHA256SUMS.txt"
CATALOG = ROOT / "internal" / "catalog" / "catalog.go"
EXPECTED_RELEASE = "v0.2.5"
EXPECTED_UPCOMING = [
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
    "lionhead_rabbit",
    "african_fat_tailed_gecko",
    "shoebill",
]


def fail(message: str) -> None:
    print(f"verify_page_release: {message}", file=sys.stderr)
    raise SystemExit(1)


def one(pattern: str, html: str, label: str) -> str:
    matches = re.findall(pattern, html, flags=re.S)
    if not matches:
        fail(f"missing {label}")
    if len(matches) > 1:
        fail(f"expected one {label}, found {len(matches)}")
    match = matches[0]
    if isinstance(match, tuple):
        return next(part for part in match if part)
    return match


def release_tag(asset: str, html: str, label: str) -> str:
    return one(
        rf'releases/download/(v[0-9][^"/]*)/{re.escape(asset)}',
        html,
        label,
    )

def runtime_variant_ids() -> list[str]:
    catalog = CATALOG.read_text(encoding="utf-8")
    block = one(
        r"var runtimeVariantIDs = \[\]string\{(.*?)\}",
        catalog,
        "runtime variant list",
    )
    ids = re.findall(r'"([^"]+)"', block)
    if not ids:
        fail("runtime variant list is empty")
    return ids


def current_page_variant_ids(html: str) -> list[str]:
    block = one(
        r'<div class="current-grid">(.*?)</div>\s*</section>',
        html,
        "current animal grid",
    )
    ids = re.findall(r'assets/animal-icons/current-([a-z0-9_]+)\.png', block)
    if not ids:
        fail("current animal grid has no icons")
    return ids


def upcoming_page_ids(html: str) -> list[str]:
    block = one(
        r'<div class="future-grid">(.*?)</div>\s*</section>',
        html,
        "upcoming animal grid",
    )
    ids = re.findall(r'assets/upcoming-silhouettes/([a-z0-9_]+)\.png', block)
    if not ids:
        fail("upcoming animal grid has no silhouettes")
    return ids


def verify_asset_refs(html: str) -> None:
    refs = sorted(set(re.findall(r'(?:src|href)="(assets/[^"]+)"', html)))
    if not refs:
        fail("no local asset references found")
    missing = [ref for ref in refs if not (ROOT / "docs" / ref).exists()]
    if missing:
        fail(f"missing local page assets: {missing}")


def main() -> None:
    html = INDEX.read_text(encoding="utf-8")

    if "releases/download/v0.2.0" in html:
        fail("blocked public Windows v0.2.0 release link remains")

    windows_tag = release_tag(WINDOWS_AMD64_ASSET, html, "Windows amd64 download version tag")
    windows_386_tag = release_tag(WINDOWS_386_ASSET, html, "Windows 386 download version tag")
    checksum_tag = release_tag(CHECKSUM_ASSET, html, "SHA256SUMS download version tag")
    mac_arm64_tag = release_tag(MAC_ARM64_ASSET, html, "macOS arm64 download version tag")
    mac_amd64_tag = release_tag(MAC_AMD64_ASSET, html, "macOS amd64 download version tag")
    if windows_tag != EXPECTED_RELEASE:
        fail(f"Windows amd64 download tag {windows_tag} does not match expected {EXPECTED_RELEASE}")

    windows_badge = one(
        r"<span[^>]*>\s*Windows版\s*<strong>(v[0-9][^<]*)</strong>\s*</span>",
        html,
        "visible Windows version badge",
    )
    release_badge = one(
        r"<strong data-release-version>(v[0-9][^<]*)</strong>",
        html,
        "public release version badge",
    )

    for label, got in {
        "Windows 386 download": windows_386_tag,
        "SHA256SUMS download": checksum_tag,
        "Windows badge": windows_badge,
        "release badge": release_badge,
    }.items():
        if got != windows_tag:
            fail(f"{label} {got} does not match Windows amd64 tag {windows_tag}")

    if mac_arm64_tag != mac_amd64_tag:
        fail(f"macOS download tags differ: {mac_arm64_tag} != {mac_amd64_tag}")

    runtime_ids = runtime_variant_ids()
    page_ids = current_page_variant_ids(html)
    if page_ids != runtime_ids:
        fail(f"current animal grid {page_ids} does not match runtime variants {runtime_ids}")

    upcoming_ids = upcoming_page_ids(html)
    if upcoming_ids != EXPECTED_UPCOMING:
        fail(f"upcoming animal grid {upcoming_ids} does not match expected priority {EXPECTED_UPCOMING}")

    for required in ("ライオンラビット", "ハシビロコウ", "低モーション", "lionhead", "shoebill", "low-motion"):
        if required not in html:
            fail(f"missing future roadmap text: {required}")
    for required in (
        'data-i18n="versions.v025.title"',
        "v0.2.5 / 2026-06-29",
        "v0.2.5 / June 29, 2026",
        'data-i18n="versions.v024.title"',
        "v0.2.4 / 2026-06-28",
        "v0.2.4 / June 28, 2026",
        "v0.2.3 / 2026-06-28",
        "v0.2.3 / June 28, 2026",
    ):
        if required not in html:
            fail(f"missing version history text: {required}")

    verify_asset_refs(html)

    print(
        f"Pages release links verified: Windows {windows_tag}, macOS {mac_arm64_tag}, "
        f"current animals {len(page_ids)}, upcoming silhouettes {len(upcoming_ids)}"
    )


if __name__ == "__main__":
    main()

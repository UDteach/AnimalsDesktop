#!/usr/bin/env python3
"""Verify release-sensitive fields on the public Pages HTML."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
INDEX = ROOT / "docs" / "index.html"
WINDOWS_AMD64_ASSET = "AnimalsDesktop-windows-amd64.zip"
WINDOWS_AMD64_NO_NETWORK_ASSET = "AnimalsDesktop-windows-amd64-no-network.zip"
WINDOWS_386_ASSET = "AnimalsDesktop-windows-386.zip"
MAC_ARM64_ASSET = "AnimalsDesktop-macos-arm64.zip"
MAC_AMD64_ASSET = "AnimalsDesktop-macos-amd64.zip"
CHECKSUM_ASSET = "SHA256SUMS.txt"
CATALOG = ROOT / "internal" / "catalog" / "catalog.go"
EXPECTED_RELEASE = "v0.2.15"
EXPECTED_UPCOMING: list[str] = []


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
    matches = re.findall(
        rf'releases/download/(v[0-9][^"/]*)/{re.escape(asset)}',
        html,
        flags=re.S,
    )
    if not matches:
        fail(f"missing {label}")
    tags = sorted(set(matches))
    if len(tags) > 1:
        fail(f"{label} points to multiple versions: {', '.join(tags)}")
    return tags[0]


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
        r"const animalCatalog = \[(.*?)\n\s*\];",
        html,
        "JavaScript animal catalog",
    )
    ids = re.findall(r'\["([a-z0-9_]+)",\s*"[^"]+"\]', block)
    if not ids:
        fail("JavaScript animal catalog has no IDs")
    return ids


def upcoming_page_ids(html: str) -> list[str]:
    matches = re.findall(
        r'<div class="future-grid">(.*?)</div>\s*</section>',
        html,
        flags=re.S,
    )
    if not matches:
        return []
    if len(matches) > 1:
        fail(f"expected one upcoming animal grid, found {len(matches)}")
    block = matches[0]
    ids = re.findall(r'assets/upcoming-silhouettes/([a-z0-9_]+)\.png', block)
    return ids


def verify_asset_refs(html: str, page_ids: list[str]) -> None:
    refs = set(re.findall(r'(?:src|href)="(assets/[^"]+)"', html))
    refs.update(re.findall(r'url\(["\']?(assets/[^"\')]+)', html))
    refs.update(f"assets/animal-icons/current-{variant_id}.png" for variant_id in page_ids)
    if not refs:
        fail("no local asset references found")
    normalized = sorted({re.split(r"[?#]", ref, maxsplit=1)[0] for ref in refs})
    missing = [ref for ref in normalized if not (ROOT / "docs" / ref).exists()]
    if missing:
        fail(f"missing local page assets: {missing}")


def main() -> None:
    html = INDEX.read_text(encoding="utf-8")

    if "releases/download/v0.2.0" in html:
        fail("blocked public Windows v0.2.0 release link remains")

    windows_tag = release_tag(WINDOWS_AMD64_ASSET, html, "Windows amd64 download version tag")
    windows_no_network_tag = release_tag(
        WINDOWS_AMD64_NO_NETWORK_ASSET,
        html,
        "Windows amd64 no-network download version tag",
    )
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
        "Windows no-network download": windows_no_network_tag,
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

    for required in (
        '<time datetime="2026-07-04">2026-07-04</time>',
        '<time datetime="2026-07-04">July 4, 2026</time>',
        "固定表示とランダム表示",
        "種類で絞り込め",
        "fixed animals and random slots",
        "filtered by animal type",
        "v0.2.14 / 2026-07-03",
        "v0.2.14 / July 3, 2026",
        "ハシビロコウが別の動物として復元",
        "stable IDs",
        "groups the Windows and Mac animal pickers by type",
        "v0.2.13 / 2026-07-03",
        "v0.2.13 / July 3, 2026",
        "roster additions do not remap saved choices",
        "v0.2.12 / 2026-07-03",
        "v0.2.12 / July 3, 2026",
        "AnimalsDesktop-windows-amd64-no-network.zip",
        "Smart App Control FAQ",
        "SmartScreen reputation",
        "GlobalSignのSmart App Control事例",
        "GlobalSign Smart App Control case study",
        "v0.2.11 / 2026-07-02",
        "v0.2.11 / July 2, 2026",
        "v0.2.10 / 2026-07-02",
        "v0.2.10 / July 2, 2026",
        "v0.2.9 / 2026-07-02",
        "v0.2.9 / July 2, 2026",
        "v0.2.8 / 2026-07-01",
        "v0.2.8 / July 1, 2026",
        "v0.2.7 / 2026-07-01",
        "v0.2.7 / July 1, 2026",
        "v0.2.6 / 2026-06-30",
        "v0.2.6 / June 30, 2026",
        "v0.2.5 / 2026-06-29",
        "v0.2.5 / June 29, 2026",
        "v0.2.4 / 2026-06-28",
        "v0.2.4 / June 28, 2026",
        "v0.2.3 / 2026-06-28",
        "v0.2.3 / June 28, 2026",
        "v0.2.2 / 2026-06-27",
        "v0.2.2 / June 27, 2026",
        "v0.2.1 / 2026-06-27",
        "v0.2.1 / June 27, 2026",
        "v0.1.5 / 2026-06-21",
        "v0.1.5 / June 21, 2026",
        "Pagesに残っていた候補12件は v0.2.9 で追加済みです。次の追加候補は準備中です。",
        "オカメインコ",
        "Cockatiel - normal gray",
        "v0.2.15の詳細",
        "About v0.2.15",
        "最大10匹",
        "Show up to 10 animals",
        "木、種、餌",
        "trees, seeds, and food",
    ):
        if required not in html:
            fail(f"missing version history text: {required}")

    verify_asset_refs(html, page_ids)

    print(
        f"Pages release links verified: Windows {windows_tag}, macOS {mac_arm64_tag}, "
        f"current animals {len(page_ids)}, upcoming silhouettes {len(upcoming_ids)}"
    )


if __name__ == "__main__":
    main()

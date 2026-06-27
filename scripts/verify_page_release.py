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


def main() -> None:
    html = INDEX.read_text(encoding="utf-8")

    if "releases/download/v0.2.0" in html:
        fail("blocked public Windows v0.2.0 release link remains")

    windows_tag = release_tag(WINDOWS_AMD64_ASSET, html, "Windows amd64 download version tag")
    windows_386_tag = release_tag(WINDOWS_386_ASSET, html, "Windows 386 download version tag")
    checksum_tag = release_tag(CHECKSUM_ASSET, html, "SHA256SUMS download version tag")
    mac_arm64_tag = release_tag(MAC_ARM64_ASSET, html, "macOS arm64 download version tag")
    mac_amd64_tag = release_tag(MAC_AMD64_ASSET, html, "macOS amd64 download version tag")

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

    print(f"Pages release links verified: Windows {windows_tag}, macOS {mac_arm64_tag}")


if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""Verify release-sensitive fields on the public Pages HTML."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
INDEX = ROOT / "docs" / "index.html"
MAC_ARM64_ASSET = "AnimalsDesktop-macos-arm64.zip"
MAC_AMD64_ASSET = "AnimalsDesktop-macos-amd64.zip"


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


def main() -> None:
    html = INDEX.read_text(encoding="utf-8")

    for term in [
        "AnimalsDesktop-windows",
        "releases/download/v0.2.0",
        "SHA256SUMS.txt",
    ]:
        if term in html:
            fail(f"blocked public Windows v0.2.0 release link remains: {term}")

    mac_arm64_tag = one(
        rf'releases/download/(v[0-9][^"/]*)/{re.escape(MAC_ARM64_ASSET)}',
        html,
        "macOS arm64 download version tag",
    )
    mac_amd64_tag = one(
        rf'releases/download/(v[0-9][^"/]*)/{re.escape(MAC_AMD64_ASSET)}',
        html,
        "macOS amd64 download version tag",
    )
    if mac_arm64_tag != mac_amd64_tag:
        fail(f"macOS download tags differ: {mac_arm64_tag} != {mac_amd64_tag}")

    print(f"Pages release links verified: macOS {mac_arm64_tag}, Windows v0.2.0 hidden")


if __name__ == "__main__":
    main()

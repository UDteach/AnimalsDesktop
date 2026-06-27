#!/usr/bin/env python3
"""Verify release-sensitive fields on the public Pages HTML."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
INDEX = ROOT / "docs" / "index.html"
WINDOWS_ASSET = "AnimalsDesktop-windows-amd64.zip"


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

    if "Windows版 Windows 10/11" not in html:
        fail("missing Windows download label")
    if WINDOWS_ASSET not in html:
        fail(f"missing {WINDOWS_ASSET}")

    windows_tag = one(
        rf'releases/download/(v[0-9][^"/]*)/{re.escape(WINDOWS_ASSET)}',
        html,
        "Windows download version tag",
    )
    windows_badge = one(
        r"<span>\s*Windows版\s*<strong>(v[0-9][^<]*)</strong>\s*</span>",
        html,
        "visible Windows version badge",
    )
    release_badge = one(
        r"<strong data-release-version>(v[0-9][^<]*)</strong>",
        html,
        "public release version badge",
    )

    if windows_badge != windows_tag:
        fail(f"Windows badge {windows_badge} does not match download tag {windows_tag}")
    if release_badge != windows_tag:
        fail(f"release badge {release_badge} does not match Windows download tag {windows_tag}")

    print(f"Windows page release verified: {windows_tag}")


if __name__ == "__main__":
    main()

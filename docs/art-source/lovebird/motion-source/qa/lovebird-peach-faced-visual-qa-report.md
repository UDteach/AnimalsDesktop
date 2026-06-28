# Lovebird Peach-Faced Set00 Visual QA

Date: 2026-06-28

Variant: `lovebird_peach_faced`

Verdict: accepted as an asset-only `set00` motion source.

Mechanical QA:

- `go run ./cmd/auditframes -frames-dir docs/art-source/lovebird/motion-source/accepted-frames/set00 -strict -artifact-warnings -report docs/art-source/lovebird/motion-source/accepted-frames/set00-auditframes-report.json`
- Result: `valid=62 missing=0 invalid=0 warnings=92`
- `go run ./cmd/assemblemotion -frames-dir docs/art-source/lovebird/motion-source/accepted-frames/set00 -out docs/art-source/lovebird/motion-source/sheets/lovebird-peach-faced-source-set00.png -report docs/art-source/lovebird/motion-source/accepted-frames/set00-assemblemotion-report.json`

Visual review:

- `53-61`, `45-61`, and full `00-61` checker contacts were reviewed before promotion.
- The set preserves a peach/orange face, green body, visible feet/contact, and stable baseline.
- The final `53-61` reaction/recover frames keep a consistent `47px` height and baseline `56`.
- Earlier one-frame size jumps around `26` and `43` were resolved before final acceptance.
- No text, scenery, props, duplicate animal, cropped head/tail/feet, or obvious background contamination was found in the promoted frames.

Scope:

- Asset/source-frame promotion only.
- Runtime variants, GitHub Pages, release notes, tags, downloads, and deploy settings were not changed for this promotion.

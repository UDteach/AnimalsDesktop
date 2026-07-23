# Adding an Animal

This is the canonical integration path for a new AnimalsDesktop species,
breed, morph, or coat. It keeps source-art acceptance, local runtime preview,
public Pages, and release publication as separate decisions.

## Prerequisites

- Go 1.24 or newer
- Python 3.10 or newer
- `python -m pip install -r requirements-dev.txt`

## Boundaries

- A new species or materially different silhouette needs its own ImageGen
  source family. Do not create it by recoloring another species.
- A coat variant may share a species only after that species' silhouette and
  motion set are stable. Accepted production frames still come from one
  ImageGen call per frame.
- Source-art acceptance does not automatically make an animal selectable.
- Runtime selection does not automatically update Pages, tags, downloads, or a
  GitHub Release. Those remain explicit release work.
- `cmd/importsheet` is the legacy Degu compatibility importer. Use
  `cmd/importanimals` for current animal families.

## 1. Prepare accepted source art

Create exactly 62 transparent `96x64` frames:

```text
docs/art-source/<family>/motion-source/accepted-frames/set00/
  frame-00.png
  ...
  frame-61.png
```

Run the mechanical gate and assemble the source sheet:

```powershell
go run ./cmd/auditframes `
  -frames-dir docs\art-source\<family>\motion-source\accepted-frames\set00 `
  -strict -artifact-warnings -motion-warnings

go run ./cmd/assemblemotion `
  -frames-dir docs\art-source\<family>\motion-source\accepted-frames\set00 `
  -out docs\art-source\<family>\motion-source\sheets\<sprite-base>-source-set00.png `
  -report docs\art-source\<family>\motion-source\qa\assemble-set00.json
```

Mechanical success is not visual acceptance. Review light, dark, and checker
contacts for species read, anatomy, baseline, crop, matte, and continuity
before registering the source as accepted.

## 2. Register the catalog entry

Update `internal/catalog/catalog.go`:

1. Add `SpeciesList`, picker-group, and default-motion-profile entries only
   when the species is new.
2. Add the source frame and `set00` sheet path constants.
3. Add one `acceptedMotionVariant(...)` entry with stable snake-case ID and
   `SpriteBase`, English and Japanese labels, breed/morph, color, popularity
   tier, and motion profile.
4. Do not add the ID to `runtimeVariantIDs` until the animal is approved for
   the selectable runtime roster.

Every importer run calls `catalog.Validate()`. It rejects duplicate or malformed
IDs, duplicate sprite bases, unknown species/profiles/statuses, unsafe paths,
incomplete metadata, invalid motion-source combinations, and invalid runtime
roster entries before assets are written.

## 3. Check only the new animal

Use the targeted check before writing runtime assets:

```powershell
go run ./cmd/importanimals -variant <variant_id> -check
go run ./cmd/validatemotion -variant <variant_id> -require-accepted
go test -buildvcs=false ./internal/catalog ./cmd/importanimals
```

`-check` runs the existing import pipeline in a temporary directory. It checks
the catalog entry, source PNG, motion-sheet dimensions, alpha/content, and
output generation without changing repository files.

## 4. Import only the new animal

```powershell
go run ./cmd/importanimals -variant <variant_id>
```

Targeted mode writes only that animal's normalized generated source and ten
runtime sprite sheets. It intentionally does not replace the aggregate
`seed-import-report.json` or `animalsdesktop-seed-preview.png` unless explicit
`-report` or `-preview` paths are supplied.

Confirm the targeted diff contains only the intended generated source and
`<sprite-base>_set00.png` through `<sprite-base>_set09.png`.

A single `set00` source sheet is allowed as a clearly warned local preview
fallback. Once any `set01`-`set09` source is added, the family must contain all
ten sheets; a partial family is rejected instead of silently duplicating
`set00`.

## 5. Opt into runtime and release gates

After parent visual approval, add the variant ID once to `runtimeVariantIDs` at
the intended picker position. The Go catalog, Windows runtime, and page-asset
builder all consume that curated order; do not copy a second roster into tests
or scripts.

Run the full importer only at the integration/release boundary:

```powershell
go run ./cmd/importanimals
go run ./cmd/validatemotion -runtime-only -require-accepted
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go build -buildvcs=false -ldflags="-H=windowsgui" `
  -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

The no-flag importer remains the authority for the aggregate report and
preview. Run it twice and confirm the second run produces no asset diff.

Pages, version text, release workflows, tags, downloadable ZIPs, and GitHub
Releases are separate explicit-publication work. Do not infer publication
permission from local animal integration.

## Addition checklist

- [ ] New species uses its own source family; coat-only reuse is justified.
- [ ] `62/62` one-frame source files pass mechanical and parent visual QA.
- [ ] Catalog metadata and paths pass `catalog.Validate()`.
- [ ] Targeted `importanimals -check` passes without repository writes.
- [ ] Targeted import changes only the intended animal assets.
- [ ] Runtime roster inclusion is an explicit parent decision.
- [ ] Full import is deterministic.
- [ ] Runtime validation, tests, vet, Windows build, and whitespace checks pass.
- [ ] Pages/release state is unchanged unless separately authorized.

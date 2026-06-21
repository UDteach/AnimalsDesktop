# AnimalsDesktop

AnimalsDesktop is a Windows desktop pet app currently preparing its first animal-complete release.

Public page: <https://udteach.github.io/AnimalsDesktop/>

## Current Status

The project is in Coming Soon mode.

- The first formal pet target is chinchilla.
- Unfinished animals and prototype assets are not listed as released pets.
- The next release tag waits until chinchilla has full motion coverage, settings/runtime selection, local QA, and GitHub workflow success.
- After chinchilla, the same completion loop will move to the next animal family.

## Release Gate

A pet is release-ready only when it has:

- accepted ImageGen source art
- 96x64 transparent sprite coverage
- 10 motion sets x 62 frames
- settings and runtime selection
- local tests and build checks
- updated public page
- `go run ./cmd/validatemotion -runtime-only -require-accepted` passing

## Development

Useful checks during local development:

```powershell
go run ./cmd/importsheet
go run ./cmd/importanimals
go run ./cmd/auditframes -root docs\art-source\chinchilla\motion-source\accepted-frames
go run ./cmd/validatemotion -variant chinchilla_standard_gray
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/degu
git diff --check
```

Run `cmd/prepareframe` only on one-pose candidates, outside the standard QA loop. It rejects checker/noisy backgrounds; prepared output still needs visual review before it counts as accepted art:

```powershell
go run ./cmd/prepareframe -src path\to\candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00.png
```

Run `cmd/assemblemotion` only for a set after `cmd/auditframes -strict` passes for that set:

```powershell
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -strict
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-set00-draft.png
```

In-progress sets are expected to fail strict mode until all 62 standalone transparent frames exist.

Do not create a release tag until the current animal target is complete.

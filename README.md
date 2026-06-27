# AnimalsDesktop

AnimalsDesktop is a small desktop pet app for Windows and macOS. Chinchilla,
hamster, Djungarian hamster, Campbell hamster, macaroni mouse, sugar glider,
rabbit, Holland Lop, Netherland Dwarf, Himalayan rabbit, gecko, guinea pig,
fancy rat, albino chipmunk, Richardson's ground squirrel, and Yorkshire Terrier
sprites walk near the Windows taskbar or along the bottom edge above the Mac
Dock.

Public page: <https://udteach.github.io/AnimalsDesktop/>

## Current Status

`v0.1.5` is an early public test release for sixteen accepted initial
62-frame ImageGen motion animals:

- chinchilla standard gray
- golden Syrian hamster
- Djungarian hamster
- Campbell hamster
- tan macaroni mouse
- gray sugar glider
- chestnut agouti rabbit
- broken orange Holland Lop
- chestnut Netherland Dwarf
- Himalayan rabbit
- gray-brown gecko
- tricolor guinea pig
- hooded fancy rat
- albino chipmunk
- Richardson's ground squirrel
- longcoat Yorkshire Terrier

This is not the final full-motion completion gate. The test release is for
checking desktop behavior, scale, direction handling, and click interaction
before expanding each animal to the full 10-set motion contract.

## Runtime Scope

The v0.1.5 preview intentionally exposes only the sixteen accepted initial
runtime animals listed above. Unverified candidate species should not appear in
the runtime picker until their source art and motion behavior pass the release QA
loop. As animals graduate into a release, remove them from the future queue and
move them into the current-animal page section. Coming-soon silhouettes should be
page-specific generated art, not repurposed runtime/prototype images.

Future queue candidates use current popular-pet signals, then get verified per
animal before production starts. The current page queue starts with leopard
gecko, blue White's tree frog, and striped chipmunk, then prioritizes
budgerigar, cockatiel, and Java sparrow. The cat-breed page queue follows the
2026 iPet/Nyanpedia cat breed ranking top five: mixed cat, Scottish Fold,
Munchkin, Ragdoll, and Minuet. Re-check ranking sources before starting a cat
runtime lane.

Each animal promoted into the current runtime/page list should move the preview
version forward by a small patch bump, with page text, workflow checks, and
download artifacts kept in sync for that version.

The typing wheel is intentionally limited to chinchilla and hamster. Other
runtime animals continue to react to typing with movement, but do not enter the
wheel state.

Foraging props such as hay, twigs, seeds, and other small food/debris items are
disabled in v0.1.5 so the preview shows only the animal sprites and their core
movement.

## Release Gate

A full animal-family release is ready only when it has:

- accepted ImageGen source art
- 96x64 transparent sprite coverage
- 10 motion sets x 62 frames
- settings and runtime selection
- local tests and build checks
- updated public page
- `go run ./cmd/validatemotion -runtime-only -require-accepted` passing

`v0.1.5` is an explicit test-preview exception for the sixteen initial-motion
animals plus Mac distribution. Future full-content releases should
still satisfy the full gate unless a new preview exception is documented.

## Development

Useful checks during local development:

```powershell
go run ./cmd/importsheet
go run ./cmd/importanimals
go run ./cmd/auditframes -root docs\art-source\chinchilla\motion-source\accepted-frames
go run ./cmd/validatemotion -runtime-only -require-accepted
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

Windows release ZIPs are built on a Windows machine and can be uploaded to the
same release tag as the Mac ZIPs when they are ready:

```powershell
New-Item -ItemType Directory -Force dist | Out-Null
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -buildvcs=false -ldflags="-H=windowsgui -s -w -X main.appVersion=v0.1.5" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
Compress-Archive -Path dist\AnimalsDesktop.exe,README.md -DestinationPath dist\AnimalsDesktop-windows-amd64.zip -Force
```

macOS release ZIPs are built with:

```bash
VERSION=v0.1.5 GOARCH=arm64 scripts/build_macos.sh
VERSION=v0.1.5 GOARCH=amd64 scripts/build_macos.sh
```

Run `cmd/prepareframe` only on one-pose candidates, outside the standard QA loop.
It rejects checker/noisy backgrounds; prepared output still needs visual review
before it counts as accepted art:

```powershell
go run ./cmd/prepareframe -src path\to\candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\<set-id>\frame-00.png
go run ./cmd/prepareframe -background chroma-green -src path\to\green-candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\<set-id>\frame-00.png
```

Run `cmd/assemblemotion` only for a set after `cmd/auditframes -strict` passes:

```powershell
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\<set-id> -strict
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\<set-id> -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-<set-id>-draft.png
```

In-progress sets are expected to fail strict mode until all 62 standalone
transparent frames exist.

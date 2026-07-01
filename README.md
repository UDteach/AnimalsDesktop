# AnimalsDesktop

AnimalsDesktop is a small desktop pet app for Windows and macOS. The current
preview exposes 40 selectable animal sprites that walk near the Windows taskbar
or along the bottom edge above the Mac Dock.

Public page: <https://udteach.github.io/AnimalsDesktop/>

Current app version: `v0.2.7`

## Current Status

`v0.2.7` is an early public test release for 40 accepted 62-frame ImageGen
motion animals. It keeps the v0.2.6 runtime roster but temporarily removes
true albino chipmunk from the public picker until white-background contrast and
red-eye readability are repaired. It keeps the
original sixteen preview animals, the current GitHub Pages priority wave, and
the latest accepted asset lanes:

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
- black-eyed white chipmunk
- Richardson's ground squirrel
- longcoat Yorkshire Terrier
- striped chipmunk
- leopard gecko
- blue White's tree frog
- green-yellow budgerigar
- normal gray cockatiel
- normal Java sparrow
- green parrotlet
- blue-green parrotlet
- peach-faced lovebird
- seal bicolor ragdoll
- silver tabby Scottish Fold
- fawn French Bulldog
- brown tabby Maine Coon
- calico domestic shorthair
- blue British Shorthair
- apricot Toy Poodle
- brown tabby Munchkin
- Roborovski hamster
- Russian smoke white guinea pig
- quokka
- Miniature Schnauzer
- Japanese giant salamander
- white wagtail
- tabby-white domestic shorthair

This is not the final full-motion completion gate. The test release is for
checking desktop behavior, scale, direction handling, and click interaction
before expanding each animal to the full 10-set motion contract.

## Runtime Scope

The v0.2.7 preview intentionally exposes only the 40 accepted runtime animals
listed above. `true_albino_chipmunk` remains cataloged as accepted source
evidence but is excluded from runtime until its no-pattern repair lane passes
white-background visibility review. Unverified candidate species should not appear in
the runtime picker until their source art and motion behavior pass the release QA
loop. As animals graduate into a release, remove them from the future queue and
move them into the current-animal page section. Coming-soon silhouettes should be
page-specific generated art, not repurposed runtime/prototype images.

Future queue candidates use current popular-pet signals, then get verified per
animal before production starts. After the v0.2.7 release, the remaining Pages
queue is leucistic sugar glider, African dormouse, Netherland Dwarf Himalayan,
American flying squirrel, black-and-white long-haired hamster, yellow
Djungarian hamster, pearl white Djungarian hamster, fancy rat blue hooded,
fancy rat chocolate self, fancy rat cream agouti, gray rabbit, and African
fat-tailed gecko. The next newly requested lanes are a white lionhead-pattern
rabbit and a special low-motion shoebill.

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
animals plus Mac distribution. `v0.2.7` is an explicit test-preview exception
for the forty-animal runtime roster before the full 10-set motion gate.
Future full-content releases should still satisfy the full gate unless a new
preview exception is documented.

## Development

Useful checks during local development:

```powershell
go run ./cmd/importsheet
go run ./cmd/importanimals
go run ./cmd/auditframes -root docs\art-source\chinchilla\motion-source\accepted-frames -artifact-warnings -motion-warnings
go run ./cmd/validatemotion -runtime-only -require-accepted
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go run ./cmd/winresicon -src docs/assets/animalsdesktop-preview.png -out winres/icon.png
go run github.com/tc-hib/go-winres@v0.3.1 make --arch amd64 --out cmd/animalsdesktop/rsrc --file-version v0.2.7 --product-version v0.2.7
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
git diff --check
```

Windows release ZIPs are built on a Windows machine and can be uploaded to the
same release tag as the Mac ZIPs when they are ready:

```powershell
New-Item -ItemType Directory -Force dist | Out-Null
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -buildvcs=false -ldflags="-H=windowsgui -s -w -X main.appVersion=v0.2.7" -o dist\AnimalsDesktop.exe ./cmd/animalsdesktop
Compress-Archive -Path dist\AnimalsDesktop.exe,README.md -DestinationPath dist\AnimalsDesktop-windows-amd64.zip -Force
```

macOS release ZIPs are built with:

```bash
VERSION=v0.2.7 GOARCH=arm64 scripts/build_macos.sh
VERSION=v0.2.7 GOARCH=amd64 scripts/build_macos.sh
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
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\<set-id> -strict -artifact-warnings -motion-warnings
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\<set-id> -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-<set-id>-draft.png
```

In-progress sets are expected to fail strict mode until all 62 standalone
transparent frames exist.

## Windows Trust Notes

The Windows release workflow embeds file metadata, product metadata, a Windows 10+ manifest, and an app icon into `AnimalsDesktop.exe`. It also publishes `SHA256SUMS.txt` next to the release ZIPs and includes `SECURITY.txt` inside each ZIP with the expected EXE hash and false-positive submission notes.

For the best Microsoft Defender SmartScreen and McAfee outcome, release builds should be Authenticode-signed with a timestamped public-trust code-signing certificate. The preferred CI path is Microsoft Azure Artifact Signing because the signing key stays in Microsoft-managed HSMs instead of being stored as a repository secret.

Configure these GitHub Secrets for Azure Artifact Signing:

- `AZURE_CLIENT_ID`
- `AZURE_TENANT_ID`
- `AZURE_SUBSCRIPTION_ID`
- `AZURE_ARTIFACT_SIGNING_ENDPOINT`
- `AZURE_ARTIFACT_SIGNING_ACCOUNT_NAME`
- `AZURE_ARTIFACT_SIGNING_CERTIFICATE_PROFILE_NAME`

If you have a legacy or private `.pfx` signing certificate, the workflow can use these fallback secrets:

- `WINDOWS_CERTIFICATE_BASE64`: base64-encoded `.pfx`
- `WINDOWS_CERTIFICATE_PASSWORD`: `.pfx` password

If signing secrets are missing, the workflow still builds and publishes checksums, but the EXE remains unsigned and may continue to receive reputation-based warnings until Microsoft/McAfee reputation or allowlisting catches up.

`v0.2.0` is retained as a mistaken Windows prerelease. `v0.2.1` is the main-line Windows trust-hardening release. `v0.2.2` keeps that trust-hardening work and adds Mac parity for animal selection, size controls, language, and display settings. `v0.2.3` is a Windows settings UI hotfix with the same sixteen-animal scope. `v0.2.4` expands the selectable animal roster to 35. `v0.2.5` fixes Windows mixed-DPI multi-monitor overlay size and placement for the 35-animal roster. `v0.2.6` expands it to 41 and adds true albino chipmunk, Miniature Schnauzer, Japanese giant salamander, white wagtail, tabby-white cat, and blue-green parrotlet. `v0.2.7` temporarily removes true albino chipmunk from the public runtime while the no-pattern albino repair lane improves white-background readability, leaving 40 selectable animals and the same remaining Pages queue plus lionhead-pattern rabbit and low-motion shoebill as the next lanes.

Do not create a stable/final release tag until the current animal target is complete.

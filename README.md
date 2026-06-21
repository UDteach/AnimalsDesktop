# Animals Desktop

Animals Desktop is a Go desktop-pet app based on the DeguDesktop taskbar pet. It keeps the lightweight Win32 overlay, tray menu, keyboard reaction, random stroll, wheel, names, and foraging behavior, then expands the selectable art from degu coats into multiple real animal species.

Repository: <https://github.com/UDteach/AnimalsDesktop>

## Current Species

The runtime catalog currently includes:

- Degu: wild agouti, black, blue/slate gray, gray, white/cream, sand/champagne, chocolate, black pied, agouti pied, blue pied, cream pied
- Chinchilla: standard gray
- Macaroni mouse / fat-tailed gerbil: tan
- Rabbit: chestnut agouti
- Small dog: cream and tan
- Cat: brown tabby
- Gecko: gray brown
- Hamster: golden Syrian

Non-degu species are seed-stage assets generated from source-truth still images. They are selectable in the app and have deterministic runtime sheets, but they are not yet full species-specific motion sets.

![Seed animal preview](docs/assets/animalsdesktop-seed-preview.png)

## Windows Features

- Transparent always-on-top pet layer above the Windows taskbar
- Tray menu and Japanese/English settings window
- 1-10 visible pets
- Fixed, per-pet, or random animal/color selection
- Optional per-pet names with hover labels
- Keyboard reaction and random stroll modes
- Typing wheel behavior
- Foraging, carrying, eating, digging, gnawing, and grooming behavior
- GitHub Release based update check and installer path

The executable entrypoint is still `./cmd/degu` while the codebase is being migrated from DeguDesktop naming.

## Asset Pipeline

Degu motion assets use the original frame importer:

```powershell
go run ./cmd/importsheet
```

It reads `assets/source/frames/wild_agouti`, coat guides, forage art, and wheel/icon sources, then writes `assets/sprites/degu_*.png`, `assets/tray.ico`, `docs/assets/degu-preview.png`, and `assets/source/import-report.json`.

Seed animal assets use:

```powershell
go run ./cmd/importanimals
```

It reads source-truth images recorded under `docs/art-source`, `docs/art-intake`, and `docs/source-truth`, then writes:

- `assets/sprites/<animal>_set00.png` through `set09.png`
- `assets/source/animals/seed-import-report.json`
- `docs/assets/animalsdesktop-seed-preview.png`

The shared runtime registry is in `internal/catalog`.

## Development

```powershell
go run ./cmd/importsheet
go run ./cmd/importanimals
go test -buildvcs=false ./...
go vet -buildvcs=false ./...
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/degu
```

Run the app locally:

```powershell
go run ./cmd/degu
```

## Release

Push a `v*` tag to build Windows release ZIPs via GitHub Actions. The Windows app checks `UDteach/AnimalsDesktop` Releases for the latest matching architecture zip.

Expected release assets:

- `AnimalsDesktop-windows-amd64.zip`
- `AnimalsDesktop-windows-386.zip`

macOS packaging scripts are still present from the baseline and have been renamed to AnimalsDesktop, but current multi-species validation is Windows-first.

## Cloudflare Pages

`wrangler.jsonc` points at `docs/` for a static Pages output directory. The default GitHub workflow builds a `docs/` artifact but does not enable or deploy GitHub Pages automatically. Pages or Cloudflare deployment should be enabled separately when the repository visibility and hosting target are settled.

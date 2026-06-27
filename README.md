# AnimalsDesktop

AnimalsDesktop is a Windows desktop pet app currently preparing its first animal-complete release.

Public page: <https://udteach.github.io/AnimalsDesktop/>

Current Windows app version: `v0.2.0`

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
go run ./cmd/winresicon -src docs/assets/degu-preview.png -out winres/icon.png
go run github.com/tc-hib/go-winres@v0.3.1 make --arch amd64 --out cmd/degu/rsrc --file-version v0.2.0 --product-version v0.2.0
go build -buildvcs=false -ldflags="-H=windowsgui" -o dist\AnimalsDesktop.exe ./cmd/degu
git diff --check
```

Run `cmd/prepareframe` only on one-pose candidates, outside the standard QA loop. It rejects checker/noisy backgrounds; prepared output still needs visual review before it counts as accepted art:

```powershell
go run ./cmd/prepareframe -src path\to\candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00.png
go run ./cmd/prepareframe -background chroma-green -src path\to\green-candidate.png -out docs\art-source\chinchilla\motion-source\prepared-candidates\set00\frame-00.png
```

Run `cmd/assemblemotion` only for a set after `cmd/auditframes -strict` passes for that set:

```powershell
go run ./cmd/auditframes -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -strict
go run ./cmd/assemblemotion -frames-dir docs\art-source\chinchilla\motion-source\accepted-frames\set00 -out docs\art-source\chinchilla\motion-source\sheets\chinchilla-standard-gray-source-set00-draft.png
```

In-progress sets are expected to fail strict mode until all 62 standalone transparent frames exist.

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

`v0.2.0` is treated as a Windows preview release for runtime testing, release packaging, and trust hardening. It is marked as a GitHub prerelease because the current runtime animals still use draft one-set motion sources.

Do not create a stable/final release tag until the current animal target is complete.

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

## Development

Useful checks during local development:

```powershell
go test -buildvcs=false ./...
git diff --check
```

Do not create a release tag until the current animal target is complete.

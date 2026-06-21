# Chinchilla Motion Source

Status: in_progress
Generated: 2026-06-21

## Source Batch

- Tool: ImageGen
- Thread image folder: `C:\Users\user\.codex\generated_images\019ee7c8-bfbb-7832-abcb-b8def1a72f27`
- Accepted extraction source: first 16-pose chinchilla sheet
- Rejected generation: second sheet contained a non-animal symbol and is not used

The accepted extraction source was generated with a checker background rather than a true transparent alpha channel. The raw checker source was kept outside the tracked tree under `.codex/tmp/chinchilla-imagegen-raw/`.

## Tracked Outputs

- `frames/chinchilla-standard-gray-pose-00.png` through `pose-15.png`
- `chinchilla-standard-gray-pose-frames-check.png`

Each extracted pose frame is:

- 96x64
- transparent PNG
- one chinchilla per frame
- generated from the same ImageGen pose batch

## Current Use

These frames are source material for the chinchilla motion pass. They are not yet a release-ready 62-frame runtime set.

Before release, this family still needs:

- 62 accepted frames for the selected chinchilla variant
- generated `set00` through `set09` runtime sheets
- importer integration
- visual QA on light and dark backgrounds
- local Go tests, vet, Windows build, and GitHub Actions success

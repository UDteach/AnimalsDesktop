# Chinchilla Motion Source

Status: in_progress
Generated: 2026-06-21

## Source Batch

- Tool: ImageGen
- Thread image folder: `C:\Users\user\.codex\generated_images\019ee7c8-bfbb-7832-abcb-b8def1a72f27`
- Draft extraction source: first 16-pose chinchilla sheet
- Rejected generation: second sheet contained a non-animal symbol and is not used

The draft extraction source was generated with a checker background rather than a true transparent alpha channel. The raw checker source was kept outside the tracked tree under `.codex/tmp/chinchilla-imagegen-raw/`.

## Tracked Outputs

- `frames/chinchilla-standard-gray-pose-00.png` through `pose-15.png`
- `chinchilla-standard-gray-pose-frames-check.png`
- `chinchilla-standard-gray-frame-map.csv`
- `sheets/chinchilla-standard-gray-source-set00-draft.png`
- `chinchilla-standard-gray-source-set00-draft-check.png`

Each extracted pose frame is:

- 96x64
- transparent PNG
- one chinchilla per frame
- generated from the same ImageGen pose batch

## Current Use

These frames are source material for the chinchilla motion pass. They are not yet a release-ready motion family.

The `source-set00-draft` sheet maps the 16 extracted poses into the 62-frame slot contract for review. `cmd/importanimals` can emit 10 runtime sheets from this one draft sheet, but this is only a temporary integration path.

This family is intentionally marked draft because it still reuses key poses and has not passed the accepted-motion bar. Mechanically shifted draft `set01` through `set09` sheets were rejected because they looked wrong and are not a substitute for accepted ImageGen motion sources.

Accepted replacement work must be generated as one pose per PNG. Do not promote a multi-pose sheet, grid, contact sheet, or baked checker-background image as source art, even if it can be cut into 96x64 frames. A batch may contain up to four separate single-pose PNGs for review, but each accepted source frame must be a standalone transparent PNG before it is assembled into the 62-frame runtime sheet.

Before release, this family still needs:

- 62 accepted frames with enough real pose variation for the selected chinchilla variant
- accepted ImageGen `set00` through `set09` motion source variation
- visual QA on light and dark backgrounds
- local Go tests, vet, Windows build, and GitHub Actions success

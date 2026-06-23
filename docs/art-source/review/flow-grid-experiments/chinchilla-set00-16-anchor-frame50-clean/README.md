# Chinchilla Flow Nano Banana Pro 16-cell test

Date: 2026-06-23

## Setup

- Service: Google Flow
- Model: Nano Banana Pro
- Output: square, x4, 2K download
- Inputs added to prompt:
  - `style-anchors.png`
  - `grid-seed.png`
- Prompt entry note: browser-level paste/type did not update Flow's internal prompt state. OS-level virtual key input activated the prompt and enabled generation.

## Outputs

- Original downloads:
  - `flow-nanobanana-pro-2ref-output-01.jpeg`
  - `flow-nanobanana-pro-2ref-output-02.jpeg`
  - `flow-nanobanana-pro-2ref-output-03.jpeg`
  - `flow-nanobanana-pro-2ref-output-04.jpeg`
- PNG conversions:
  - `flow-nanobanana-pro-2ref-output-01.png` through `-04.png`
- Green-normalized PNGs:
  - `flow-nanobanana-pro-2ref-output-01-greenfixed.png` through `-04-greenfixed.png`

## Parse Results

- Original `batch.json` fixed guide manifest: all candidates parsed `0/16`.
- Actual 512px grid manifest without green cleanup: all candidates rejected due boundary chroma noise.
- Actual 512px grid with green cleanup:
  - candidate 01: `11/16`
  - candidate 02: `9/16`
  - candidate 03: `7/16`
  - candidate 04: `9/16`

Remaining failures are mainly disconnected components and chroma pinholes. The generated sheets are visually consistent, but most cells duplicate the same pose, so this is not accepted as a motion-production workflow.

## Variant A: Explicit Limb Delta Prompt

Follow-up date: 2026-06-23

- Prompt file: `prompt-variant-a-short.txt`
- Prompt notes and next variants: `prompt-variants.md`
- Consultation packet: `failure-consultation-variant-a.md`
- Cache-recovered outputs:
  - `flow-nanobanana-pro-variant-a-cache-141219.jpeg`
  - `flow-nanobanana-pro-variant-a-cache-141429.jpeg`

Result:

- Pose variation improved versus the first Nano Banana Pro prompt.
- The visible `grid-seed.png` guide was still reproduced as cell lines, so raw output is rejected as `GUIDE_INK`.
- Diagnostic green normalization was used only to inspect animal content; it is not an acceptance cleanup path.
- `141219` green-normalized parse: `3/16`, mostly rejected for `MAT_PINHOLE` or detached alpha.
- `141429` green-normalized parse: `0/16`, rejected for boundary/grid residue.

Current conclusion:

Do not retry the same visible grid seed with only pose wording changes. The next prompt experiment should change the guide strategy: style-anchor-only invisible grid, a smaller 4-frame strip, or independent images.

## Review Images

- `flow-nanobanana-pro-2ref-contact.png`: four Flow outputs.
- `parsed-pro-01-greenfixed-contact.png`: best parsed candidate preview.
- `flow-nanobanana-pro-variant-a-cache-contact.png`: Variant A recovered outputs.
- `flow-nanobanana-pro-variant-a-141219-qa-contact.png`: Variant A diagnostic parse view.
- `flow-nanobanana-pro-variant-a-141429-qa-contact.png`: Variant A diagnostic parse view.

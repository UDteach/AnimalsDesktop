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

## Review Images

- `flow-nanobanana-pro-2ref-contact.png`: four Flow outputs.
- `parsed-pro-01-greenfixed-contact.png`: best parsed candidate preview.

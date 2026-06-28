# Budgerigar Green-Yellow Visual QA - 2026-06-27

Status: accepted `set00` source asset.

Accepted source paths:

- `docs/art-source/budgerigar/motion-source/accepted-frames/set00/`
- `docs/art-source/budgerigar/motion-source/sheets/budgerigar-green-yellow-source-set00.png`

Mechanical QA after promotion:

```text
audited 62 frame slots: valid=62 missing=0 invalid=0 warnings=50
```

Promoted review contacts:

- `docs/art-source/budgerigar/motion-source/contacts/budgerigar-green-yellow-set00-frame-55-61-fallback-contact.png`
- `docs/art-source/budgerigar/motion-source/contacts/budgerigar-green-yellow-set00-frame-38-43-fallback-contact.png`
- `docs/art-source/budgerigar/motion-source/contacts/budgerigar-green-yellow-set00-full-contact.png`

Accepted in this pass:

- Tail sequence `55-61` has visible feet/contact, stable enough size/baseline,
  and consistent yellow/green/blue coloring after the parent fallback for `60`.
- Eat/chew sequence `38-43` has visible feet/contact, stable enough
  size/baseline, and no remaining gray/white face drift after parent fallbacks
  for `40-41`.

Acceptance notes:

- `frame-60.png` is a local fallback from neighbors, documented in
  `qa/frame-60-parent-fallback-note.md`.
- `frame-40.png` and `frame-41.png` are local fallbacks from neighbors,
  documented in `qa/frame-40-41-parent-fallback-note.md`.
- These fallback frames are accepted for source-art continuity after parent
  contact-sheet review because they keep the yellow head, green body, visible
  feet/contact, baseline, and scale better than the rejected ImageGen retries.
- This pass does not add the variant to runtime, Pages, release notes, tags, or
  downloadable artifacts.

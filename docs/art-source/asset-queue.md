# AnimalsDesktop Asset Queue

This queue tracks asset production rules without exposing unverified animal
families to the public page or runtime picker. Promote or name a new family only
when the parent thread explicitly accepts it as the next production target.

## Current Released State

- Latest test-preview release: `v0.1.4`
- Preview scope: six accepted initial runtime animals with Mac ZIPs; Windows ZIP prepared on a Windows machine
- Current preview animals: chinchilla, hamster, macaroni mouse, sugar glider, rabbit, gecko
- Public page scope: the six preview animals above, plus a "coming soon" silhouette queue
- Coming-soon queue: guinea pig, fancy rat, albino chipmunk, Richardson's ground squirrel, Yorkshire terrier
- Next full animal-content release: hold until one accepted family reaches
  DeguDesktop-level 10-set asset coverage
- Wheel-capable runtime animals: chinchilla and hamster only

## Batch Strategy

One thread should produce more images, but for fewer animals:

- Source candidate thread: one named family only after parent approval.
- Motion thread: one species family, 1-3 accepted variants, 62 source frames per accepted variant.
- Integration thread: parent only; verifies alpha bounds, 96x64 readability, deterministic import, catalog coverage, and runtime behavior.
- Frame progress thread: one species family only; uses `cmd/auditframes` to track valid/missing/invalid standalone PNG frames without promoting unfinished or opaque ImageGen outputs.

This gives each thread enough output to compare candidates without pushing 600+
unreviewed frames through one context.

## Completion Contract

Use `docs/art-source/motion-contract.md` as the promotion and release gate. The
runtime target stays DeguDesktop-compatible: 10 motion sets, 62 frames per set,
96x64 transparent frames, and the DeguDesktop frame slots preserved. A family is
not complete just because it is selectable; it is complete for a version bump
only when its source-truth art, 62-frame motion set, importer output,
settings/tray reflection, visual review, and QA pass.

Work one family at a time for full releases. The `v0.1.4` preview is allowed to
publish six accepted initial-motion animals for runtime testing, but the next
full-content version should finish a concrete 10-set slice before claiming
DeguDesktop-level coverage. Public coming-soon silhouettes are release-managed:
when an animal ships, remove it from the queue image/list and move it into the
current-animal section.

## Profile-Specific Notes

- Chinchilla: heavy soft body, rounded ears, fluffy tail, gentle scurry, no degu recolor.
- Macaroni mouse / fat-tailed gerbil: compact small body, thick tail, tiny feet, low fast scurry.
- Hamster: round cheek-forward body, very short tail, cheek/groom actions, wheel capable.
- Rabbit: readable ears, hop contact points, stable baseline.
- Sugar glider / momonga-style runtime family: membrane silhouette, low skitter, no long airborne scene.

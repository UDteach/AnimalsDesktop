# AnimalsDesktop Asset Queue

This queue tracks future art families without changing the fixed 100-variant runtime manifest. Promote an item into `internal/catalog` only when the parent thread intentionally expands or swaps the selectable catalog.

## Current Released State

- Latest infrastructure/page release: `v0.1.2`
- Next animal-content release: hold until one animal family reaches DeguDesktop-level asset coverage
- Selectable variants: exactly 100
- Accepted full degu motion sets: 11 variants
- Seed-stage non-degu variants: 89 variants
- Public priority families: chinchilla, macaroni mouse / fat-tailed gerbil, hamster, gecko

## Batch Strategy

One thread should produce more images, but for fewer animals:

- Source candidate thread: 4-6 variants, 4 source candidates per variant, 2 preview backgrounds for the best candidates.
- Motion thread: 1 species family, 1-3 accepted variants, 62 source frames per accepted variant.
- Integration thread: parent only; verifies alpha bounds, 96x64 readability, deterministic import, catalog coverage, and runtime behavior.
- Frame progress thread: one species family only; uses `cmd/auditframes` to track valid/missing/invalid standalone PNG frames without promoting unfinished or opaque ImageGen outputs.

This gives each thread enough output to compare candidates without pushing 600+ unreviewed frames through one context.

## Completion Contract

Use `docs/art-source/motion-contract.md` as the promotion and release gate. The runtime target stays DeguDesktop-compatible: 10 motion sets, 62 frames per set, 96x64 transparent frames, and the DeguDesktop frame slots preserved. A family is not complete just because it is selectable; it is complete for a version bump only when its source-truth art, 62-frame motion set, importer output, settings/tray reflection, visual review, and QA pass.

Work one family at a time. The next accepted-art version should finish a small first slice such as `chinchilla_standard_gray`; only then bump the release version and continue to the next chinchilla or macaroni mouse variant.

## Priority Queue

| Priority | Family | Candidate variants | Notes |
| --- | --- | --- | --- |
| P0 | Chinchilla | standard gray, beige, ebony, white mosaic | Replace tint/prototype sheets with accepted source-truth motion. |
| P0 | Macaroni mouse / fat-tailed gerbil | tan, gray, cream | Keep the round tail/body silhouette distinct from degu and hamster. |
| P0 | Hamster | golden Syrian, cream, black banded, white, cinnamon | Wheel-capable small-rodent profile remains valid. |
| P0 | Gecko | gray brown, leopard, tangerine, blizzard, albino, crested gecko harlequin | Low-crawler motion, no wheel. |
| P1 | Momonga family | Japanese dwarf flying squirrel / momonga, sugar glider gray, sugar glider mosaic, sugar glider leucistic | Needs glide membrane silhouette and skitter motion. |
| P1 | Small birds | sakura buncho / Java sparrow, budgerigar green, budgerigar blue, white buncho, cockatiel gray, lovebird peach-faced, zebra finch | New bird-perch / hop / flutter motion profile needed before runtime promotion. |
| P2 | Popular dogs | French Bulldog, Labrador, Golden Retriever, German Shepherd, Dachshund, Poodle, Beagle, Bulldog, Shiba Inu, Pomeranian, Corgi | Use breed silhouette families, not just recolors. |
| P2 | Popular cats | Maine Coon, Ragdoll, Persian, British Shorthair, Siamese, Sphynx, Scottish Fold, Bengal, Calico, Tuxedo | Use cat-stalk motion but vary silhouette where needed. |
| P2 | Rabbits | Holland Lop, Netherland Dwarf, Mini Rex, Lionhead, Dutch | Ear shape and hop contact points are the key differentiators. |

## Bird Motion Profile Draft

Bird variants should not use degu wheel, rabbit hop, or low-crawler motion. Draft profile:

- idle perch: subtle breathing and head turn
- hop: short two-foot ground hops with small vertical lift
- flutter: brief wing-open lift without leaving the 96x64 canvas
- preen: beak-to-wing grooming
- turn: quick body flip with tail alignment

Keep feet and tail visible. Avoid scenery, cages, branches, text, and multiple birds in the source frame.

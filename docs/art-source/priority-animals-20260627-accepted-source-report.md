# Priority animals accepted source report

Date: 2026-06-27

This parent pass promotes five completed one-frame ImageGen runs into
frame-addressable accepted `set00` source assets. These variants are accepted
as source art only; they are not added to `runtimeVariantIDs` or public
downloads in this pass.

## Verdict

Accepted as slot-based `set00` motion sources:

- `chipmunk_striped` under `docs/art-source/chipmunk/motion-source/`
- `gecko_leopard` under `docs/art-source/leopard-gecko/motion-source/`
- `whites_tree_frog_blue` under `docs/art-source/whites-tree-frog-blue/motion-source/`
- `cockatiel_normal_gray` under `docs/art-source/cockatiel/motion-source/`
- `java_sparrow_normal` under `docs/art-source/java-sparrow/motion-source/`

`budgerigar` remains in-progress. The local runs have partial frame coverage,
so it is not accepted or cataloged by this report.

## Mechanical QA

All five accepted directories were revalidated after copying into this branch.

| Variant | Frames | Audit result | Warnings | Sheet |
| --- | ---: | --- | ---: | --- |
| `chipmunk_striped` | 62/62 | valid=62 missing=0 invalid=0 | 52 | `sheets/chipmunk-striped-source-set00.png` |
| `gecko_leopard` | 62/62 | valid=62 missing=0 invalid=0 | 11 | `sheets/leopard-gecko-source-set00.png` |
| `whites_tree_frog_blue` | 62/62 | valid=62 missing=0 invalid=0 | 72 | `sheets/whites-tree-frog-blue-source-set00.png` |
| `cockatiel_normal_gray` | 62/62 | valid=62 missing=0 invalid=0 | 47 | `sheets/cockatiel-normal-gray-source-set00.png` |
| `java_sparrow_normal` | 62/62 | valid=62 missing=0 invalid=0 | 90 | `sheets/java-sparrow-normal-source-set00.png` |

Commands run:

```sh
go run ./cmd/auditframes -frames-dir <family>/motion-source/accepted-frames/set00 -strict -artifact-warnings -report <family>/motion-source/accepted-frames/set00-auditframes-report.json
go run ./cmd/assemblemotion -frames-dir <family>/motion-source/accepted-frames/set00 -out <family>/motion-source/sheets/<sheet>.png -report <family>/motion-source/accepted-frames/set00-assemblemotion-report.json
```

## Visual Review

- `chipmunk_striped`: pass. Striped coat, compact body, and tail stay readable.
  Turn frames are the main angle change; no persistent crop, prop, or duplicate
  animal was found.
- `gecko_leopard`: pass. Low crawler anatomy keeps four underside legs, spotted
  yellow-tan body, and a thick tail. No feet appear on the back/top of the body.
- `whites_tree_frog_blue`: pass with chroma-adjacent warning note. The high
  warning count is accepted after contact-sheet review because the blue body,
  toe pads, and legs remain intact with no visible green matte damage.
- `cockatiel_normal_gray`: pass. Crest, orange cheek patch, gray body, and long
  tail remain identifiable. Low fast frames stay bird-like rather than becoming
  a rodent silhouette.
- `java_sparrow_normal`: pass with angle note. Black head, white cheek, red
  beak, and gray body remain readable across the set. Higher artifact warnings
  are tied to small feet and fine alpha details, not visible body breakage.

## Release Status

These sources are ready for the next runtime/version lane. They should remain
off the public current-animal list until the corresponding app version,
Pages copy, release notes, and downloadable ZIP artifacts are aligned.

# Leucistic Sugar Glider 00-04 Parent Review

## Scope

- Variant: `leucistic_sugar_glider`
- Source variant: `sugar_glider_gray`
- Frames: `00-04`
- Status: parent-review pilot only; not accepted source and not runtime scoped.

## Method

The pilot uses a deterministic color-only transform from the accepted
`sugar_glider_gray` `set00` frames. Alpha, bbox, scale, baseline, pose,
gliding membrane shape, tail shape, and contact points are preserved exactly
from the accepted source frames.

Rejected attempts:

- `rejected/pink-overcast-v1/`: body and tail shifted too pink.
- `rejected/dark-source-markings-v2/`: mechanics were clean, but the dark
  source markings remained too strong for a leucistic read.

Current v3 candidate:

- white to warm-cream body and tail
- pale pink ears, feet, and nose details
- black eyes preserved
- sugar-glider morphology, membrane, tail, and baseline preserved

## QA: 00-04 Gate

`go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/frames -artifact-warnings -motion-warnings -report docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/qa/auditframes-00-04.json`

Result:

- `valid=5`
- `missing=57`
- `invalid=0`
- `warnings=0`

Contact sheet:

`docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/qa/leucistic-sugar-glider-00-04-contact.png`

## QA: 62-Frame Review Candidate

After the parent mechanics gate, the same v3 transform was expanded to all
`00-61` source frames for review.

Commands:

- `go run ./cmd/auditframes -frames-dir docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/frames -strict -artifact-warnings -motion-warnings -report docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/qa/auditframes-set00-strict.json`
- `go run ./cmd/assemblemotion -frames-dir docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/frames -out docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/sheets/leucistic-sugar-glider-source-set00.png -report docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/qa/assemblemotion-report.json`

Result:

- `valid=62`
- `missing=0`
- `invalid=0`
- `warnings=79`
- The accepted `sugar_glider_gray` source frames also report `warnings=79`
  with the same audit settings, so the deterministic color transform did not
  introduce additional alpha, bbox, contact, or motion warning count.

Full contact sheet:

`docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/qa/leucistic-sugar-glider-set00-contact-checker.png`

Source sheet:

`docs/art-source/one-frame-method-fullrun-20260628/leucistic-sugar-glider-color-pilot-00-04/sheets/leucistic-sugar-glider-source-set00.png`

## Parent Verdict

Mechanics pass for the 00-04 pilot: no scale jump, no contact drift, no membrane
loss, no cropped tail/ears/feet, and no species drift into a rodent or flying
squirrel.

The full 62-frame candidate also passes the mechanical gate and preserves the
accepted source motion exactly.

Visual risk to resolve before promotion: the source sugar-glider facial line is
still visible. This may be acceptable as a readable sugar-glider cue, but if the
target is a stricter pure white leucistic coat, the next step should be an
ImageGen/edit pilot or a stronger color transform that preserves only the black
eyes.

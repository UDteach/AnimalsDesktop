# Prompt Variants

Date: 2026-06-23

The first Nano Banana Pro run failed mainly as `POSE_DUP`: the face/style were stable, but almost every cell repeated the same right-facing pose. These variants change prompt wording only; the anchor set and guide images should stay the same for the next comparison.

## Variant A: Explicit Limb Delta

Goal: force visible motion differences while keeping the same source family.

```text
Create one single square 4x4 sprite sheet, 2048x2048.
Use the attached chinchilla style anchor only for identity, face, coat, outline, camera size, and shading. Do not copy its exact pose into every cell.
Use the attached green grid seed only for placement. Do not draw grid lines, labels, borders, ticks, shadows, floor, scenery, props, text, or checkerboard.

Each cell must contain one complete right-facing standard gray chinchilla on pure flat #00ff00.
All 16 cells must be different motion frames. Repeating the same body pose is a failure.
Keep the same scale, contact baseline, head size, eye, muzzle, ears, tail thickness, and fur style across all cells.

Motion plan, left to right then top to bottom:
1 idle: neutral standing, ears relaxed, all paws near ground.
2 idle: tiny head lift, ears slightly forward, body unchanged.
3 idle: tiny weight shift backward, tail tip slightly higher.
4 idle: tiny weight shift forward, front paw barely raised.
5 walk: front right paw forward, rear right paw back, body centered.
6 walk: front paw under chest, rear paw lifting, body slightly lower.
7 walk: front left paw forward, rear left paw back, tail counterbalanced.
8 walk: paws pass under body, body slightly higher.
9 walk: repeat gait phase with opposite paws, not a duplicate of cell 5.
10 walk: mid-step with one front paw lifted and one rear paw pushing.
11 walk: low contact step, head slightly forward.
12 walk: recovery step, tail slightly curled upward.
13 scurry: low stretched body, head forward, rear legs extended back.
14 scurry: body compressed, front paws tucked, rear paws under body.
15 scurry: long reach, front paws forward, tail trailing.
16 scurry: push-off, rear legs extended, body rising.

Clean antialiased edges. No white fringe, no green halo, no green dots inside the animal, no missing pixels, no disconnected body parts. Leave green padding around ears, whiskers, feet, toes, tail, and body in every cell.
```

## Variant B: Row-As-Animation Blocks

Goal: make the model think in four mini animation rows instead of one repeated contact sheet.

```text
Generate a 4 row x 4 column animation sprite sheet. One animal per cell. Pure #00ff00 background.
The attached chinchilla image is the locked art style and character identity. Keep the same side-view camera, size, face, eye, muzzle, ears, tail, outline, shading, and gray coat. The attached grid image is only a placement guide.

Rows:
Row 1 is an idle loop with four near-hold frames. Only ears, head angle, and weight shift change.
Row 2 is a walking loop with four distinct grounded step frames. Paw positions must visibly change in every cell.
Row 3 is a second walking loop with the opposite gait phases. It must not duplicate Row 2.
Row 4 is a fast low scurry loop with four distinct stretch/compress frames. Body length and paw reach must visibly change.

For every row, frame 1, frame 2, frame 3, and frame 4 must be different. If a motion frame is uncertain, change only paws, shoulders, haunch, and tail; do not change the face or species.

No grid lines, no labels, no numbers, no floor, no shadow, no prop, no scenery, no duplicate animals, no crossing cell boundaries, no cropped ears, no cropped feet, no cropped whiskers, no cropped tail. Clean sprite edges with no pinholes or green specks inside the body.
```

## Variant C: Anti-Duplication Contract

Goal: make duplication an explicit rejection condition in the prompt.

```text
Create a 2048x2048 square 4x4 sprite sheet for review. The output must be one sheet, not separate images.
Use the attached style anchor for the chinchilla identity only. Use the attached green grid seed for cell placement only.

Important rule: do not paste or repeat the same chinchilla pose. The 16 cells are an animation sequence. Every neighboring cell must change at least two of these: front paw position, rear paw position, shoulder height, hip height, head angle, tail curve, body stretch/compression.

Keep the face stable and calm in all cells. Keep the same right-facing direction, same camera distance, same body size, same coat, same outline, same soft shaded sprite style, same baseline.

Cell order:
1-4 idle micro-movement.
5-12 grounded walk cycle with eight different paw phases.
13-16 fast scurry with low stretch, compress, reach, push-off.

Pure flat chroma green #00ff00 in all empty space. No text, labels, grid lines, borders, safe boxes, floor, shelf, contact shadow, scenery, props, checkerboard, white fringe, green halo, noisy matte, pinholes, detached specks, or cropped anatomy.
```

## Current Preference

Variant A was tried in Google Flow / Nano Banana Pro on 2026-06-23 with `style-anchors.png` and `grid-seed.png`.

Result:

- Pose variation improved compared with the first Nano Banana Pro run.
- The visible grid guide was still reproduced as cell lines, so the raw output is rejected as `GUIDE_INK`.
- Local green normalization was useful only for diagnostics, not acceptance.
- Cache output `141219` parsed `3/16` after green normalization; the other 13 rejected mainly for `MAT_PINHOLE` or small detached alpha.
- Cache output `141429` parsed `0/16` after green normalization because boundary/grid residue remained.
- Visual faces stayed mostly coherent, but the matte defects and guide ink make the workflow not accepted.

Next variable should not be another prompt-only pose retry with the same visible grid seed. Change the guide strategy first.

## Variant D: Invisible Grid, Style Anchor Only

Goal: remove `GUIDE_INK` by not uploading a visible grid seed. This risks weaker placement, but it tests whether the model can keep the 4x4 layout from text alone.

```text
Create one square 4 by 4 sprite sheet, 2048 by 2048. Use the attached chinchilla image only for identity, face, coat, outline, shading, and camera size.

Do not draw or include any grid lines, borders, labels, dividers, safe boxes, ticks, baselines, numbers, text, floor, shadow, scenery, or props. The 4 by 4 grid is invisible layout only.

Pure flat chroma green background in all empty space. Place one complete right-facing standard gray chinchilla in each invisible cell, centered with green padding around ears, whiskers, feet, toes, tail, and body. No animal may touch the edge of its invisible cell.

All 16 cells must be different motion frames. Cells 1 to 4 are idle micro movement. Cells 5 to 12 are a grounded walk cycle with visible paw changes. Cells 13 to 16 are fast low scurry with stretch, compress, reach, and push-off.

Keep the face stable, scale stable, baseline stable, and coat stable. Clean smooth sprite edges. No white fringe, no green halo, no missing pixels, no green dots inside the animal, no detached specks, and no cropped anatomy.
```

## Variant E: Four-Frame Motion Block

Goal: reduce batch pressure. Use this before another 16-cell attempt if guide ink or pinholes repeat.

```text
Create one horizontal 4-frame sprite strip on pure flat chroma green. No visible cell borders, no grid lines, no dividers, no labels, no text, no floor, no shadow.

Use the attached chinchilla image only as the locked art style and identity. Keep the same side-view camera, body size, face, eye, muzzle, ears, tail thickness, gray coat, outline, and soft sprite shading.

The four frames are one small motion block: stretch, contact, compress, and push-off. Every frame must visibly change front paw position, rear paw position, body length, shoulder height, or hip height. Do not repeat the same pose.

One complete right-facing chinchilla per frame, with green padding around all anatomy. Clean antialiased outline, no white fringe, no green halo, no green dots inside the animal, no transparent pinholes, no detached specks, and no cropped anatomy.
```

## Variant F: Independent Images

Goal: avoid all cell-boundary and grid-line failure. Use this if Flow can return multiple independent images from one prompt.

```text
Generate exactly four separate images, not a sheet, grid, collage, or strip.

Each image is one complete right-facing standard gray chinchilla on pure flat chroma green, using the attached chinchilla image only for identity, face, coat, outline, camera size, and shading.

The four images are consecutive motion frames: stretch, contact, compress, and push-off. Keep the same scale, baseline, face, body size, and coat in all four images. The only changes should be paws, shoulders, hips, tail curve, and body stretch/compression.

No text, labels, borders, grid lines, floor, shadow, scenery, props, duplicate animals, cropped anatomy, white fringe, green halo, green dots inside the animal, missing pixels, transparent pinholes, or detached specks.
```

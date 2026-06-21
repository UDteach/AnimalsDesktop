# Chinchilla Single-Frame ImageGen Prompts

Use one prompt per output image. Accepted frames must come from standalone single-pose PNGs, not a sheet cut into frames.

## Base Prompt

Create one single side-view standard gray chinchilla desktop pet sprite pose, right-facing, `{pose_description}`. One complete animal only, realistic cute chinchilla anatomy, rounded ears, fluffy tail visible, tiny feet visible, whiskers visible, consistent 2D sprite style for a 96x64 desktop pet. Transparent background with real alpha channel. No checkerboard, no white background, no gray background, no shadow, no ground, no scenery, no text, no border, no props, no costume. Center the animal with safe margins so ears, feet, whiskers, and tail are not cropped.

## Pose Slot Descriptions

- Frames 00-03 idle: gentle breathing, tiny ear adjustment, calm weight shift.
- Frames 04-11 walk: soft four-foot walk, stable body volume, feet contacting a shared baseline.
- Frames 12-19 scurry: lower faster travel pose, tail following body, no jumping.
- Frames 20-25 forage/nibble: head lowers slightly, paws near mouth, compact body.
- Frames 26-31 action/hop-safe: short chinchilla hop or alert paw lift, no degu-like upright wheel pose.
- Frames 32-39 turn: body turn preparation and readable side-to-three-quarter transition.
- Frames 40-43 eat: chewing or nibbling pose, paws close to mouth.
- Frames 44-47 ground check: light paw scratch or substrate check, no scenery.
- Frames 48-51 stand/rest: seated alert or compact rest pose.
- Frames 52-55 groom: face groom or whisker clean with paws.
- Frames 56-61 reaction: wheel-safe small-rodent burst or alert reaction only after visual review.

## Current Attempt Log

- 2026-06-21: parent-thread single-pose idle prompt produced a `1536x1024` opaque PNG with `AlphaMin=255`; rejected as reference-only and not copied into the repo.

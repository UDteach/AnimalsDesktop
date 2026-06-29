# True Albino Chipmunk Motion Source

`true_albino_chipmunk` has an accepted `set00` source produced as a mechanical
coat-variant conversion from the accepted striped chipmunk source silhouette.

This promotion is source-art only. Runtime catalog, generated runtime sprite
sheets, GitHub Pages current-animal cards, release artifacts, tags, and public
release versions were not changed by this source promotion.

Why this path was used: ImageGen true-albino candidates preserved red eyes and
removed stripe pigment, but their tails drifted toward an albino rat/mouse
read. The accepted source keeps the established chipmunk body and tail
silhouette while flattening stripe pigment and adding red/pink eyes.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/true-albino-chipmunk-source-set00.png`
- Full contact: `contacts/true-albino-chipmunk-set00-full-contact.png`
- Animated preview: `contacts/true-albino-chipmunk-set00-preview.gif`
- Local final review: `qa/true-albino-chipmunk-local-final-review-20260629.md`

Validation:

```bash
go run ./cmd/auditframes -frames-dir docs/art-source/true-albino-chipmunk/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/assemblemotion -frames-dir docs/art-source/true-albino-chipmunk/motion-source/accepted-frames/set00 -out docs/art-source/true-albino-chipmunk/motion-source/sheets/true-albino-chipmunk-source-set00.png
```

`cmd/validatemotion -variant true_albino_chipmunk -require-accepted` is deferred
until catalog integration because the variant is not registered yet.

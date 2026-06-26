# Albino Chipmunk motion source

`albino_chipmunk` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/albino-chipmunk-source-set00.png`
- Full contact: `contacts/albino-chipmunk-set00-full-contact.png`
- Animated preview: `contacts/albino-chipmunk-set00-preview.gif`
- Visual QA: `qa/albino-chipmunk-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/albino-chipmunk/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant albino_chipmunk -require-accepted
```

This is accepted at current runtime-preview parity. It is not a full content
release set yet because only `set00` exists; a full release gate still requires
accepted `set00` through `set09`.

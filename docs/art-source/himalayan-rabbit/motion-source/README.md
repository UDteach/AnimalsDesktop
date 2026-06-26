# Himalayan Rabbit motion source

`himalayan_rabbit` has an accepted `set00` source generated with the one-frame
ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/himalayan-rabbit-source-set00.png`
- Full contact: `contacts/himalayan-rabbit-set00-full-contact.png`
- Animated preview: `contacts/himalayan-rabbit-set00-preview.gif`
- Visual QA: `qa/himalayan-rabbit-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/himalayan-rabbit/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant himalayan_rabbit -require-accepted
```

This is accepted at current runtime-preview parity. It is not a full content
release set yet because only `set00` exists; a full release gate still requires
accepted `set00` through `set09`.

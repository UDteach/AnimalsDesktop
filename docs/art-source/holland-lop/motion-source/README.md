# Holland Lop motion source

`holland_lop_broken_orange` has an accepted `set00` source generated with the
one-frame ImageGen workflow.

Accepted source:

- Frames: `accepted-frames/set00/frame-00.png` through `frame-61.png`
- Sheet: `sheets/holland-lop-broken-orange-source-set00.png`
- Full contact: `contacts/holland-lop-set00-full-contact.png`
- Animated preview: `contacts/holland-lop-set00-preview.gif`
- Visual QA: `qa/holland-lop-visual-qa-report.md`

Validation:

```sh
go run ./cmd/auditframes -frames-dir docs/art-source/holland-lop/motion-source/accepted-frames/set00 -strict -artifact-warnings
go run ./cmd/validatemotion -variant holland_lop_broken_orange -require-accepted
```

This is accepted at current runtime-preview parity. It is not a full content
release set yet because only `set00` exists; a full release gate still requires
accepted `set00` through `set09`.

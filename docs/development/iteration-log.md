# Iteration Log

## 2026-06-21

- Restored `kdevelopk.pages.dev` from the recent icon-gallery Cloudflare deployment and added AnimalsDesktop with a chinchilla-based icon.
- Made `UDteach/works-gallery` the public GitHub source for the icon gallery.
- Replaced the AnimalsDesktop GitHub Pages draft with a Coming Soon page. Pet lists will appear only after chinchilla is complete.
- Removed GitHub Pages download ZIP generation for the Coming Soon phase so unfinished builds are not published through the site.
- Narrowed the GitHub Pages artifact to the public Coming Soon `index.html` only.
- Scoped the Windows runtime selectable variants to `chinchilla_standard_gray` only, so unfinished backlog variants and degu variants are not selectable in AnimalsDesktop builds.
- Generated the first ImageGen chinchilla pose-source batch, rejected the non-transparent raw sheet as final art, and extracted 16 transparent 96x64 pose frames for the motion-source pass.
- Built a non-release `set00` 62-frame draft sheet from the 16 extracted poses to make the next motion-review pass concrete.
- Recorded the chinchilla-first release loop and the no-degu rule in `.codex/tasks/chinchilla-v010-release.md`.
- Simplified the public AnimalsDesktop GitHub Pages page to plain Coming Soon only and strengthened workflow checks against unfinished animal names, images, download links, and work-in-progress lists.
- Connected `chinchilla_standard_gray` to the 62-frame draft motion source sheet in `cmd/importanimals`; runtime sprite sheets now come from ImageGen-derived motion frames and remain explicitly marked non-release draft until accepted `set00` through `set09` variation exists.
- Added draft `set01` through `set09` chinchilla motion-source sheets and taught `cmd/importanimals` to import a complete `set00` through `set09` source family when present.
- Added `cmd/validatemotion` and wired the release workflow to fail when runtime motion sources are still draft instead of accepted.
- Rejected the mechanically shifted chinchilla draft `set01` through `set09` sheets after runtime review; accepted motion work must come from real ImageGen pose/frame generation and visual QA.
- Rejected worker-generated chinchilla sheets and opaque/checker-background PNGs as review-only, and tightened importer/validator checks so full-opaque motion frames cannot pass as accepted transparent sources.
- Added `cmd/assemblemotion` to assemble exactly 62 standalone 96x64 transparent PNGs into one 5952x64 motion source sheet, blocking wrong dimensions, empty frames, and opaque/checker backgrounds.
- Tried a parent-thread one-pose ImageGen idle prompt; the result was `1536x1024` with `AlphaMin=255`, so it was rejected and kept out of the repo. Added single-frame prompt and accepted-frame staging docs for the next clean generation pass.
- Added `cmd/auditframes` so partial one-pose PNG progress can be measured without promoting bad ImageGen output: valid, missing, invalid, and edge-warning counts are reported per set.
- Added `cmd/prepareframe` for one-pose ImageGen candidates: true-alpha input is fitted to 96x64, uniform edge backgrounds can be removed, and checker/noisy backgrounds are rejected before any visual-review promotion.
- Tested a pure green single-pose ImageGen fallback. The first prepared output visibly retained green background, so `cmd/prepareframe` was tightened to fail when cleaned content still touches the source canvas edge; the candidate remains rejected.
- Added explicit `chroma-green` preparation, transparent-RGB cleanup, and green despill. The first visually reviewed chinchilla idle frame was promoted to `accepted-frames/set00/frame-00.png`; `cmd/auditframes` now reports `valid=1 missing=619`.

## 2026-06-25

- Ran a Codex built-in ImageGen-only review pass for `hamster_golden_syrian`, `macaroni_mouse_tan`, `sugar_glider_gray`, and `rabbit_chestnut_agouti` under `docs/art-source/external-ai-trials/codex-imagegen-20260625/`, keeping all outputs out of `accepted-frames` and runtime/catalog paths.
- Produced and locally split 16-cell trial sheets for all four variants; all split outputs had non-empty alpha and no edge-alpha cells after 96x64 normalization.
- Attempted 62-frame sheets only for the strongest 16-cell directions, `macaroni_mouse_tan` and `sugar_glider_gray`, then rejected both because visual contact sheets showed row-wrapped/cropped fragments despite passing mechanical empty/edge checks.
- Recorded prompts, saved paths, visual decisions, and parent-thread integration recommendations in `docs/art-source/external-ai-trials/codex-imagegen-20260625/report.md`.
- Tested a Codex built-in ImageGen one-frame retry-until-pass method for eight high-risk `sugar_glider_gray` poses under the same isolated external trial directory. All eight attempt-01 outputs produced complete 96x64 review candidates with no empty alpha and no edge-alpha; the method appears stronger than sheet extraction for high-risk pose repair, with style/scale matching as the remaining integration risk. See `docs/art-source/external-ai-trials/codex-imagegen-20260625/one-frame-method-report.md`.

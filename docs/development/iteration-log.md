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

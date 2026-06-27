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

## 2026-06-26

- Promoted the original `guinea_pig_tricolor` one-frame `set00` run to accepted source art under `docs/art-source/guinea-pig/motion-source/`, assembled `guinea-pig-tricolor-source-set00.png`, and changed the catalog entry from shape-only to accepted motion source. The template-lock alternate was rejected because late frames were visibly damaged.
- Verified the guinea pig source with `auditframes` (`valid=62 missing=0 invalid=0` with artifact warnings), `assemblemotion`, `validatemotion -variant guinea_pig_tricolor -require-accepted`, `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, and `git diff --check`. `release_ready` remains false because only one runtime set exists.
- Rechecked interrupted fancy-rat generation threads and reconciled local state to frames `00-19`; current `auditframes -strict -artifact-warnings` reports the expected incomplete state `valid=20 missing=42 invalid=0 warnings=10`.
- Confirmed the GitHub Pages / asset queue planned animals as guinea pig, fancy rat, albino chipmunk, Richardson's ground squirrel, and Yorkshire terrier. Started four scoped local Codex generation threads for the remaining planned work, each limited to its own `docs/art-source/one-frame-method-fullrun-20260626/` run directory and no catalog/runtime/page/Git edits.
- Created `.codex/tasks/20260626-release-asset-parity.md` to define "implemented asset parity" as the current preview-runtime level: accepted 62-frame `set00`, assembled source sheet, catalog accepted source status, runtime validation with one accepted set, and prepared page/queue updates without publishing. Added heartbeat monitor `animaldesktop-asset-parity-monitor` for child generation recovery and parent-side QA.
- Parent-side monitoring found no mechanical QA issues in current partial canonical frames, but detected duplicate `attempt=01` manifest rows in the fancy-rat child run. Sent all active child threads a correction to continue only from missing frames and use higher attempt numbers for retries.
- Performed a stricter visual audit for `guinea_pig_tricolor`, especially the tricolor coat pattern. Added 2x and 4x pattern-review sheets plus rough color-ratio metrics; no pattern-break outliers or visual rejects were found. Frames 53-55 have the strongest posture-driven marking angle change but still read as the same tricolor guinea pig.
- Added `guinea_pig_tricolor` to the local preview-runtime list, regenerated runtime sprites/import report, refreshed the current animal icon and preview/page assets, and removed guinea pig from the coming-soon silhouette image. This remains a local release-prep state; publishing still requires a next-version decision and matching release ZIPs.
- Added parent-monitor progress contact sheets for the four parallel generation lanes. Fancy rat, albino chipmunk, and Richardson's ground squirrel are mechanically clean so far; albino chipmunk still needs final tail/species-read review. The initial Yorkshire terrier run was rejected at partial QA as too shepherd/spitz-like, and the child thread was instructed to start a silkier toy-dog retry run.
- Started a separate Yorkshire terrier silky retry thread `019f02d9-1a2a-7250-9562-69ae3047ee5a` because the original Yorkshire child continued advancing before receiving the queued correction. Latest monitor snapshot: fancy rat 53/62, albino chipmunk 27/62, Richardson's ground squirrel 27/62, Yorkshire silky retry 0/62.
- Re-ran parent monitor contacts after compaction/resume. Latest snapshot recorded in `parent-monitor/progress-summary.json`: fancy rat 58/62, albino chipmunk 33/62, Richardson's ground squirrel 32/62, Yorkshire original 26/62 rejected-reference, Yorkshire silky 8/62 rejected-reference.
- Rejected the Yorkshire silky retry direction at partial visual QA because it still reads too much like an upright-eared large black-and-tan dog. Added `yorkshire-terrier-longcoat-retry-set00-00-61.md` and started thread `019f02e5-096a-7c21-bd45-1592f7258087` for a long-coated low rectangular toy Yorkie retry with an early 00-04 stop gate.
- Passed the Yorkshire longcoat retry through the parent early visual gate for frames 00-04. It reads as a low, long-coated toy Yorkie rather than shepherd/spitz, so the child thread was told to continue frames 05-61. Fancy rat is now waiting on only frame 61 before full review.
- Updated the public page layout so current animal icons use fixed image slots and the hero preview uses right-facing, edge-spaced animals. The sugar glider page asset was mirrored for page consistency only; runtime sprite sheets were not changed.
- Added Himalayan rabbit, Djungarian hamster, Campbell hamster, grayish chestnut Netherland Dwarf, and Holland Lop to the coming-soon queue. Started five scoped generation threads for them, each limited to its own run directory and required to generate real color ImageGen source frames rather than silhouette-only or reused existing images.
- Rebuilt `docs/assets/animalsdesktop-coming-soon-silhouettes.png` from a fresh Page-specific ImageGen source stored at `docs/art-source/one-frame-method-fullrun-20260626/page-coming-soon/coming-soon-eight-animals-imagegen-source.png`, then converted only that generated source into black silhouettes. Existing runtime/source/prototype animal images are not used for the coming-soon silhouette.
- Recorded the release policy change: completed animals should move through small preview version bumps after parent QA and matching artifacts, while full 10-set DeguDesktop-level completion remains a separate release gate.
- Deleted heartbeat automation `animaldesktop-asset-parity-monitor` during parent-thread takeover, then revalidated fancy rat and accepted albino chipmunk. Albino chipmunk passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=70`; visual QA accepts it as a faint-striped albino chipmunk rather than the earlier rejected white-squirrel seed.
- Promoted `albino_chipmunk` to accepted `set00` source art under `docs/art-source/albino-chipmunk/motion-source/`, assembled `albino-chipmunk-source-set00.png`, added it to the local runtime preview list, regenerated runtime sprites/page preview assets, and moved it from the coming-soon queue to current-page/runtime preview prep. External publication still requires explicit release/version approval and matching ZIP artifacts.
- Accepted Richardson's ground squirrel after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=57`; visual QA accepts it as a low-tailed ground squirrel without chipmunk stripes or tree-squirrel tail drift.
- Added a separate `ground_squirrel` catalog species, promoted `richardsons_ground_squirrel` to accepted `set00` source art under `docs/art-source/richardsons-ground-squirrel/motion-source/`, regenerated runtime sprites/page preview assets for 10 local preview animals, and moved Richardson's ground squirrel out of the coming-soon queue.
- Accepted the Yorkshire Terrier longcoat retry after it reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=92`; visual QA accepts it as a compact long-coated Yorkie rather than the earlier shepherd/spitz-like rejected runs.
- Promoted `yorkshire_terrier_longcoat` to accepted `set00` source art under `docs/art-source/yorkshire-terrier/motion-source/`, regenerated runtime sprites/page preview assets for 11 local preview animals, and moved Yorkshire Terrier out of the coming-soon queue.
- Accepted Djungarian hamster after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=85`; visual QA accepts it as a gray-white winter-white dwarf hamster with a visible dorsal stripe. Promoted `djungarian_hamster` to accepted `set00` source art under `docs/art-source/djungarian-hamster/motion-source/`, regenerated runtime sprites/page preview assets for 12 local preview animals, and moved Djungarian hamster out of the coming-soon queue.
- Accepted Holland Lop after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=62`; visual QA accepts it as a broken orange lop rabbit with consistently dropped ears. Promoted `holland_lop_broken_orange` to accepted `set00` source art under `docs/art-source/holland-lop/motion-source/`, regenerated runtime sprites/page preview assets for 13 local preview animals, moved Holland Lop out of the coming-soon queue, and kept the page-only sugar glider icon right-facing per user request without changing runtime sprites.
- Accepted Netherland Dwarf after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=81`; visual QA accepts it as a compact short-eared grayish chestnut dwarf rabbit. Promoted `netherland_dwarf_chestnut` to accepted `set00` source art under `docs/art-source/netherland-dwarf/motion-source/`, regenerated runtime sprites/page preview assets for 14 local preview animals, and moved Netherland Dwarf out of the coming-soon queue.
- Accepted Himalayan rabbit after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=28`; visual QA accepts it as a cream-bodied Himalayan rabbit with dark ears, nose, feet, and tail. Promoted `himalayan_rabbit` to accepted `set00` source art under `docs/art-source/himalayan-rabbit/motion-source/`, regenerated runtime sprites/page preview assets for 15 local preview animals, and moved Himalayan rabbit out of the coming-soon queue.
- Accepted Campbell hamster after the child run reached 62 frames. It passed `oneframe_run.py review` with no issues and `auditframes -strict -artifact-warnings` with `valid=62 missing=0 invalid=0 warnings=71`; visual QA accepts it as a warm gray-brown Campbell dwarf hamster with a visible dorsal stripe. Frames `60` and `61` had a localized warm color drift, so the parent preserved originals, applied a color-only correction toward frames `56-59`, regenerated contact sheets, and reran audit before promotion. Promoted `campbell_hamster` to accepted `set00` source art under `docs/art-source/campbell-hamster/motion-source/`, regenerated runtime sprites/page preview assets for 16 local preview animals, prepared v0.1.5 page/docs/workflow checks, and moved the future queue to text-only normal striped chipmunk (`シマリス`) without reusing albino chipmunk art.
- Added `scripts/build_page_assets.py` to rebuild page icons and the hero preview from runtime sheets. It keeps `sugar_glider_gray` right-facing for page assets only, leaving runtime sprites unchanged.

## 2026-06-27

- Promoted five completed priority-animal one-frame runs to parent-owned
  accepted `set00` source assets without adding them to `runtimeVariantIDs`:
  `chipmunk_striped`, `gecko_leopard`, `whites_tree_frog_blue`,
  `cockatiel_normal_gray`, and `java_sparrow_normal`. Each has 62 accepted
  frames, an assembled source sheet, copied one-frame review JSON, contact
  sheet, and branch-local `auditframes` / `assemblemotion` reports.
- Added catalog metadata for those five accepted sources and introduced the
  `bird-hop` motion profile for cockatiel / Java sparrow source variants. The
  public v0.2.2 runtime list and Pages current-animal list remain at 16 animals;
  the next version bump/release lane should decide when to move accepted
  sources into runtime. Budgerigar remains in-progress and was not promoted.
- Reordered the GitHub Pages Coming Soon grid for the next animal wave while
  keeping the public downloads on `v0.2.2`. The display priority is now striped
  chipmunk, blue White's tree frog, leopard gecko, then the bird wave, followed
  by popular cat and dog candidates already covered by the page-specific
  silhouette source sheet.
- Split `scripts/build_page_assets.py` upcoming silhouette handling into the
  ImageGen source-sheet layout and the public display order, so future Pages
  ordering changes do not accidentally crop the wrong source cell. Strengthened
  `scripts/verify_page_release.py` to verify the Coming Soon priority order and
  all local `assets/...` references.
- Updated the public GitHub Pages Coming Soon section with page-only silhouette
  cards. The first cards are quick release candidates for
  `gecko_leopard`, `whites_tree_frog_blue`, and `chipmunk_striped`; the next
  generation wave prioritizes budgerigar, cockatiel, and Java sparrow, with
  lovebird and parrotlet as the 4-5 thread cap fillers.
- Kept the published v0.1.5 page scoped to sixteen current animals so the live
  download links do not claim unreleased animals. Leopard gecko, blue White's
  tree frog, and striped chipmunk remain local accepted `set00` candidates until
  a follow-up release/version step updates artifacts.
- Extended `scripts/build_page_assets.py` and both Pages workflows so upcoming
  silhouette PNGs are generated, verified, copied into the Pages artifact, and
  kept separate from runtime/accepted source frames.
- Corrected the GitHub Pages Coming Soon silhouette flow so it no longer uses
  deterministic hand-drawn silhouettes. `scripts/build_page_assets.py` now
  requires the page-specific ImageGen source sheet at
  `docs/art-source/one-frame-method-fullrun-20260627/page-coming-soon/coming-soon-fifteen-animals-imagegen-source.png`
  and extracts black silhouettes from that generated source only.
- Simplified the Coming Soon cards to name plus silhouette only. Added cat
  breed candidates from the 2026 iPet/Nyanpedia cat breed ranking top five:
  mixed cat, Scottish Fold, Munchkin, Ragdoll, and Minuet. The parent release
  policy is to bump preview versions in small increments as each animal
  graduates into the current runtime/page list.
- Updated the deployable public Pages links to the published `v0.2.1` Windows
  and Mac release assets. Leopard gecko, blue White's tree frog, and striped
  chipmunk are shown as Coming Soon silhouettes until their release artifacts
  and version links are aligned.
- Added JP/EN language switching for the static GitHub Pages site without a
  frontend framework. The page uses JP/EN buttons, saves the selection in
  `localStorage`, updates page metadata, navigation, download copy, animal
  names, upcoming names, feature copy, version notes,
  and keeps the release-verifier-sensitive Windows badge source markup intact.
- Added macOS JP/EN language persistence and menu/settings switching. The
  existing JSON settings now stores `language`, the Cocoa status menu can switch
  between Japanese and English, and the settings window is rebuilt after a
  language change so labels, tabs, option titles, animal names, and placeholders
  refresh consistently. Windows code and runtime assets were not changed.
- Expanded the macOS animal picker from the old five hardcoded labels to the
  runtime catalog, added fixed/selected/random animal menus in the status item,
  added per-pet size controls, and corrected the macOS bundle identifier to
  `com.udteach.animalsdesktop`.
- Added first-pass macOS multi-monitor support. The Mac overlay now persists a
  display ID, starts on the saved screen when it is available, falls back to the
  main screen when it is not, and exposes the display selector in both the status
  menu and the settings Motion tab. Local installed app
  `/Users/kyota/Applications/AnimalsDesktop.app` was checked and is still the
  published `v0.2.1` arm64 build until explicitly replaced.
- Strengthened macOS release QA for `v0.2.2`: Darwin tests now require the
  exact 16 release-scoped runtime animals, exercise every animal through fixed
  and per-pet selection, and verify every visible size step from 70% through
  120%. `scripts/verify_page_release.py` now also checks that the Pages current
  animal grid matches `catalog.RuntimeVariants()` exactly.
- Prepared the `v0.2.2` release docs and Pages copy for the Mac parity release.
  This release keeps leopard gecko, blue White's tree frog, and normal striped
  chipmunk in Coming Soon until their runtime assets are promoted in a later
  animal-addition lane. Verified `go run ./cmd/importanimals`,
  `python3 scripts/verify_page_release.py`, `go run ./cmd/validatemotion
  -runtime-only -require-accepted`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, `git diff --check`, and macOS arm64/amd64
  `VERSION=v0.2.2` ZIP builds. The motion validator remains `release_ready=false`
  for all runtime animals because this is a one-set preview, not the full
  10-set release gate.
- Reprioritized the GitHub Pages Coming Soon queue from the parent takeover
  thread. The first wave now favors シマリス, リューシスティックモモンガ,
  アフリカヤマネ, ネザーランドドワーフ（ヒマラヤン）, アメリカモモンガ,
  hamster / Djungarian color variants, fancy rat coat variants, and グレーうさぎ
  before the earlier frog / gecko / bird candidates. Reptile morph expansion is
  tracked as lower priority.
- Rebuilt Coming Soon silhouettes from a new page-specific 18-animal ImageGen
  source sheet at
  `docs/art-source/one-frame-method-fullrun-20260627/page-coming-soon/coming-soon-eighteen-animals-imagegen-source.png`.
  `scripts/build_page_assets.py`, `scripts/verify_page_release.py`, and both
  Pages workflows now verify the 18-card order.
- Prepared a bounded post-budgerigar asset planning lane in
  `docs/development/post-budgerigar-asset-queue-20260627.md`. The prep keeps
  Lane A's budgerigar run untouched, preserves the 16-animal runtime/public
  release boundary, queues cockatiel and Java sparrow accepted-source promotion
  before a new lovebird source lane, and records parallelization boundaries for
  later bird, cat, and dog work.
- Redirected the current phase from a Pages-specific goal to reusable asset
  output. Added `scripts/export_upcoming_asset_pack.py` and exported the current
  18 upcoming animals to `assets/source/upcoming/20260627/` as transparent color
  cutouts, normalized black silhouettes, copied source sheet, contact sheets,
  and `manifest.json`. The exporter uses connected-component extraction instead
  of equal grid cells so animals remain complete even when the ImageGen sheet
  is not a strict grid.
- Checked current ranking sources and recorded the implementation priority
  evidence in `docs/development/popularity-priority-sources-20260627.md`. Bird
  production keeps the user's explicit order of budgerigar, cockatiel, and Java
  sparrow before later lovebird/parrotlet work; cat and dog queues are recorded
  as breed-specific candidates for later lanes.

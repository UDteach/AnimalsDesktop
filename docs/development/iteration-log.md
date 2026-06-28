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
- Promoted `budgerigar_green_yellow` as a 62-frame accepted `set00` source
  asset. Parent-side visual QA replaced gray-face drift frames `40` and `41`
  with documented max-alpha blends from accepted neighbors and kept the earlier
  `60` neighbor fallback; originals were preserved under the fullrun
  `rejected/` directory. Final asset-only QA reports
  `auditframes -strict -artifact-warnings` as
  `valid=62 missing=0 invalid=0 warnings=50`, assembled
  `docs/art-source/budgerigar/motion-source/sheets/budgerigar-green-yellow-source-set00.png`
  sized `5952x64`, added catalog metadata as an accepted motion source, and
  regenerated the seed import outputs. This does not add the variant to
  runtime, Pages, release notes, tags, or downloadable artifacts.
- Reconfirmed the current production target as reusable source assets, not
  GitHub Pages deployment. Parent-gated `lovebird_peach_faced` frames `00-12`:
  `auditframes` reports `valid=13 missing=49 invalid=0 warnings=17`, every
  generated frame is `96x64`, and alpha bbox checks show a stable `47px` height
  and baseline `56` across the generated range. Checker/dark contacts showed no
  obvious background contamination, missing feet, clipping, or sudden scale
  jump, so the lovebird child lane was advanced to generate only `13-20` while
  preserving `00-12`. The parrotlet lane remains in progress after its `00-04`
  parent gate; no `05-12` canonical frames are visible yet.
- Started a third bounded asset-only lane for the top cat-breed queue:
  `scottish_fold_silver_tabby` frames `00-04` only under
  `docs/art-source/one-frame-method-fullrun-20260627/scottish-fold-silver-tabby-set00-oneframe-62/`.
  The lane is constrained to run-local raw/normalized frames, QA notes, and
  light/dark/checker contacts; parent review will gate folded ears, silver tabby
  markings, four attached paws, baseline stability, and non-bouncy `cat-stalk`
  motion before any `05-12` continuation.
- Parent-gated `parrotlet_green` frames `00-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=11`, every generated frame is
  `96x64`, and alpha bbox checks show stable `52px` height and baseline `58`
  across the generated range. Checker/dark contacts showed compact green
  parrotlet identity, visible attached feet, no background contamination, no
  clipping, and no sudden scale jump. The parrotlet child lane was advanced to
  generate only `13-20` while preserving `00-12`. The Scottish Fold lane has
  begun raw generation but has not reached the `00-04` parent gate yet.
- Early-rejected the first Scottish Fold canonical `frame-00` for scale, not
  mechanics: it was `96x64` and transparent, but the visible bbox was only
  `48x27` with baseline `57`, much smaller than the existing cat 96x64 preview
  evidence around `83x58`. A recovery instruction was sent to rebuild the
  anchor larger before generating or claiming `00-04`, targeting a cat-sized
  bbox around `72-86px` wide, `45-56px` high, and baseline `58-60`.
- Parent-gated `lovebird_peach_faced` frames `00-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=24`. Frames `13-19` keep the accepted
  lovebird scale, colors, baseline, and feet/contact. `frame-20` is much lower
  (`54x29` bbox) but was accepted as the start of sniff/nibble/forage rather
  than a scale failure; the lovebird child lane was advanced to generate only
  `21-28`, with `21-25` required to stay visually continuous with the low
  forage posture.
- Parent-gated corrected `scottish_fold_silver_tabby` frames `00-04`:
  `auditframes` reports `valid=5 missing=57 invalid=0 warnings=0`; corrected
  idle bboxes are `83-84px` wide and `45-47px` high with baseline `59`.
  `frame-04` is lower as cat-stalk start and was accepted. The child lane was
  advanced to generate only `05-12` while preserving `00-04`.
- Started a fourth bounded asset-only lane for `ragdoll_seal_bicolor` frames
  `00-04` under
  `docs/art-source/one-frame-method-fullrun-20260627/ragdoll-seal-bicolor-set00-oneframe-62/`.
  This lane uses the same cat-size gate as Scottish Fold and remains run-dir
  only: no catalog, runtime, Pages, release, tag, or Git operations.
- Parent-gated `parrotlet_green` frames `00-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=15`. Frames `13-17` and `19` preserve
  compact green parrotlet scale, baseline, and feet/contact; `frame-18` is a
  little longer and `frame-20` is wider/lower, but `20` was accepted as the
  start of sniff/nibble/forage rather than a scale failure. The parrotlet child
  lane was advanced to generate only `21-28`, with `21-25` required to stay
  visually continuous with the low forage posture.
- Reconfirmed the latest user correction as asset-only output, not GitHub
  Pages. Parent-gated `ragdoll_seal_bicolor` frames `00-04`:
  `auditframes` reports `valid=5 missing=57 invalid=0 warnings=0`, bbox checks
  show `00-03` around `88x48-50` with baseline `57`, and `04` is lower at
  `88x36` with the same baseline. Temporary parent checker/dark contacts show
  clear Ragdoll seal bicolor identity, attached paws, and no crop; the lane was
  advanced to generate only `05-12`. Parent also gated `lovebird_peach_faced`
  `21-28` mechanically (`valid=29 missing=33 invalid=0 warnings=40`) but
  rejected only `frame-26.png` for an abrupt wing-flare/tall-silhouette jump
  from the low forage frames. The lovebird lane was instructed to preserve
  `00-25` and `27-28`, replace only `26`, rebuild contacts, rerun
  `auditframes`, and stop before `29`. `parrotlet_green` still shows no visible
  `21-28` output after parent polls, so a light continuation reminder was sent;
  `scottish_fold_silver_tabby` has raw `05-10` visible but no canonical
  `05-12` gate yet.
- Parent-gated `parrotlet_green` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=25`. Contact review showed `21-24`
  as low forage, `25-26` as a natural rise, and `27-28` as upright recovery
  while preserving compact green parrotlet identity, baseline, and feet/contact.
  The parrotlet lane was advanced to generate only `29-36`. Parent-gated
  `scottish_fold_silver_tabby` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=5`; bbox checks and contacts show a
  stable low cat-stalk sequence with baseline `59-60`, folded ears, silver
  tabby identity, attached paws, and no abrupt scale jump. The Scottish Fold
  lane was advanced to generate only `13-20`, with stripe shimmer and vertical
  bounce called out as risks to avoid.
- Parent-gated `lovebird_peach_faced` frames `21-28` after retrying only
  `frame-26`: attempt 01 was rejected for an abrupt wing-flare/tall-silhouette
  jump, attempt 02 was rejected for staying too low (`54x27`) and preserving a
  `26 -> 27` size jump, and attempt 03 passed as a usable intermediate posture
  (`54x36`, baseline `56`). Final `auditframes` for the current partial run
  reports `valid=29 missing=33 invalid=0 warnings=40`; contact review shows a
  smoother low-forage-to-upright transition, so the lovebird lane was advanced
  to generate only `29-36` while preserving `00-28`.
- Parent-gated `parrotlet_green` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=29`; bbox checks show stable height
  `54` and baseline `59`, and contacts preserve compact green parrotlet
  identity with visible feet/contact and no abrupt scale jump. The parrotlet
  lane was advanced to generate only `37-44`. Parent-gated
  `ragdoll_seal_bicolor` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=1`; bbox checks show stable baseline
  `57`, and contacts preserve Ragdoll seal bicolor identity, white paws, seal
  mask/tail, long tail, and no crop. The Ragdoll lane was advanced to generate
  only `13-20`.
- Parent-gated `scottish_fold_silver_tabby` frames `13-20`: `auditframes`
  reports `valid=21 missing=41 invalid=0 warnings=8`; bbox checks show stable
  width `86`, height `37-41`, and baseline `59-60`. Contact review, including
  a parent-created temporary `00-20` full-range contact because the child did
  not save one, showed stable Scottish Fold silver tabby identity, folded ears,
  tail, attached paws, and no abrupt scale jump. The Scottish Fold lane was
  advanced to generate only `21-28`, with an explicit reminder to save
  full-range contacts next time.
- Parent-gated `lovebird_peach_faced` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=49`; bbox checks show stable
  baseline `56`, and contacts preserve peach-faced lovebird color/body
  identity, visible feet/contact, and no abrupt wing flare or vertical scale
  jump. The lovebird lane was advanced to generate only `37-44`.
- Parent checked `parrotlet_green` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=34`, but contact review rejected only
  `frame-43.png` because it jumps upright between low `42` and low `44`,
  recreating the abrupt-size-change issue. The parrotlet lane was instructed to
  preserve `00-42` and `44`, replace only `43` with a low eat/ground-check
  intermediate, rebuild contacts, and stop before `45`. Parent-gated
  `ragdoll_seal_bicolor` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=5`; contacts preserve low
  cat-stalk/sniff posture, baseline `57`, seal mask/tail, white paws, long
  tail, and no crop. The Ragdoll lane was advanced to generate only `21-28`.
- Parent-gated the retried `parrotlet_green` `frame-43`: replacement bbox is
  `82x40` with baseline `59`, `37-44` contact no longer has the one-frame
  upright jump, and `auditframes` now reports
  `valid=45 missing=17 invalid=0 warnings=35`. The parrotlet lane was advanced
  to generate only `45-52`.
- Parent-gated `scottish_fold_silver_tabby` frames `21-28`: `auditframes`
  reports `valid=29 missing=33 invalid=0 warnings=11`; bbox checks show stable
  width `86` and baseline `59-60`. Contact review, including a parent-created
  temporary `00-28` full-range contact because the child again did not save one,
  showed stable folded ears, silver tabby identity, tail, attached paws, and no
  abrupt scale jump. The Scottish Fold lane was advanced to generate only
  `29-36`, with another explicit reminder to save full-range contacts.
- Parent checked `lovebird_peach_faced` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=58`, but contact review rejected only
  `frame-43.png` because it jumps upright between low `42` and low `44`,
  matching the same abrupt-size-change issue found in parrotlet. The lovebird
  lane was instructed to preserve `00-42` and `44`, replace only `43` with a low
  eat/ground-check intermediate, rebuild contacts, and stop before `45`.
- Reconfirmed the current output target as asset/source-frame production, not
  GitHub Pages. Parent-gated `ragdoll_seal_bicolor` frames `21-28`:
  `auditframes` reports `valid=29 missing=33 invalid=0 warnings=6`; bbox checks
  show stable baseline `56-57`, and contacts preserve Ragdoll seal bicolor
  identity, white paws, seal mask/tail, long tail, and no abrupt one-frame scale
  jump. The Ragdoll lane was advanced to generate only `29-36`.
- Parent-gated `parrotlet_green` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=41`; contact review shows `48-52` as
  a stable upright/rest band after the ground-check frames, with compact green
  parrotlet identity, visible feet/contact, and no isolated one-frame size
  spike. The parrotlet lane was advanced to generate final range `53-61`, then
  stop for parent final review.
- Parent-gated the retried `lovebird_peach_faced` `frame-43`: replacement bbox
  is `54x26` with baseline `56`, `37-44` contact no longer has the one-frame
  upright jump, and `auditframes` still reports
  `valid=45 missing=17 invalid=0 warnings=58`. The lovebird lane was advanced
  to generate only `45-52`.
- Parent-gated `scottish_fold_silver_tabby` frames `29-36`: `auditframes`
  reports `valid=37 missing=25 invalid=0 warnings=14`; bbox checks show stable
  baseline `59-60`. Contact review shows `34-36` as a continuous upright/rise
  band, not an isolated body-size spike, while preserving folded ears, silver
  tabby identity, tail, and attached paws. The lane was advanced to generate
  only `37-44`, with full `00-44` contact saving required.
- Parent final-gated `parrotlet_green` as a 62/62 draft source:
  `auditframes` reports `valid=62 missing=0 invalid=0 warnings=42`; parent
  temporary contacts for `53-61`, `45-61`, and `00-61` show stable compact
  green parrotlet identity, visible feet/contact, no pale/gray drift, and no
  one-frame size spike. The lane was told to save missing final contacts only,
  with no frame changes or promotion.
- `parrotlet_green` final contact artifact completion passed. The lane saved
  `53-61`, `45-61`, and `00-61` light/dark/checker contacts, and parent reran
  `auditframes` as `valid=62 missing=0 invalid=0 warnings=42`. This is ready
  for parent accepted-source promotion review as an asset-only source.
- Promoted `parrotlet_green` into the parent branch as an asset-only accepted
  source. The 62 frames, contact evidence, QA notes, assembled source sheet,
  catalog metadata, generated seed source, and 10 sprite sheets are present.
  Validation passed: `validatemotion -variant parrotlet_green
  -require-accepted`, `importanimals`, targeted Go tests, and
  `git diff --check`. Runtime variants, Pages, release, tags, and deploy were
  not changed.
- Parent-gated `lovebird_peach_faced` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=71`; contact review shows `47-52` as
  a stable upright/rest band after the low ground-check frames, preserving
  peach/orange face, green body, visible feet/contact, and no isolated
  one-frame size spike. The lovebird lane was advanced to generate final range
  `53-61`, then stop for parent final review.
- Parent-gated `ragdoll_seal_bicolor` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=6`; bbox checks show stable baseline
  `57`, and contacts preserve Ragdoll seal bicolor identity, white paws, seal
  mask/tail, long tail, attached paws, and no abrupt one-frame scale jump. The
  Ragdoll lane was advanced to generate only `37-44`.
- Parent checked `ragdoll_seal_bicolor` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=9`, but contact review rejected only
  `frame-43.png` because it jumps upright between low `42` and low `44`.
  The Ragdoll lane was instructed to preserve `00-42` and `44`, replace only
  `43` with a low ground-check intermediate, save the missing `37-44` and
  `00-44` contacts, and stop before `45`.
- Parent checked `scottish_fold_silver_tabby` frames `37-44`: `auditframes`
  reports `valid=45 missing=17 invalid=0 warnings=17`, but labeled contact
  review rejected only `frame-43.png` because it jumps upright between low `42`
  and low `44`. The Scottish Fold lane was instructed to preserve `00-42` and
  `44`, replace only `43` with a low head-down stalk/ground-check intermediate,
  save the missing full `00-44` contacts, and stop before `45`.
- Parent-gated the retried `ragdoll_seal_bicolor` `frame-43`: replacement bbox
  is `87x27` with top `31` and baseline `57`, saved `37-44` and `00-44`
  contacts no longer have the one-frame upright jump, and `auditframes` still
  reports `valid=45 missing=17 invalid=0 warnings=9`. The Ragdoll lane was
  advanced to generate only `45-52`.
- Reconfirmed the current deliverable as reusable local assets/source frames,
  not GitHub Pages. Parent work continues on
  `/Users/kyota/.codex/worktrees/42d3/AnimalDesktop`
  (`codex/upcoming-asset-pack`) to avoid mixing in unrelated dirty Pages/release
  changes from the older `/Volumes/.../AnimalDesktop` `main` worktree. Parent
  poll found `lovebird_peach_faced` still at canonical `00-52` and
  `ragdoll_seal_bicolor` still at canonical `00-44`, so reminder prompts were
  sent for their already-approved next ranges only. Parent-gated the retried
  `scottish_fold_silver_tabby` `frame-43`: current bbox is `86x39` with top
  `22` and baseline `60`, matching the low `42` / `44` neighbor band. The
  `37-44` checker contact no longer has the one-frame upright jump, and
  `auditframes` reports `valid=45 missing=17 invalid=0 warnings=17`
  (missing `45-61` expected). The Scottish Fold lane was advanced to generate
  only `45-52`.
- Promoted `lovebird_peach_faced` into the parent branch as an asset-only
  accepted source. Parent verified the final `53-61`, `45-61`, and full
  `00-61` contacts; `auditframes` reports
  `valid=62 missing=0 invalid=0 warnings=92`, and
  `assemblemotion` produced
  `docs/art-source/lovebird/motion-source/sheets/lovebird-peach-faced-source-set00.png`.
  The 62 accepted frames, contact evidence, QA notes, catalog metadata,
  generated seed source, and 10 sprite sheets are present. Validation passed:
  `validatemotion -variant lovebird_peach_faced -require-accepted`,
  `importanimals`, and targeted Go tests. Runtime variants, Pages, release,
  tags, downloads, and deploy were not changed.
- Ragdoll and Scottish Fold lanes remain active after lovebird promotion.
  Ragdoll has advanced to canonical `45-49` with low stable bboxes through
  `48` and a moderate `49` rise, but no `45-52` contacts yet. Scottish Fold
  has raw `frame-45-attempt-01.png` and no new canonical `45` yet. To keep the
  child pool near the requested 4-5 lane cap, two new first-gate asset lanes
  were queued: `maine_coon_brown_tabby` (`00-04` only, pending worktree
  `local:19a9757d-ffaa-4346-836c-1923ef1f7c3c`) and
  `french_bulldog_fawn` (`00-04` only, pending worktree
  `local:6a7f3365-2213-4d6f-9a11-f53659e2c3b4`). Both are source-frame lanes
  only, with no catalog/runtime/Pages/release/Git edits allowed in child
  threads.
- Parent-gated `french_bulldog_fawn` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=3`, bboxes are stable around
  `56-60x52` with baseline `57`, and the parent temporary checker contact
  shows a coherent fawn French Bulldog with dark muzzle, upright ears, sturdy
  body, attached paws, and no scale jump into `04`. The lane was instructed to
  save missing `00-04` contacts/audit evidence, preserve `00-04`, generate
  only `05-12`, save `05-12` and `00-12` contacts, rerun `auditframes`, and
  stop before `13`.
- Parent recovery prompt sent for `ragdoll_seal_bicolor` `50-52` after the
  child rejected `frame-50-attempt-01-height-spike` and
  `frame-51-attempt-01-upright-raw`. The prompt specified a gradual bbox
  ladder from accepted `49` into alert/rest: `50` about `35-38px` high, `51`
  about `40-44px`, and `52` about `45-49px`, all keeping baseline `57`.
  Scottish Fold now has raw `45-48` but no canonical `45` yet. Maine Coon has
  raw `00-04` and was asked to normalize or reject, with no weak canonical
  frames accepted by the parent.
- Parent-gated `maine_coon_brown_tabby` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=0`. Canonical bboxes keep baseline
  `57`; `04` is low at `88x33` with `4px` side margins and reads as the
  cat-stalk/walk start rather than a crop. Parent checker review shows a
  coherent Maine Coon brown tabby with longhair body, dark tabby stripes,
  fluffy tail/chest, prominent ears, attached paws, and no scale jump. The lane
  was instructed to save proper canonical `00-04` contacts/audit evidence,
  preserve `00-04`, generate only `05-12`, save `05-12` and `00-12` contacts,
  rerun `auditframes`, and stop before `13`.
- Parent-gated `scottish_fold_silver_tabby` frames `45-52`: `auditframes`
  reports `valid=53 missing=9 invalid=0 warnings=19`. Contact review shows
  `49-52` as a coherent upright/rest band after low `45-48`, not an isolated
  one-frame spike, while preserving folded ears, silver tabby markings, tail,
  attached paws, and baseline. The lane was instructed to save missing full
  `00-52` contacts, preserve `00-52`, generate only final `53-61`, save
  `53-61` / `45-61` / `00-61` contacts, rerun `auditframes`, and stop for
  final review.
- Parent-gated `french_bulldog_fawn` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=3`. Bboxes keep height `52` and
  baseline `57`; parent temporary `05-12` and `00-12` checker contacts show a
  stable fawn French Bulldog with dark muzzle, upright ears, sturdy body,
  attached paws, and no body-size spike. The lane was instructed to save
  missing `05-12` / `00-12` contacts, preserve `00-12`, generate only `13-20`,
  save `13-20` / `00-20` contacts, rerun `auditframes`, and stop before `21`.
- Ragdoll recovery is partially improved: `frame-50` is now accepted at
  `86x38 top=20 baseline=57`, but `frame-51` attempts 02 and 03 were rejected
  for landing too tall and recreating the abrupt-rise problem. Canonical frames
  currently stop at `50`; continue only with a low-mid `51` that meets the
  ladder target.
- Ragdoll `frame-51` repeated failure reached the split-lane threshold for
  this asset queue. Attempt 04 was also rejected as too low/short, while
  attempt 02 was too tall. Parent paused the original Ragdoll thread at
  canonical `00-50` and opened focused rescue thread
  `019f09e8-f775-7ea2-b9f5-70ba261d4735` for `frame-51` only, using existing
  worktree `/Users/kyota/.codex/worktrees/359e/AnimalDesktop`. The rescue
  target anchors to accepted `frame-37.png` (`88x40 top=18 baseline=57`) and
  must produce width `86-88`, height `40-42`, top `16-18`, baseline `57`,
  then stop for parent gate.
- Maine Coon `05-12` remains on hold: parent still saw old `frame-11.png` at
  `82x52 top=6 baseline=57`, a one-frame upright spike between `10` and `12`.
  The lane was re-instructed to retry only `frame-11` with target height
  `42-45`, top `13-16`, baseline `57`, then rebuild `05-12` and `00-12`
  contacts before parent gate.
- Promoted `ragdoll_seal_bicolor` as a reusable accepted `set00` source asset,
  not as a Pages/runtime/release change. Parent-side final QA accepted the
  `53-61`, `45-61`, and full `00-61` contacts; the 8-column full grid showed
  stable seal-bicolor mask, pale cream body, white paws, fluffy tail, baseline,
  and no isolated one-frame size spike. Mechanical QA passed with
  `auditframes` `valid=62 missing=0 invalid=0 warnings=20`, and
  `assemblemotion` produced
  `docs/art-source/ragdoll/motion-source/sheets/ragdoll-seal-bicolor-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `ragdoll_seal_bicolor` variant as `motion_source_accepted`; the runtime
  variant list, Pages, release notes, tags, downloads, and deploy outputs were
  not changed.
- Revalidated the parent branch after Ragdoll promotion:
  `go run ./cmd/validatemotion -variant ragdoll_seal_bicolor -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, `go run ./cmd/validatemotion -runtime-only
  -require-accepted`, and `git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets still have one set, not the
  full 10-set release gate.
- Cleaned AppleDouble `._*` files from the Git common object directory after
  they caused `git diff --check` to fail with `non-monotonic index`. A follow-up
  `find "$(git rev-parse --git-common-dir)" -name '._*'` returned empty.
- Remaining active asset lanes after the Ragdoll promotion: Scottish Fold is
  accepted through `00-54`, Maine Coon through `00-20`, and French Bulldog
  through `00-21`. Continue monitoring by local file state only and avoid
  large thread transcript reads.
- Scottish Fold reached mechanical `62/62` (`valid=62 missing=0 invalid=0
  warnings=19`) but was not promoted. Parent visual QA rejected `55-61` because
  `55` became too pale/low-contrast and `56-61` drifted into a lighter
  bullseye/swirl pattern that did not match the darker accepted `00-54`
  silver-tabby coat. The child lane was instructed to preserve `00-54` and
  retry only `55-61` with the prior dark stripe style.
- Queued a new bounded first-gate lane for `domestic_shorthair_calico` under
  pending worktree `local:e7c9af71-ba0a-4403-96ab-eca3a7f5fc23`. This covers
  the popular mixed-cat / domestic shorthair family using an existing catalog
  ID, with first parent gate `00-04` only and no catalog/runtime/Pages/release
  writes in the child lane.
- Promoted `scottish_fold_silver_tabby` as a reusable accepted `set00` source
  asset, not as a Pages/runtime/release change. Parent-side final QA accepted
  the retried `55-61`, `45-61`, and full `00-61` contacts; the final band
  restored darker silver-tabby contrast, folded ears, attached paws, ringed
  tail, and stable baseline. Mechanical QA passed with `auditframes`
  `valid=62 missing=0 invalid=0 warnings=19`, and `assemblemotion` produced
  `docs/art-source/scottish-fold/motion-source/sheets/scottish-fold-silver-tabby-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `scottish_fold_silver_tabby` variant as `motion_source_accepted`; runtime
  variants, Pages, release notes, tags, downloads, and deploy outputs were not
  changed.
- Revalidated the parent branch after Scottish Fold promotion:
  `go run ./cmd/validatemotion -variant scottish_fold_silver_tabby
  -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`,
  `go run ./cmd/validatemotion -runtime-only -require-accepted`, and
  `COPYFILE_DISABLE=1 git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets currently have one set, not
  the full 10-set release gate.
- Parent-gated `maine_coon_brown_tabby` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=9`. Parent reviewed the `29-36` and
  `21-36` checker contacts; `frame-34` is narrow (`48x52`) but reads as a
  front-turn frame between `33` and `35`, not an isolated scale jump. The lane
  was instructed to preserve `00-36`, generate only `37-44`, and stop before
  `45`.
- Parent-gated `french_bulldog_fawn` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=4`. Parent reviewed the `29-36` and
  `21-36` checker contacts; `34-35` form the front-turn band and `36` returns
  toward right-facing side view with stable fawn coat, dark mask, upright ears,
  attached paws, and baseline. The lane was instructed to preserve `00-36`,
  generate only `37-44`, and stop before `45`.
- Parent-gated `domestic_shorthair_calico` frames `05-12`: `auditframes`
  reports `valid=13 missing=49 invalid=0 warnings=5`. Parent reviewed the
  `05-12` and `00-12` checker contacts; `07` and `12` are taller walking
  poses but not body-scale spikes, and the calico patches remain coherent. The
  lane was instructed to preserve `00-12`, generate only `13-20`, and stop
  before `21`.
- Parent-gated `maine_coon_brown_tabby` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=13`. Parent reviewed the `37-44`
  and `29-44` checker contacts; `40-41` are the low ground-check band and
  `42-43` return to standing without an isolated scale jump. The lane was
  instructed to preserve `00-44`, generate only `45-52`, and stop before `53`.
- Parent-gated `domestic_shorthair_calico` frames `13-20`: `auditframes`
  reports `valid=21 missing=41 invalid=0 warnings=11`. Parent reviewed the
  `13-20`, `05-20`, and `00-20` checker contacts; `13-17` are a low fast
  movement band, `18-19` return upward, and `20` starts sniffing with the head
  lowered. The lane was instructed to preserve `00-20`, generate only `21-28`,
  and stop before `29`.
- `french_bulldog_fawn` has generated through `42`, but the `37-44` parent gate
  is still incomplete. Latest local audit is `valid=43 missing=19 invalid=0
  warnings=4`; a focused reminder was sent for only `43-44`, preserving
  `00-42` and stopping before `45`.
- Parent-gated `french_bulldog_fawn` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=5`. Parent reviewed the `37-44` and
  `29-44` checker contacts; `40-44` read as a low ground-check band with stable
  fawn coat, dark mask, upright ears, attached paws, and baseline. The lane was
  instructed to preserve `00-44`, generate only `45-52`, and stop before `53`.
- Current child-lane polling state after the French gate: `maine_coon_brown_tabby`
  remains at `valid=45` waiting for `45-52`, `french_bulldog_fawn` is also
  `valid=45` waiting for `45-52`, and `domestic_shorthair_calico` has generated
  through `25` (`valid=26`) while waiting for `26-28` before the next parent
  gate.
- Parent-gated `maine_coon_brown_tabby` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=17`. Parent reviewed the `45-52`,
  `37-52`, and `00-52` checker contacts; `45-46` are a low ground-check band,
  `47-51` return through alert/rest scale, and `52` is acceptable as the
  face-groom start. The lane was instructed to preserve `00-52`, generate only
  `53-61`, then create `52-61`, `45-61`, and `00-61` contacts plus final
  `auditframes`.
- `french_bulldog_fawn` remains at `valid=45 missing=17 invalid=0 warnings=5`.
  A focused reminder was sent to preserve `00-44`, generate only `45-52`, and
  stop before `53`.
- `domestic_shorthair_calico` is at `valid=28 missing=34 invalid=0 warnings=20`.
  Parent partial contacts for `21-27` and `13-27` showed no immediate scale or
  anatomy blocker, so the lane was instructed to preserve `00-27`, generate
  only missing `28`, and stop before `29` for the `21-28` parent gate.
- Parent-gated `french_bulldog_fawn` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=6`. Parent reviewed the `45-52`,
  `37-52`, and `00-52` checker contacts; `45-47` are the low-to-rising
  transition, `48-51` stay stable in the alert/rest band, and narrower `52` is
  acceptable as the face-groom start. The lane was instructed to preserve
  `00-52`, generate only `53-61`, then create `52-61`, `45-61`, and `00-61`
  contacts plus final `auditframes`.
- Parent-gated `domestic_shorthair_calico` frames `21-28`: `auditframes`
  reports `valid=29 missing=33 invalid=0 warnings=20`. Parent reviewed the
  `21-28`, `13-28`, and `00-28` checker contacts; `21-22` and `28` are low
  sniff/action poses, with coherent calico patches, visible paws/contact, intact
  tail, baseline continuity, and no isolated body-size spike. The lane was
  instructed to preserve `00-28`, generate only `29-36`, and stop before `37`.
- Promoted `french_bulldog_fawn` as a reusable accepted `set00` source asset,
  not as a runtime, Pages, release, tag, download, or deploy change. Parent
  final QA accepted the `52-61`, `45-61`, and full `00-61` contacts on checker,
  dark, and light backgrounds; the fawn coat, dark mask, upright ears, compact
  bulldog body, attached paws, and baseline remain readable. Mechanical QA
  passed with `auditframes` `valid=62 missing=0 invalid=0 warnings=6`, and
  `assemblemotion` produced
  `docs/art-source/french-bulldog/motion-source/sheets/french-bulldog-fawn-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `french_bulldog_fawn` variant as `motion_source_accepted`; runtime variants,
  Pages, release notes, tags, downloads, and deploy outputs were not changed.
- Revalidated the parent branch after French Bulldog promotion:
  `go run ./cmd/validatemotion -variant french_bulldog_fawn -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, and
  `go run ./cmd/validatemotion -runtime-only -require-accepted` all passed.
  `release_ready=false` remains expected because accepted source assets
  currently have one set, not the full 10-set release gate.
- Promoted `maine_coon_brown_tabby` as a reusable accepted `set00` source
  asset, not as a runtime, Pages, release, tag, download, or deploy change.
  Parent final QA accepted the `52-61`, `45-61`, and full `00-61` contacts on
  checker, dark, and light backgrounds; brown tabby stripes, longhair body,
  fluffy ringed tail, ears, attached paws, and baseline remain readable.
  Mechanical QA passed with `auditframes`
  `valid=62 missing=0 invalid=0 warnings=21`, and `assemblemotion` produced
  `docs/art-source/maine-coon/motion-source/sheets/maine-coon-brown-tabby-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `maine_coon_brown_tabby` variant as `motion_source_accepted`; runtime
  variants, Pages, release notes, tags, downloads, and deploy outputs were not
  changed.
- Revalidated the parent branch after Maine Coon promotion:
  `go run ./cmd/validatemotion -variant maine_coon_brown_tabby -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`, and
  `go run ./cmd/validatemotion -runtime-only -require-accepted` all passed.
  `release_ready=false` remains expected because accepted source assets
  currently have one set, not the full 10-set release gate.
- Parent-gated `domestic_shorthair_calico` frames `29-36`: `auditframes`
  reports `valid=37 missing=25 invalid=0 warnings=37`. Parent reviewed the
  `29-36`, `21-36`, and `00-36` checker contacts; `32-36` form a coherent
  turn/front-facing band rather than an isolated body-size spike, while calico
  patches, visible paws/contact, tail, and baseline remain stable. The lane was
  instructed to preserve `00-36`, generate only `37-44`, and stop before `45`.
- Queued three asset-only child lanes, all capped at the first `00-04` parent
  gate and explicitly scoped away from Pages, runtime lists, releases, tags,
  downloads, and deploy output: `munchkin_brown_tabby` pending worktree
  `local:989cd7ad-752a-4242-9a78-cfcd7e2b8799`,
  `toy_poodle_apricot` pending worktree
  `local:b9029752-44bf-4a78-b43d-165529ca92e8`, and
  `british_shorthair_blue` pending worktree
  `local:fd8a46c7-2a07-450d-be98-cab5af70a4a6`.
- The queued lanes materialized as active child threads:
  `munchkin_brown_tabby` thread `019f0a7f-2439-77f3-974e-65025430c586` in
  `/Users/kyota/.codex/worktrees/1fe1/AnimalDesktop`,
  `toy_poodle_apricot` thread `019f0a7f-2449-7a92-8a59-c1c53ef5b278` in
  `/Users/kyota/.codex/worktrees/6572/AnimalDesktop`, and
  `british_shorthair_blue` thread `019f0a7f-2439-77f3-974e-64fd29322719` in
  `/Users/kyota/.codex/worktrees/f06b/AnimalDesktop`.
- Parent-gated `british_shorthair_blue` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=4`; bboxes keep height `52` and
  baseline `57`. The checker contact shows a coherent blue-gray British
  Shorthair with round head/body and visible paws/contact. The wider `02` and
  `04` frames are stride-width changes, not body-size spikes. The lane was
  instructed to preserve `00-04`, generate only `05-12`, and stop before `13`.
- Parent-gated `munchkin_brown_tabby` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=2`; bboxes are stable at `88x41` with
  baseline `57`. The checker contact shows a coherent short-legged brown tabby
  Munchkin with attached paws and ringed tail. Adjacent pixel diffs confirm the
  frames are not exact duplicates despite the stable silhouette. The lane was
  instructed to preserve `00-04`, generate only `05-12`, and stop before `13`.
- Parent-gated `domestic_shorthair_calico` frames `37-44`: `auditframes`
  reports `valid=45 missing=17 invalid=0 warnings=54`; all `37-44` frames keep
  height `52` and baseline `57`, with `39` as the wider stride frame. Parent
  reviewed the `37-44`, `29-44`, and `00-44` checker contacts; calico patches,
  visible paws/contact, tail, and baseline remain stable with no sudden
  body-size spike. The lane was instructed to preserve `00-44`, generate only
  `45-52`, and stop before `53`.
- Parent-gated `british_shorthair_blue` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=8`; `05-11` keep height `50-52` and
  baseline `57`, while `12` is an acceptable low sniff / ground-check pose at
  height `36` and the same baseline. Parent reviewed the `05-12` and `00-12`
  checker contacts and instructed the lane to preserve `00-12`, generate only
  `13-20`, and stop before `21`.
- Parent-gated `toy_poodle_apricot` frames `00-04`: `auditframes` reports
  `valid=5 missing=57 invalid=0 warnings=0`; bboxes stay around `58-60x52`
  with baseline `57`. The checker contact shows a readable apricot Toy Poodle
  with curly coat, rounded muzzle/ears, attached paws, and no one-frame
  body-size spike. The lane was instructed to preserve `00-04`, generate only
  `05-12`, and stop before `13`.
- Parent-gated `munchkin_brown_tabby` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=9`; bboxes stay `88x37-41` with
  baseline `57`. Parent reviewed the `05-12` and `00-12` checker contacts;
  `06` and `10` read as low walking phases, with stable stripes, ringed tail,
  attached paws, and no body-size spike. The lane was instructed to preserve
  `00-12`, generate only `13-20`, and stop before `21`.
- Parent-gated `toy_poodle_apricot` frames `05-12`: `auditframes` reports
  `valid=13 missing=49 invalid=0 warnings=0`; bboxes keep height `52` and
  baseline `57`, with `12` as a wider stride frame. Parent reviewed the
  `05-12` and `00-12` checker contacts; apricot coat, curly texture, rounded
  muzzle/ears, attached paws, and scale remain stable. The lane was instructed
  to preserve `00-12`, generate only `13-20`, and stop before `21`.
- Parent-gated `domestic_shorthair_calico` frames `45-52`: `auditframes`
  reports `valid=53 missing=9 invalid=0 warnings=64`; all `45-52` frames keep
  height `52` and baseline `57`. Parent reviewed the `45-52`, `37-52`, and
  `00-52` checker contacts. `49-52` read as a coherent reaction / upright band
  after the low `45-46` and walk `47-48` poses, not as an isolated scale jump.
  The lane was instructed to preserve `00-52`, generate only final `53-61`,
  and stop for final review.
- Parent-gated `british_shorthair_blue` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=14`; bboxes keep width `88`,
  baseline `57`, and height `35-44` across the low walking / ground-check band.
  Parent reviewed the `13-20` and `00-20` checker contacts and instructed the
  lane to preserve `00-20`, generate only `21-28`, and stop before `29`.
- Parent-gated `munchkin_brown_tabby` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=17`; bboxes keep width `88`,
  baseline `57`, and height `29-41`. Parent reviewed the `13-20` and `00-20`
  checker contacts; `20` reads as a low crouch / ground-check frame, not as a
  size jump. The lane was instructed to preserve `00-20`, generate only
  `21-28`, and stop before `29`.
- Parent-gated `toy_poodle_apricot` frames `13-20`: `auditframes` reports
  `valid=21 missing=41 invalid=0 warnings=0`; bboxes keep height `52` and
  baseline `57`, with `20` as a wider low ground-check start frame. Parent
  created and reviewed the `13-20` and `00-20` checker contacts and instructed
  the lane to preserve `00-20`, generate only `21-28`, continue naturally from
  the low `20` pose, and stop before `29`.
- Promoted `domestic_shorthair_calico` as a reusable accepted `set00` source
  asset, not as a runtime, Pages, release, tag, download, or deploy change.
  Parent final QA accepted the `53-61`, `45-61`, and full `00-61` contacts on
  checker, dark, and light backgrounds; calico patches, visible paws/contact,
  tail, ears, and baseline remain readable. Mechanical QA passed with
  `auditframes` `valid=62 missing=0 invalid=0 warnings=70`, and
  `assemblemotion` produced
  `docs/art-source/domestic-shorthair/motion-source/sheets/domestic-shorthair-calico-source-set00.png`
  at `5952x64`. Catalog metadata now marks the existing
  `domestic_shorthair_calico` variant as `motion_source_accepted`; runtime
  variants, Pages, release notes, tags, downloads, and deploy outputs were not
  changed.
- Revalidated the parent branch after Calico promotion:
  `go run ./cmd/validatemotion -variant domestic_shorthair_calico -require-accepted`,
  `go run ./cmd/importanimals`, `go test -buildvcs=false ./...`,
  `go vet -buildvcs=false ./...`,
  `go run ./cmd/validatemotion -runtime-only -require-accepted`, and
  `COPYFILE_DISABLE=1 git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets currently have one set, not
  the full 10-set release gate.
- Parent-gated `british_shorthair_blue` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=18`; `21-25` form the low
  ground-check band and `26-28` recover toward upright poses without a
  one-frame size spike. Parent reviewed the `21-28` and `00-28` checker
  contacts and instructed the lane to preserve `00-28`, generate only `29-36`,
  and stop before `37`.
- Parent-gated `toy_poodle_apricot` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=1`; `21-23` form a low ground-check
  band and `24-28` recover toward upright poses. Parent created and reviewed
  the `21-28` and `00-28` checker contacts and instructed the lane to preserve
  `00-28`, generate only `29-36`, and stop before `37`.
- Parent-gated `munchkin_brown_tabby` frames `21-28`: `auditframes` reports
  `valid=29 missing=33 invalid=0 warnings=22`; bboxes keep width `88`,
  baseline `57`, and height `27-37` across the low ground-check band. Parent
  reviewed the `21-28` and `00-28` checker contacts and instructed the lane to
  preserve `00-28`, generate only `29-36`, and stop before `37`.
- Parent-gated `british_shorthair_blue` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=22`; `32-35` form a coherent
  front-facing band, and `36` returns toward the right-facing pose. Parent
  reviewed the `29-36` and `00-36` checker contacts and instructed the lane to
  preserve `00-36`, generate only `37-44`, and stop before `45`.
- Parent-gated `munchkin_brown_tabby` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=25`; `33-34` form a front-turn band
  and `35-36` return toward the right-facing pose. Parent reviewed the `29-36`
  and `00-36` checker contacts and instructed the lane to preserve `00-36`,
  generate only `37-44`, and stop before `45`.
- Parent-gated `toy_poodle_apricot` frames `29-36`: `auditframes` reports
  `valid=37 missing=25 invalid=0 warnings=1`; `32-36` form a coherent
  front-facing band after the upright walk frames. Parent reviewed the `29-36`
  and `00-36` checker contacts and instructed the lane to preserve `00-36`,
  generate only `37-44`, and stop before `45`.
- Parent-gated `british_shorthair_blue` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=27`; `40-44` form a coherent low
  ground-check band. Parent reviewed the `37-44` and `00-44` checker contacts
  and instructed the lane to preserve `00-44`, generate only `45-52`, and stop
  before `53`.
- Parent-gated `munchkin_brown_tabby` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=26`; bboxes keep width `88`,
  baseline `57`, and height `26-41`. Parent reviewed the `37-44` and `00-44`
  checker contacts; `44` is a low ground-check pose, not a body-size spike. The
  lane was instructed to preserve `00-44`, generate only `45-52`, and stop
  before `53`.
- Parent-gated `toy_poodle_apricot` frames `37-44`: `auditframes` reports
  `valid=45 missing=17 invalid=0 warnings=1`; all frames keep height `52` and
  baseline `57`, with width changes matching the turn and ground-check poses.
  Parent reviewed the `37-44` and `00-44` checker contacts and instructed the
  lane to preserve `00-44`, generate only `45-52`, and stop before `53`.
- Parent-gated `british_shorthair_blue` frames `45-52`: `auditframes` reports
  `valid=53 missing=9 invalid=0 warnings=31`; `45-48` stay low and `49-52`
  recover into upright alert/rest poses with stable blue-gray coat, round
  head/body, tail attachment, and visible paws/contact. Parent reviewed the
  `45-52` and `00-52` checker contacts and instructed the lane to preserve
  `00-52`, generate only final `53-61`, and stop for final review.
- Revised the active asset/release target after user confirmation: the target
  set is now 19 animals, adding `quokka`, `roborovski_hamster`, and
  `guinea_pig_russian_smoke_white` to the previous 16-count plan. The user also
  explicitly requested the lane to proceed through asset QA, settings UI/menu
  confirmation for every target animal, and Release. Child generation lanes
  remain asset-only; the parent lane now owns runtime/UI/page/release
  integration after all 19 assets pass parent QA.
- Promoted `british_shorthair_blue` as the 14th accepted source asset in the
  19-target plan. Parent accepted the full `00-61` contact evidence from the
  child lane and copied the 62 source frames to
  `docs/art-source/british-shorthair/motion-source/accepted-frames/set00/`.
  Mechanical QA passed with `auditframes` `valid=62 missing=0 invalid=0
  warnings=31`; `assemblemotion` produced
  `docs/art-source/british-shorthair/motion-source/sheets/british-shorthair-blue-source-set00.png`.
  Catalog metadata now marks `british_shorthair_blue` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Revalidated the parent branch after British Shorthair promotion:
  `go run ./cmd/validatemotion -variant british_shorthair_blue
  -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`, and
  `COPYFILE_DISABLE=1 git diff --check` all passed. `release_ready=false`
  remains expected because accepted source assets currently have one set, not
  the full 10-set release gate.
- Sent final-frame child continuations for `munchkin_brown_tabby` and
  `toy_poodle_apricot`, each scoped to generate only `53-61`, create
  `53-61`/`45-61`/`00-61` contacts, run strict `auditframes`, and stop for
  parent final QA without catalog/runtime/UI/Pages/release edits.
- Added mechanical motion-consistency review to `cmd/auditframes` via
  `-motion-warnings`. The new warnings flag adjacent or isolated changes in
  bbox width/height, contact baseline, body alpha area, and bbox fill ratio so
  sudden size jumps and body-ratio outliers are review candidates before parent
  acceptance. This does not make warnings fatal; the parent must still inspect
  flagged frames on contact sheets and record whether they are intentional
  crouch/turn/groom poses or true rejects. `go test -buildvcs=false
  ./cmd/auditframes` passed, and British Shorthair's new motion report flags
  frames `12`, `25`, `26`, `40`, and `49` for contact-sheet review.
- Promoted `toy_poodle_apricot` as the 15th accepted source asset in the
  19-target plan. Parent `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=1`;
  the lone `frame-27` lower-ledge warning was reviewed on the `21-28` and full
  `00-61` contacts and accepted as connected curly body/leg fur, not a floor,
  prop, shadow, or detached shelf. The `45-61` contact keeps stable size, body
  ratio, and baseline. The accepted frames now live under
  `docs/art-source/toy-poodle/motion-source/accepted-frames/set00/`, and
  `assemblemotion` produced
  `docs/art-source/toy-poodle/motion-source/sheets/toy-poodle-apricot-source-set00.png`.
  Catalog metadata now marks `toy_poodle_apricot` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Revalidated after Toy Poodle promotion: `go run ./cmd/validatemotion -variant
  toy_poodle_apricot -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./cmd/auditframes ./cmd/importanimals
  ./internal/catalog`, and `COPYFILE_DISABLE=1 git diff --check` passed.
- Promoted `munchkin_brown_tabby` as the 16th accepted source asset in the
  19-target plan. Parent `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=37`;
  most warnings are lower ledge/floor candidates from the low, long Munchkin
  body. Motion warnings at `20` and `44` were reviewed on slot contacts and
  accepted as lower sniff/ground-check pose changes, not sudden scale errors.
  The accepted frames now live under
  `docs/art-source/munchkin/motion-source/accepted-frames/set00/`, and
  `assemblemotion` produced
  `docs/art-source/munchkin/motion-source/sheets/munchkin-brown-tabby-source-set00.png`.
  Catalog metadata now marks `munchkin_brown_tabby` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Revalidated after Munchkin promotion: `go run ./cmd/validatemotion -variant
  munchkin_brown_tabby -require-accepted`, `go run ./cmd/importanimals`,
  `go test -buildvcs=false ./cmd/auditframes ./cmd/importanimals
  ./internal/catalog`, and `COPYFILE_DISABLE=1 git diff --check` passed.
- Parent-gated new 19-target lanes through `05-12` where usable:
  `roborovski_hamster` passed with parent `auditframes` `valid=13 missing=49
  invalid=0 warnings=3` and stable sandy Roborovski read;
  `guinea_pig_russian_smoke_white` passed with `valid=13 missing=49 invalid=0
  warnings=22`, where lower-body/floor heuristic warnings were accepted after
  checker/dark contact review confirmed a white guinea pig with gray ear/nose
  and no detached shadow. `quokka` current `00-12` was rejected at 4x review as
  too rodent/kangaroo-rat-like because of the long thin tail and low rat-like
  body; the worker was interrupted to preserve it as rejected-reference and
  restart `00-04` with stronger wallaby/quokka constraints.
- Parent-gated the remaining three 19-target lanes through later partial
  ranges. `roborovski_hamster` now passes through `45-52` with
  `auditframes -artifact-warnings -motion-warnings` reporting
  `valid=53 missing=9 invalid=0 warnings=63`; the `45-52` contact confirms the
  width changes are turn/upright pose changes, not isolated size jumps. The user
  noted the Roborovski tone should be corrected, so a parent candidate
  color-only correction for frames `29-52` was prepared toward the brighter
  `00-28` sandy coat while preserving white belly/cheek/eyebrow and pink
  ear/foot/nose pixels. It will be applied only after final `53-61` exists, with
  originals and a correction report preserved. `quokka` regenerated art now
  passes through `29-36` with `valid=37 missing=25 invalid=0 warnings=44`; the
  stout attached tail and upright wallaby/quokka read avoid the rejected
  long-tail rodent direction. `guinea_pig_russian_smoke_white` now passes
  through `45-52` with `valid=53 missing=9 invalid=0 warnings=101`; width stays
  `84`, fill ratio remains about `806-848` permille, and the warnings are broad
  white-body lower-edge heuristics rather than detached floor or body-size
  spikes. `numpy 2.5.0` was installed into the user-site Python 3.14 package
  directory after direct system pip install was blocked by the externally
  managed Python environment.
- Promoted `roborovski_hamster` as the 17th accepted source asset in the
  19-target plan. Parent final `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=81`.
  Frames `29-52` received a documented color-only correction toward the early
  `00-28` sandy coat after the user called out tone drift; originals and the
  correction report were preserved under
  `docs/art-source/roborovski-hamster/motion-source/qa/`. The accepted frames
  now live under
  `docs/art-source/roborovski-hamster/motion-source/accepted-frames/set00/`,
  and `assemblemotion` produced
  `docs/art-source/roborovski-hamster/motion-source/sheets/roborovski-hamster-source-set00.png`.
  Catalog metadata now marks `roborovski_hamster` as
  `motion_source_accepted`; runtime variants, Pages, release notes, tags,
  downloads, and deploy outputs were not changed in this step.
- Promoted `guinea_pig_russian_smoke_white` as the 18th accepted source asset
  in the 19-target plan. Parent final `auditframes -strict -artifact-warnings
  -motion-warnings` passed with `valid=62 missing=0 invalid=0 warnings=119`;
  the high warning count was accepted after contact review because the broad
  white low body triggers lower-edge heuristics while width, baseline, fill
  ratio, and Russian-smoke gray ear/nose read remain stable. The accepted frames
  now live under
  `docs/art-source/guinea-pig-russian-smoke-white/motion-source/accepted-frames/set00/`,
  and `assemblemotion` produced
  `docs/art-source/guinea-pig-russian-smoke-white/motion-source/sheets/guinea-pig-russian-smoke-white-source-set00.png`.
  Catalog metadata now marks `guinea_pig_russian_smoke_white` as a distinct
  `motion_source_accepted` variant rather than reusing existing tricolor guinea
  pig art. Validation passed: both new variants pass
  `validatemotion -require-accepted`, `go run ./cmd/importanimals` imported
  `105` seed variants, and `go test -buildvcs=false ./cmd/auditframes
  ./cmd/importanimals ./internal/catalog` passed.
- Completed the 19/19 accepted-source wave by promoting `quokka`. Parent final
  `auditframes -strict -artifact-warnings -motion-warnings` passed with
  `valid=62 missing=0 invalid=0 warnings=90`; frame `57` needed two rejected
  replacements before the final stout, attached-tail quokka read was accepted.
  Accepted frames, contact sheets, QA notes, assembled source sheet, catalog
  metadata, generated seed source, and 10 runtime sprite sheets now exist for
  `quokka`.
- Per the user pre-release visual flags, repaired the `gecko_leopard` middle
  turn scale jump by resizing/re-anchoring frames `33-35`, and repaired
  `domestic_shorthair_calico` contact drift by re-anchoring frames `07`, `27`,
  and `53`. Originals and JSON reports were preserved under each animal's
  `motion-source/qa/` repair directory. Post-repair `auditframes` passed for
  both accepted source sets, with contacts reviewed on checker/light/dark
  backgrounds.
- Integrated all 19 new accepted animals into the v0.2.4 runtime/page release
  scope, raising runtime selection to 35 variants. GitHub Pages now shows 35
  current animal cards, 12 remaining upcoming silhouettes, v0.2.4 download
  links, v0.2.4 version-history copy, and explicit future roadmap text for the
  remaining Pages candidates, the white lionhead-pattern rabbit, and the
  special low-motion shoebill. Release workflow prerelease handling and
  `docs/releases/v0.2.4.md` were updated.
- Release-prep validation passed locally: `go run ./cmd/importsheet`, `go run
  ./cmd/importanimals`, `python3 scripts/build_page_assets.py`,
  `python3 scripts/verify_page_release.py`, JS syntax check for
  `docs/index.html`, `go run ./cmd/validatemotion -runtime-only
  -require-accepted` (35/35 accepted source; 35 expected one-set preview
  warnings), `go test -buildvcs=false ./...`, `go vet -buildvcs=false ./...`,
  Windows amd64/386 cross-builds, Windows amd64 compile-only test, macOS
  arm64/amd64 ZIP builds, and `COPYFILE_DISABLE=1 git diff --check`. Playwright
  browser QA confirmed the JP and EN Pages views, 35 animal cards, 12 upcoming
  cards, v0.2.4 version history, and future roadmap text render.

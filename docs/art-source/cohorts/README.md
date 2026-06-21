# ImageGen Cohort Contract

AnimalsDesktop uses 10 ImageGen cohorts with 10 selectable variants each. Each cohort owns only its own directory:

`docs/art-source/cohorts/cohort-XX/`

Parent-thread integration promotes accepted cohort output into the importer manifest and regenerates `assets/sprites`. Cohort workers must not edit `internal/catalog`, `assets/sprites`, release workflows, or files outside their assigned cohort directory.

## Required Cohort Output

For each variant:

- `source-truth/<variant-id>.png`: one complete animal on transparent background.
- `motion-source/<variant-id>/`: 62 transparent source frames preserving camera, scale, anatomy, baseline, and contact points.
- `preview/<variant-id>-96x64-light.png`: 96x64 readability preview on light background.
- `preview/<variant-id>-96x64-dark.png`: 96x64 readability preview on dark background.

For each cohort:

- `cohort-report.md`: source prompts, acceptance notes, and known risks.
- `known-risks.md`: cropping, anatomy, motion, alpha, scale, or species-identity risks that parent integration must review.

## Acceptance Rules

- Each source frame contains exactly one animal, no text, no borders, no scenery, no shadows, no costume, and no human-like pose.
- Do not make chinchillas, macaroni mice, rabbits, cats, dogs, reptiles, or other species by recoloring degu sprites.
- Coat variants may share a species source family only after the species silhouette and motion set are stable.
- Low-crawler, snake, frog, and tortoise profiles must not use degu wheel or upright hopping actions.
- Parent thread decides when a cohort is promoted from `prototype_only` to an accepted ImageGen source status.

## Cohort Map

- `cohort-01`: degu core coats 1-10.
- `cohort-02`: degu cream pied, chinchillas, macaroni mice, first rabbits.
- `cohort-03`: remaining baseline rabbits, baseline dogs, first cats.
- `cohort-04`: remaining baseline cats, geckos, first hamster.
- `cohort-05`: remaining hamsters, ferrets, guinea pigs.
- `cohort-06`: hedgehogs, squirrels, foxes, red pandas, otter, sugar glider.
- `cohort-07`: capybaras, tortoises, first popular dog additions.
- `cohort-08`: remaining popular dog additions, first popular cat additions.
- `cohort-09`: remaining popular cat additions, first popular rabbit additions.
- `cohort-10`: remaining rabbit addition, small animals, reptile/amphibian additions.

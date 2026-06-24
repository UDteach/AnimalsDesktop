# Diagnostic Notes

## Summary

The DeguDesktop-style dot-family direction is visually promising but not technically passable yet.

Generated outputs:

- `d01..d10`: 10 raw one-frame candidates.
- `d01`, `d02`, `d03`, `d07`, `d08`, `d09`, and `d10` prepared successfully and produced coherent 96x64 dot-family sprites.
- `d04`, `d05`, and `d06` failed `prepareframe` because chroma-green removal created transparent pinholes.
- All prepared candidates were rejected by the technical gate because they retained 1-5 transparent pinhole pixels.

Visual read:

- The prepared sequence is much more internally consistent than the mixed existing chinchilla set00 draft.
- Face, ear, tail, body size, and baseline are stable enough to justify another dot-family attempt.
- The main blocker is not source-family drift; it is tiny internal pinholes after background removal.

Extra diagnostics:

- Bright magenta background did not work with the current `prepareframe -background auto` path; the background was not removed.
- A transparent-background request returned a white-background image. `prepareframe -background auto` removed it, but the prepared sprite still had a 1px transparent pinhole.

Next useful attempt:

- Continue the dot-family lane, but change prompting away from chroma key language.
- Prefer true transparent or white/uniform background plus `prepareframe -background auto`.
- Add stronger wording for closed solid interior pixels, no white/missing pinholes, and no background-colored pixels inside the animal.
- Do not locally fill pinholes for acceptance unless a separate repair policy is explicitly approved.

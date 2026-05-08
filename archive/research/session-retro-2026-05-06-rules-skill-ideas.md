# Session Retro: Rules And Skill Improvement Ideas

Date: 2026-05-06

Scope: retrospective on the PR #2 and PR #3 workflow, Block 13B implementation,
parallel subagent work, pi review, PR merge, and cleanup.

## Observed Process Risks

1. Subagent slice completion is not integration evidence.
   - Worker B correctly added safe command descriptors.
   - Worker A correctly added preview/doctor.
   - Integration still leaked raw command argv in preview until the parent audit
     compared the two slices.

2. pi reviews with tool access can hang or produce empty artifacts on broad PRs.
   - The useful reviews came from minimal-context, no-tool prompts over key
     files.
   - Empty or hung review artifacts must not be counted as review evidence.

3. PR merge commands can fail from a feature worktree.
   - `gh pr merge` attempted local branch cleanup while the branch was checked
     out in another worktree.
   - The safe path is to merge from the main worktree, then remove feature
     worktrees and branches explicitly.

4. PR review evidence should be committed only after it contains usable content.
   - Empty review files from hung runs should be deleted or replaced before
     staging.

5. Safety-sensitive CLI output needs negative tests, not only positive shape
   tests.
   - Block 13B needed assertions that secret-like argv values do not appear in
     preview/dry-run output.

6. Merge completion must be verified against `origin/main`, not inferred from
   command success text.
   - PR #3 merged remotely even though local cleanup failed.
   - The final state was only clear after checking PR state and fetching
     `origin/main`.

## Proposed Rule Updates

### AGENTS.md

Rejected after GLM review.

Reason: `AGENTS.md` already delegates block execution to
`sdp-trace-trust-workflow`, and adding workflow details would push the file
toward its 100-line decomposition limit. The correct home for these details is
the skill.

### `sdp-trace-trust-workflow` Skill

Add detailed workflow clauses:

- after subagents return, inspect diffs and run an integration audit before
  trusting their completion claims;
- for safety-sensitive outputs, add negative leak tests using secret-like
  markers;
- default pi review mode is `--no-tools --no-context-files` over explicitly
  named files; use tool access only for narrow questions;
- if a pi run hangs or produces an empty file, delete/replace the artifact and
  do not count it;
- commit review evidence only after the file has usable content;
- merge PRs from the main worktree, verify PR state and `origin/main`, then
  remove feature worktrees and branches;
- never count local merge/cleanup failure as product failure if the remote PR
  state proves merge succeeded, but do verify and clean up explicitly.

## Expected Behavior Next Time

When the user writes "берем блок в работу", the agent should:

1. land the current approved PR;
2. create a fresh block worktree;
3. split subagent work by disjoint write scopes;
4. run parent integration audit after subagents return;
5. run tests, schema checks, hygiene checks, and safety negative checks;
6. run minimal-context pi review;
7. fix and re-review until no critical/major findings remain;
8. create PR;
9. run a separate PR-review cycle;
10. merge from the main worktree;
11. fetch and verify `origin/main`;
12. remove the feature worktree and stale branches.

## pi Review Disposition

Reviewed with:

- MiniMax-M2.7: accepted the proposal with no critical or major findings.
- GLM-5.1: revised. Major finding: do not add workflow details to `AGENTS.md`;
  put them only in the skill. Major finding accepted. GLM also requested an
  explicit trigger for negative leak tests; accepted.

Applied disposition:

- `AGENTS.md` left unchanged.
- `sdp-trace-trust-workflow` updated with parent integration audit,
  safety-sensitive negative leak test trigger, minimal-context pi review
  default, empty review artifact hygiene, and main-worktree PR merge cleanup.

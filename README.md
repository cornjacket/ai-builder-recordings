# ai-builder-recordings

Recording storage for [ai-builder](https://github.com/cornjacket/ai-builder)
replay regression tests.

Each branch holds one regression test's recording — a git history of workspace
snapshots captured during a live pipeline run. Replay serves the pre-recorded
AI responses in place of real model calls, allowing the deterministic partition
of the pipeline (orchestrator logic, handler scripts, task management) to be
tested in isolation at zero token cost.

---

## Branch convention

- One **orphan branch** per regression test — no shared history with any other
  branch
- **Branch name = test name** (matches the directory under `tests/regression/`)
- `main` holds only this README; it is not a recording branch
- History is replaced on every re-record — each branch always contains exactly
  one run's worth of commits (no stacked runs, no duplicate invocation numbers)

---

## Regression tests

| Test | Commits | What it exercises |
|------|---------|-------------------|
| user-service | [user-service](https://github.com/cornjacket/ai-builder-recordings/commits/user-service/) | TM single-level decomposition — service decomposed into 3 components |
| platform-monolith | [platform-monolith](https://github.com/cornjacket/ai-builder-recordings/commits/platform-monolith/) | TM two-level decomposition — IAM + metrics services in one monolith |

When adding a new replay regression, add a row to this table. Branch link
format: `https://github.com/cornjacket/ai-builder-recordings/commits/<test-name>/`

---

## How to read a recording

Each commit in a recording branch corresponds to one orchestrator invocation:

```
d7131c5  recording.json                               ← manifest (always last)
521dd39  inv-21 pipeline done                         ← pipeline teardown commit
5c5c853  inv-20 LEAF_COMPLETE_HANDLER HANDLER_ALL_DONE
34f1486  inv-19 TESTER TESTER_TESTS_PASS
...
4b4470f  inv-01 ARCHITECT ARCHITECT_DECOMPOSITION_READY  ← root commit
```

Commit message format: `inv-NN ROLE OUTCOME`

- `inv-NN` — invocation number, 1-indexed, matching `recording.json`
- `ROLE` — the role that ran (e.g. `ARCHITECT`, `IMPLEMENTOR`, `TESTER`,
  `LEAF_COMPLETE_HANDLER`)
- `OUTCOME` — the outcome string from the state machine

Each commit snapshots `output/` and `target/` as they existed after that
invocation. AI invocations also write a response file to `responses/`:

```
responses/
    inv-01-ARCHITECT.txt      ← pre-recorded text served verbatim during replay
    inv-05-IMPLEMENTOR.txt
    ...
```

`recording.json` (the final commit) contains:

- `recorded_at` — ISO timestamp of the recording run
- `ai_builder_commit` — git SHA of the ai-builder repo at record time
- `task_hex_id` — 6-char hex ID of the top-level USER-TASK (pinned so replay
  produces identical task directory paths for snapshot comparison)
- `prompt_hashes` — SHA-256 of each role prompt file at record time (used by
  replay to detect prompt drift)
- `invocations` — ordered list of `{n, role, outcome, commit, ai}` entries

---

## Instructions for Claude

### Adding a new replay regression test

1. Create `tests/regression/<test-name>/record.sh` and `test-replay.sh`
   following the patterns in `tests/regression/user-service/`
2. Run `record.sh` — it initialises a fresh orphan branch, records the run,
   and pushes to this repo automatically
3. Run `test-replay.sh` once to confirm the recording replays cleanly
4. Add a row to the **Regression tests** table above:
   - Test name
   - Link: `https://github.com/cornjacket/ai-builder-recordings/commits/<test-name>/`
   - One-line description of what the test exercises
5. Commit and push the README update to `main`

### Re-recording an existing test

`record.sh --force` handles everything automatically:
- Wipes local `.git` (fresh history, no stacked runs)
- Re-records the pipeline
- Deletes the remote branch
- Pushes the new orphan branch

No README update needed — the table entry is unchanged.

---

## Further reading

- [`tests/regression/how-to-write-a-regression-test.md`](https://github.com/cornjacket/ai-builder/blob/main/tests/regression/how-to-write-a-regression-test.md) — how to add replay support to a new test
- [`tests/regression/README.md`](https://github.com/cornjacket/ai-builder/blob/main/tests/regression/README.md) — when to use replay vs live runs
- [`ai-builder/orchestrator/record-replay.md`](https://github.com/cornjacket/ai-builder/blob/main/ai-builder/orchestrator/record-replay.md) — orchestrator-level reference

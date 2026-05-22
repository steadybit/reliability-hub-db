# Automating Extension Documentation and Compatibility

Status: proposal — Phase 0 is shipped (PR #38); Phases 1–4 are pending.

## TL;DR

Every fact about a Steadybit action is currently authored in 3 places: the extension Go code, the `reliability-hub-db` YAML files, and the `docs-public` compatibility matrix. We just closed the first gap with a sync tool (Phase 0). This document proposes closing the second by making the compatibility matrix **derived from per-action capability requirements declared in the extension itself**, instead of hand-maintained across hundreds of markdown cells.

Required change: **one optional field (`Requires`) added to `ActionDescription` in `action-kit`**. Extensions adopt it at their own pace; the matrix gradually flips from hand-edited to generated.

---

## 1. The problem

Every extension's metadata is duplicated. The Go binary declares an action with its ID, label, description, parameters, and icon. The same data is hand-maintained in `reliability-hub-db`'s YAML files. The same parameters are hand-maintained a third time in `summary.mdx` markdown tables. The same compatibility information is hand-maintained a fourth time in `docs-public/quick-start/compatibility/README.md`.

When an extension author renames a parameter, updates a description, adds an action, or supports a new platform, four documents must be kept in sync by humans. They drift. Recent commits like `docs: sync action parameter tables with extension source` exist precisely because the system has no automation — they're periodic clean-ups of accumulated drift.

## 2. What we just shipped (Phase 0)

A Go sync tool in `reliability-hub-db` (PR #38) that closes the first gap: extension `--describe` JSON + a small hand-edited `sync.yml` registry → all `description.yml` files **and** the `# Parameters` table inside each action's `summary.mdx` are now derived, not authored.

- **One source of truth per kind**: action / target / advice schema comes from the extension binary; repo / GHCR / license / tags come from `sync.yml`.
- **CI-driven**: `.github/workflows/drift.yml` runs `docker pull <extension-image>; docker run --describe; sync`; when it detects drift, it opens a PR with the diff.
- **Author-curated content untouched**: `# Introduction`, `# Use Cases`, `previews.yml`, screenshots, YouTube embeds.

Concrete result: a successful first run against `extension-rabbitmq` produced a 13-file PR with real drift (longer descriptions, normalized tag order, the binary's newer multi-line SVG icon, parameter table updates). The same run a second time produces zero changes — idempotent.

## 3. What still drifts: the compatibility matrix

`docs-public/quick-start/compatibility/README.md` is ~30 KB of hand-aligned markdown tables: which actions work on Docker, Kubernetes flavors (EKS, GKE, AKS, Autopilot, OpenShift, minikube), AWS ECS on EC2, ECS on Fargate, Linux distros, Windows Server, etc., across all attack categories (Network, Resource, State, etc.). Adding a new extension or a new platform means a human edit to dozens or hundreds of cells.

Worse, the matrix is structurally derivable. A cell `(action, platform)` is ❌ for one of a small set of reasons: the action needs `tc-qdisc-netem` and Fargate doesn't load that kernel module; it needs `NET_ADMIN` and GKE Autopilot blocks privileged capabilities; it needs `kernel:nft` and only `iptables` is available. These constraints are owned by the *extension runtime* (what it tries to do) and the *platform* (what it permits) — not by the human editing the matrix.

## 4. Proposal: capability-based compatibility

Make the matrix derived, not authored.

### 4.1 One new field in `action-kit`

Add an **optional** `Supports` field to `ActionDescription` in `action-kit`. Each action declares the runtime primitives its implementation needs, per supported OS:

```go
action_kit_api.ActionDescription{
    Id:    "com.steadybit.extension_host.network_block_dns",
    Label: "Block DNS",
    Supports: []action_kit_api.SupportProfile{
        {OS: "linux",   Requires: []string{"cap:NET_ADMIN", "kernel:nft"}},
        {OS: "windows", Requires: []string{"win:windivert"}},
    },
}
```

This belongs in **`action-kit`**, not `extension-kit`. `ActionDescription` is owned by `action-kit`; `extension-kit` is the meta package that bundles the kits.

The field is **optional**. Extensions without it just don't appear in the auto-generated matrix sections (they fall back to manual matrix entries for now). This means **the change is non-breaking**: no extension repo has to move on day one.

The field flows through the existing pipeline:

1. Extension `--describe` includes `supports` in each action's JSON automatically — `--describe` just serializes whatever `ActionDescription` contains, no new code.
2. The sync tool's projector writes `supports` to `description.yml`. One-line addition to `internal/projector/action.go` in `reliability-hub-db`.
3. `description.yml` now carries the capability requirements.

### 4.2 The capability vocab

A small, namespaced vocab. Start with what explains every ❌ in the current matrix:

| Namespace   | Examples                                       | Meaning                                              |
| ----------- | ---------------------------------------------- | ---------------------------------------------------- |
| `os:`       | `os:linux`, `os:windows`                       | OS family the action implementation needs            |
| `cap:`      | `cap:NET_ADMIN`, `cap:SYS_PTRACE`              | Linux capabilities (POSIX `CAP_*`)                   |
| `kernel:`   | `kernel:tc-qdisc-netem`, `kernel:ebpf`, `kernel:nft` | Linux kernel features / modules                |
| `win:`      | `win:windivert`, `win:job-object-suspend`      | Windows-specific primitives                          |
| `runtime:`  | `runtime:host-netns-share`, `runtime:privileged` | Container runtime traits (OS-agnostic)             |

Expected size: ~10–15 Linux primitives + a handful of Windows + a handful of runtime traits. **Cloud-API actions** (Hibernate EC2, Stop ECS Service, etc.) don't need this at all — those sections of the matrix are organized by target type, which we already know from `description.yml`'s `targetType` field.

### 4.3 One new file in `docs-public`

`platforms.yml`: one entry per row of the current matrix, declaring what each deployment provides:

```yaml
platforms:
  - id: ecs-ec2
    label: AWS Elastic Container Service (ECS) on EC2
    os: linux
    provides:
      - cap:NET_ADMIN
      - kernel:tc-qdisc-netem
      - kernel:nft
      - kernel:ebpf
      - runtime:target-netns-attach
      - runtime:privileged

  - id: ecs-fargate
    label: AWS Elastic Container Service (ECS) on Fargate
    os: linux
    provides:
      - cap:NET_ADMIN
      - runtime:target-netns-attach
      # NOT provided: kernel:tc-qdisc-netem (kernel netem module not loaded)
      # NOT provided: runtime:privileged
      # NOT provided: kernel:ebpf

  - id: windows-server-2022
    label: Windows Server 2022 (x64)
    os: windows
    provides:
      - win:windivert
      - win:job-object-suspend
```

Roughly **15 entries** total — bounded by the number of deployment targets you publicly support, not by actions × platforms (the current matrix has hundreds of cells).

### 4.4 One small generator

Reads action capabilities from `reliability-hub-db` (via the existing submodule), reads platforms from local `platforms.yml`, computes each cell:

```
cell(action, platform) = ✅  if  any profile P in action.Supports satisfies:
                                    P.OS == platform.OS  AND  P.Requires ⊆ platform.Provides
                         = ❌  otherwise (with missing-capability footnote)
```

Emits the markdown. The footnote on each ❌ is auto-generated from "which required capability did the platform not provide":

```markdown
| Corrupt Outgoing Packages | ✅ | ✅ | ❌[^3] |

[^3]: Missing kernel:tc-qdisc-netem on ECS Fargate (netem module not loaded).
```

When AWS one day enables netem on Fargate, you change *one line* in `platforms.yml` and every affected ❌ flips to ✅ automatically.

## 5. Required changes by repo

| Repo                  | Change                                                                     | Effort                                        |
| --------------------- | -------------------------------------------------------------------------- | --------------------------------------------- |
| `action-kit`          | Add optional `Supports []SupportProfile` on `ActionDescription`            | ~30 minutes                                   |
| Each extension repo   | Bump `action-kit` dep; add `Supports` to actions; one-line per action      | ~1 hour per extension; **fully incremental**  |
| `reliability-hub-db`  | Add `supports` projection in `internal/projector/action.go`                | ~1 hour                                       |
| `docs-public`         | New `platforms.yml` + Go generator + CI job to regenerate the README       | ~1 day                                        |

**No agent change. No new infrastructure. No new repo.**

## 6. Phased rollout

Each phase delivers standalone value — we can stop at any phase if priorities shift.

| Phase | Scope                                                                       | Value delivered                                                                                   |
| ----- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| **0** ✅ | Sync tool for `description.yml` + `summary.mdx` parameter tables          | Extension → DB drift closed. Already shipped (PR #38).                                            |
| **1** | Define vocab; write `platforms.yml`; build generator; *manual* `supports.yml` overlay for one extension (no action-kit change yet) | End-to-end demo of generated matrix for one section. Vocab gets reviewed before propagating.       |
| **2** | `action-kit` change + annotate `extension-container` first (highest-impact)  | The Containers section of the matrix becomes generated. The Fargate gotchas become explicit.       |
| **3** | Roll out across extensions one-by-one                                       | Each annotated extension's rows move from hand-edited to generated. Independent per extension.    |
| **4** | Delete the hand-maintained matrix                                           | Single source of truth achieved end-to-end. CI prevents regressions.                              |

Phase 1 is deliberately structured to **NOT require any extension-side change yet** — we prove the generator works against a stand-in `supports.yml` file (just for the prototype) before asking every extension to adopt the new field.

## 7. What we're explicitly NOT doing

- **No CI smoke-test infrastructure.** A separate workstream worth doing later — running each action on each platform in CI would catch drift between declarations and reality. But it's not required for this proposal. Declarations are the source of truth; the matrix reflects what the extension authors say their code needs.
- **No agent-side capability probing.** Conceptually appealing (agent introspects its runtime, reports what it can do) but invasive. Static declarations + manual `platforms.yml` cover the same ground at a fraction of the cost.
- **No schema versioning.** `Supports` being optional means existing extensions keep working unchanged. No breaking change to `action-kit`.

## 8. Risks and mitigations

| Risk                                                                                          | Mitigation                                                                                                                                                                                                                  |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Capability vocab is too granular or too coarse                                                | Start with the ~15 items above (what explains every ❌ in the current matrix). Expand only when a real ❌ doesn't fit. Vocab is validated in CI against `platforms.yml`.                                                       |
| Author declares `Supports` incorrectly → silently wrong matrix                                | Same risk we live with today (matrix is already author-declared, just unstructured). Layer CI smoke tests on top later if needed. For now, treat declarations as authoritative.                                              |
| `platforms.yml` is still hand-maintained                                                      | Yes, but it's ~15 rows, not hundreds of cells. Bounded by the number of deployment targets we publicly support — that number changes rarely.                                                                                 |
| Rollout takes time across N extensions                                                        | Phased rollout means each annotated extension delivers value immediately. The matrix can stay partly manual for months while incrementally moving to generated.                                                              |
| `action-kit` is shared across many extensions; adding a field could cause churn                | Optional field, additive, non-breaking. Extensions on older `action-kit` versions keep working unchanged. Extensions opt in by bumping the dep when they're ready to declare `Supports`.                                     |

## 9. Why convince now

- **Phase 0 already paid for the infrastructure** (sync tool, CI workflow, projector). Phases 1–4 are incremental on top of that, not greenfield.
- **The matrix gets larger every quarter.** Adding the K-th cloud platform means K more hand-edits per action. The economics get worse, not better, with delay.
- **Customers see the matrix.** Errors there directly affect customer trust ("docs said this works on Fargate, it doesn't"). Removing humans from the loop reduces that class of error.
- **The same data is reusable elsewhere.** `Supports` enables an experiment-design check in the hub: "this action's target is on Fargate, but the action requires `kernel:tc-qdisc-netem` which Fargate doesn't provide — pick a different target". That's a customer-facing feature that becomes essentially free once the data exists.

## 10. Decision needed

1. Approve the proposal in principle.
2. Allocate Phase 1 (~1 day of engineering) — define the vocab + build the generator against a stand-in `supports.yml`. Phase 1 alone has standalone value (we get a structured `platforms.yml` and a generator, even if extensions haven't adopted `Supports` yet).
3. After Phase 1 demo: decide whether to commit to the `action-kit` change and the per-extension rollout.

If approved, Phase 1 is one PR in `docs-public` plus one short discussion to lock in the capability vocab.

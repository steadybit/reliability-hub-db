# Content Conventions

These rules keep the [Reliability Hub](https://hub.steadybit.com) consistent: easy to search, easy to scan, and the
same look no matter who added an entry. They apply to everything in this repository — extensions, target types,
actions, advice and templates. Check new and changed entries against them before opening a pull request.

## Tags

Tags drive the hub's search bar: every tag in this repository is offered as a search suggestion, and filtering by a
tag shows every entry whose tag *contains* the search term, ignoring case. A tag is only useful if everyone spells it
the same way.

1. **Reuse before you invent.** Look at the existing tags first and pick the established spelling:
   ```shell
   find . -name description.yml -not -path './node_modules/*' -exec awk 'FNR==1{t=0} /^tags:/{t=1;next} t&&/^[[:space:]]*- /{sub(/^[[:space:]]*- +/,"");print;next} {t=0}' {} + | sort | uniq -c | sort -rn
   jq -r '.tags[]?' templates/*/template.json | sort | uniq -c | sort -rn
   ```
   Only add a new tag for a concept no existing tag covers.
2. **One spelling per concept.** Never introduce a variant that differs only in case, plural, abbreviation or wording
   from an existing tag. Known pitfalls:

   | Use               | Not                                |
   |-------------------|------------------------------------|
   | `HTTP`            | `Http`, `HTTP Route`               |
   | `Network`         | `Networking`                       |
   | `Database`        | `DB`                               |
   | `Load Balancing`  | `Load balancing`                   |
   | `Virtual Machine` | `Virtual machine`, `Virtual machines` |

3. **Title Case**, except for names with an official spelling (`gRPC`, `k6`, `PromQL`, `.NET`, `NGINX`).
4. **Singular nouns**: `Container`, `Virtual Machine`, `Database`.
5. **Prefix cloud services with their provider**: `AWS Lambda`, `Azure Service Bus`, `GCP Cloud SQL`. Also add the
   provider tag itself (`AWS`, `Azure`, `GCP`).
6. **Tag the protocol an action works on.** An action that acts on HTTP(s) traffic gets `HTTP` — also when it selects
   traffic by hostname, and also for gRPC (which runs over HTTP/2; add `gRPC` as well). Network-level actions get
   `Network`.
7. **Keep variants in sync.** Actions that exist for several targets (e.g. `extension_container.*`,
   `extension_host.*`, `extension_host_windows.*`) carry the same tags, apart from the target-specific one
   (`Container`, `Host`, `Windows`).
8. **No duplicates** within one entry, and no tag that only repeats the label.

## Labels

- **Title Case**: `Reset TCP/HTTP(s) Connection`, not `Reset TCP/HTTP(s) connection`.
- Target type labels have a singular (`one`) and a plural (`other`) form.
- Labels match the label in the extension's code (the platform shows that one). Change both together.

## Descriptions

The `description` is the one line shown at the top of an entry and in search results. Its only job is to say what the
entry does.

- **One sentence, at most 150 characters** — aim for about 100. Most descriptions are around 50.
- Start with what it does: `Adds latency to calls to a specific HTTP(s) dependency, selected by hostname.`
- **No details**: parameters, modes, prerequisites, caveats, how it is implemented and how it compares to other actions
  belong in `summary.mdx`.
- Plain text on a single line: no Markdown, no line breaks.
- If the extension's code has a longer description, shorten it there too, so the hub and the platform match.

## Summaries (`summary.mdx`)

The summary holds everything the description leaves out. For actions, use these top-level sections in this order and
leave out the ones that don't apply:

```mdx
# Introduction
# Use Cases
# Prerequisites
# Parameters
# Details
# Rollback
```

For the extension summary structure, see [`.claude/CLAUDE.md`](.claude/CLAUDE.md#summarymdx).

## Icons

- Inline SVG in the YAML (pipe `|` multiline syntax).
- **No fixed `width` or `height`** on the `<svg>` element; keep the `viewBox` so the icon can scale.
- **Use `currentColor`** for `fill` and `stroke`, so the hub can color the icon with CSS.
- Use the same icon for an extension, its target types and its actions.

## IDs and Metadata

- IDs use reverse domain notation (`com.steadybit.extension_<name>.<target>.<action>`), are identical to the directory
  name, and match the IDs in the extension's code.
- `releaseDate` is a quoted ISO date: `'YYYY-MM-DD'`.
- When you rename or remove an ID, check whether the hub needs a redirect (`next.config.js` in
  [reliability-hub](https://github.com/steadybit/reliability-hub)).

## Review Checklist

- [ ] Tags reuse existing spellings and follow the tag rules above
- [ ] Protocol tags (`HTTP`, `Network`, …) are set where they apply
- [ ] Variants of the same action (container/host/Windows) carry the same tags
- [ ] Labels are Title Case and match the extension's code
- [ ] The description is one sentence of at most 150 characters; details are in `summary.mdx`
- [ ] Icons have no fixed size and use `currentColor`

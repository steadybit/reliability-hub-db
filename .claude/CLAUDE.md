# Reliability Hub DB

This repository contains the metadata for all Steadybit extensions, actions, target types, advice, and experiment templates displayed on the [Reliability Hub](https://hub.steadybit.com/).

## Repository Structure

```
reliability-hub-db/
├── extensions/          # Extension metadata (one dir per extension)
├── actions/             # Action definitions (one dir per action)
├── targetTypes/         # Target type definitions (one dir per target type)
├── advice/              # Advice definitions (one dir per advice)
├── templates/           # Experiment templates (one dir per template)
├── maintainers/         # Maintainer metadata
└── index.json           # Master index (auto-generated, do not edit manually)
```

## How to Register a New Extension

When creating a new Steadybit extension, you need to add entries here so it appears on the Reliability Hub. Create the following:

### 1. Extension (`extensions/<id>/`)

**Directory name**: `com.steadybit.extension_<name>` (reverse domain notation)

**Files**:
- `description.yml` (required)
- `summary.mdx` (required)
- `previews.yml` (optional — YouTube videos or screenshots)
- `public/` directory (optional — preview images)

#### `description.yml`

```yaml
---
id: com.steadybit.extension_<name>
label: <Display Name>
description: <One-line description>
icon: |
  <svg>...</svg>
maintainer: com.steadybit
license: MIT
gitHub:
  owner: steadybit
  repository: extension-<name>
ghcr:
  owner: steadybit
  repository: extension-<name>
  package: extension-<name>
homepage: https://hub.steadybit.com/extension/com.steadybit.extension_<name>
installation: https://github.com/steadybit/extension-<name>#installation
changelog: https://github.com/steadybit/extension-<name>/blob/main/CHANGELOG.md
releaseDate: 'YYYY-MM-DD'
tags:
  - <Tag1>
  - <Tag2>
```

#### `summary.mdx`

Use this structure:

```mdx
# Introduction to the <Name> Extension
<One paragraph describing what the extension does.>

# Integration and Functionality
<How it integrates, what libraries/APIs it uses, what permissions it needs.>

## <Subsection for each capability group>
- Bullet points of what it can do

# Installation and Setup
To integrate the <Name> extension with your environment, follow our [setup guide](https://github.com/steadybit/extension-<name>#installation).
```

### 2. Target Type (`targetTypes/<id>/`)

**Directory name**: `com.steadybit.extension_<name>.<target>` (matches the `TargetType` constant in Go code)

#### `description.yml`

```yaml
---
id: com.steadybit.extension_<name>.<target>
label:
  one: <Singular Label>
  other: <Plural Label>
icon: |
  <svg>...</svg>
extension: com.steadybit.extension_<name>
releaseDate: 'YYYY-MM-DD'
tags:
  - <Tag1>
  - <Tag2>
```

### 3. Actions (`actions/<id>/`)

**Directory name**: matches the action ID from Go code (e.g., `com.steadybit.extension_<name>.<target>.stop`)

One directory per action. Each action needs a `description.yml`.

#### `description.yml`

```yaml
---
id: com.steadybit.extension_<name>.<target>.<action>
label: <Action Label>
description: <One-line description>
icon: |
  <svg>...</svg>
kind: attack    # or "check"
targetType: com.steadybit.extension_<name>.<target>
extension: com.steadybit.extension_<name>
releaseDate: 'YYYY-MM-DD'
tags:
  - <Tag1>
  - <Tag2>
```

**`kind` values**: `attack` (for actions that modify state) or `check` (for monitoring/validation actions).

### 4. Advice (`advice/<id>/`) — Optional

Only needed if your extension provides advice (operational guidance).

#### `description.yml`

```yaml
---
id: com.steadybit.extension_<name>.advice.<advice-name>
label: <Advice Label>
description: <One-line description>
icon: |
  <svg>...</svg>
targetTypes:
  - com.steadybit.extension_<name>.<target1>
  - com.steadybit.extension_<name>.<target2>
tags:
  - <Tag1>
extension: com.steadybit.extension_<name>
releaseDate: 'YYYY-MM-DD'
```

## Conventions

- **IDs** use reverse domain notation: `com.steadybit.extension_<name>.<target>.<action>`
- **Icons** are inline SVG in the YAML, using the pipe (`|`) multiline syntax
- **Tags** are an array of strings, typically the technology name + category
- **`releaseDate`** is a quoted ISO date string: `'YYYY-MM-DD'`
- All extensions by Steadybit use `maintainer: com.steadybit` and `license: MIT`
- The SVG icon should be the same across the extension, target type, and all actions

## Checklist for Adding a New Extension

- [ ] `extensions/com.steadybit.extension_<name>/description.yml`
- [ ] `extensions/com.steadybit.extension_<name>/summary.mdx`
- [ ] `targetTypes/com.steadybit.extension_<name>.<target>/description.yml` (one per target type)
- [ ] `actions/com.steadybit.extension_<name>.<target>.<action>/description.yml` (one per action)
- [ ] `advice/com.steadybit.extension_<name>.advice.<name>/description.yml` (if applicable)
- [ ] IDs match what's in the extension's Go code (`TargetType`, action IDs, advice IDs)
- [ ] Create a PR on this repo

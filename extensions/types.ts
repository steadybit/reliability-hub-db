/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

import { MaintainerId } from "../maintainers/types";

// reverse domain name, e.g., com.example
export type ExtensionId = string;

// SPDX license identifier
// https://spdx.org/licenses/
export type SpdxLicenseIdentifier = string;

export interface GitHub {
  owner: string;
  repository: string;
}

export interface DockerHub {
  owner: string;
  repository: string;
}

export interface GitHubContainerRegistry {
  owner: string;
  repository: string;
  package: string;
}

export interface ExtensionDescription {
  id: ExtensionId;
  label: string;
  description: string;
  icon?: string | null;
  license?: SpdxLicenseIdentifier;
  maintainer: MaintainerId;
  // will result in GitHub statistics presentation
  gitHub?: GitHub;
  // will result in Docker statistics presentation
  dockerHub?: DockerHub;
  // will result in Docker statistics presentation
  ghcr?: GitHubContainerRegistry;
  homepage?: string;
  installation?: string;
  changelog?: string;
  tags: string[];
  releaseDate?: string;
}

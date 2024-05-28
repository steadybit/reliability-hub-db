/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

import { MaintainerId } from '../maintainers/types';
import { DockerHub, GitHub, GitHubContainerRegistry } from '../extensions/types';

// reverse domain name, e.g., com.example
export type ExtensionId = string;

// SPDX license identifier
// https://spdx.org/licenses/
export type SpdxLicenseIdentifier = string;

export interface TemplateDescription {
	license?: SpdxLicenseIdentifier;
	maintainer: MaintainerId;
	// will result in GitHub statistics presentation
	gitHub?: GitHub;
	// will result in Docker statistics presentation
	dockerHub?: DockerHub;
	// will result in Docker statistics presentation
	ghcr?: GitHubContainerRegistry;
	homepage?: string;
}

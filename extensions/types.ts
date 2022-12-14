import { MaintainerId } from '../maintainers/types';

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
	homepage?: string;
	changelog?: string;
	tags: string[];
}

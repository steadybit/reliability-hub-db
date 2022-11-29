// reverse domain name, e.g., com.example
export type MaintainerId = string;

export interface MaintainerDescription {
	id: MaintainerId;
	name: string;
	email: string;
	icon?: string;
	website?: string;
	avatar: string;
}

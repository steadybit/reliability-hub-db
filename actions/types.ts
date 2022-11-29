import { ActionDescription as ActionKitActionDescription } from '@steadybit/action-kit-api';

import { ExtensionId } from '../extensions/types';

export type Deprecation = ReplacementDeprecation;

export interface ReplacementDeprecation {
	type: 'replacement';
	newActionId: string;
}

export function isReplacementDeprecation(d: Deprecation): d is ReplacementDeprecation {
	return d.type === 'replacement';
}

export interface ActionDescription
	extends Pick<
		ActionKitActionDescription,
		'id' | 'label' | 'icon' | 'kind' | 'category' | 'targetType' | 'description'
	> {
	extension: ExtensionId;
	deprecation?: Deprecation;
	promotedActions: string[];
}

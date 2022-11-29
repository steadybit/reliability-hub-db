import { TargetDescription as DiscoveryKitTargetDescription } from '@steadybit/discovery-kit-api';

import { ExtensionId } from '../extensions/types';

export interface TargetTypeDescription
	extends Pick<DiscoveryKitTargetDescription, 'id' | 'label' | 'icon' | 'category'> {
	description?: string;
	extension: ExtensionId;
}

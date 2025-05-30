/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

import { ExtensionId } from "../extensions/types";

export interface TargetTypeDescription {
  id: string;
  label: { one: string; other: string };
  icon?: string;
  category?: string | null;
  description?: string;
  extension: ExtensionId;
  releaseDate?: string;
  beta?: boolean;
}

/*
 * Copyright 2024 steadybit GmbH. All rights reserved.
 */

export interface AdviceDescription {
  id: string;
  label: string;
  description: string;
  icon?: string;
  targetTypes: string[];
  tags: string[];
  extension: string;
  releaseDate: string | null;
}

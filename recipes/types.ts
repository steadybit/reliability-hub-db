/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

import { GitHub, SpdxLicenseIdentifier } from "../extensions/types";

import { MaintainerId } from "../maintainers/types";

export interface RecipeDescription {
  id: string;
  label: string;
  description: string;
  tags: string[];
  license?: SpdxLicenseIdentifier;
  maintainer: MaintainerId;
  // will result in GitHub statistics presentation
  gitHub?: GitHub;
  homepage?: string;
  releaseDate?: string;
}

export interface ExperimentRecipe extends Experiment {}

export interface Experiment {
  lanes: Lane[];
}

export interface Lane {
  steps: Step[];
}

export interface Step {
  actionType?: string;
}

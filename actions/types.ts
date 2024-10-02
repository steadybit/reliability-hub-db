/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

import { ExtensionId } from "../extensions/types";

export type Deprecation = ReplacementDeprecation | FullDeprecation;

export interface ReplacementDeprecation {
  type: "replacement";
  newActionId: string;
}

export interface FullDeprecation {
  type: "deprecated";
}

export function isReplacementDeprecation(
  d: Deprecation,
): d is ReplacementDeprecation {
  return d.type === "replacement";
}

export function isFullDeprecation(d: Deprecation): d is FullDeprecation {
  return d.type === "deprecated";
}

export interface ActionDescription {
  id: string;
  label: string;
  icon: string;
  kind: "attack" | "check" | "load_test" | "other";
  category: string | null;
  targetType?: string | null;
  description: string;
  extension: ExtensionId;
  deprecation?: Deprecation;
  promotedActions: string[];
  releaseDate?: string;
  tags: string[];
}

/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

import { ActionDescription as ActionKitActionDescription } from "@steadybit/action-kit-api";

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

export interface ActionDescription
  extends Pick<
    ActionKitActionDescription,
    "id" | "label" | "icon" | "kind" | "category" | "targetType" | "description"
  > {
  extension: ExtensionId;
  deprecation?: Deprecation;
  promotedActions: string[];
  releaseDate?: string;
}

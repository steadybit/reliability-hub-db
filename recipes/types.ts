export interface RecipeDescription {
	id: string;
	label: string;
	description: string;
	tags: string[];
}

export interface ExperimentRecipe extends Experiment {
}

export interface Experiment {
	lanes: Lane[];
}

export interface Lane {
	steps: Step[];
}

export interface Step {
	actionType?: string;
}

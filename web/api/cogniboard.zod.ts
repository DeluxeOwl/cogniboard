/**
 * Generated by orval v7.5.0 🍺
 * Do not edit manually.
 * CogniBoard
 * OpenAPI spec version: 0.0.1
 */
import { z as zod } from "zod";

/**
 * @summary Get all tasks
 */
export const getTasksResponse = zod.object({
	$schema: zod.string().url().optional(),
	tasks: zod
		.array(
			zod.object({
				assignee: zod.string().nullable(),
				completed_at: zod.string().datetime().nullable(),
				created_at: zod.string().datetime(),
				description: zod.string().nullable(),
				due_date: zod.string().datetime().nullable(),
				id: zod.string(),
				status: zod.string(),
				title: zod.string(),
			}),
		)
		.nullable(),
});

/**
 * @summary Create a task
 */
export const createTaskBodyTitleMax = 50;

export const createTaskBody = zod.object({
	assignee_name: zod.string().optional(),
	description: zod.string().optional(),
	due_date: zod.string().datetime().optional(),
	title: zod.string().min(1).max(createTaskBodyTitleMax),
});

/**
 * @summary Assign a task to someone
 */
export const assignTaskParams = zod.object({
	taskId: zod.string(),
});

export const assignTaskBody = zod.object({
	assignee_name: zod.string().min(1),
});

/**
 * @summary Change task status
 */
export const changeTaskStatusParams = zod.object({
	taskId: zod.string(),
});

export const changeTaskStatusBody = zod.object({
	status: zod.string().min(1),
});

/**
 * @summary Unassign a task
 */
export const unassignTaskParams = zod.object({
	taskId: zod.string(),
});

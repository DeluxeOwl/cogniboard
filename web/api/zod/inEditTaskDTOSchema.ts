import { z } from 'zod'

export const inEditTaskDTOSchema = z.object({
  $schema: z.string().url().describe('A URL to the JSON Schema for this object.').optional(),
  assignee_name: z.string().describe('Name of the person to assign the task to').optional(),
  description: z.string().describe("Task's description").optional(),
  due_date: z.string().datetime().describe("Task's due date").optional(),
  status: z.string().describe("Task's status").optional(),
  title: z.string().max(50).describe("Task's title").optional(),
})
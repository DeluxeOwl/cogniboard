import { z } from 'zod'

export const changeTaskStatusSchema = z.object({
  $schema: z.string().url().describe('A URL to the JSON Schema for this object.').optional(),
  status: z.string().min(1).describe('New status for the task'),
})
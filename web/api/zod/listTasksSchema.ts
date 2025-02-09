import { taskSchema } from './taskSchema.ts'
import { z } from 'zod'

export const listTasksSchema = z.object({
  $schema: z.string().url().describe('A URL to the JSON Schema for this object.').optional(),
  tasks: z.array(z.lazy(() => taskSchema)).nullable(),
})
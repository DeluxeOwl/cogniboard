import { inTaskDTOSchema } from './inTaskDTOSchema.ts'
import { z } from 'zod'

export const inTasksDTOSchema = z.object({
  $schema: z.string().url().describe('A URL to the JSON Schema for this object.').optional(),
  tasks: z.array(z.lazy(() => inTaskDTOSchema)).nullable(),
})
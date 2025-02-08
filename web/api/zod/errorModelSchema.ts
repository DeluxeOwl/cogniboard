import { errorDetailSchema } from './errorDetailSchema.ts'
import { z } from 'zod'

export const errorModelSchema = z.object({
  $schema: z.string().url().describe('A URL to the JSON Schema for this object.').optional(),
  detail: z.string().describe('A human-readable explanation specific to this occurrence of the problem.').optional(),
  errors: z
    .array(z.lazy(() => errorDetailSchema))
    .describe('Optional list of individual error details')
    .nullable()
    .nullish(),
  instance: z.string().url().describe('A URI reference that identifies the specific occurrence of the problem.').optional(),
  status: z.number().int().describe('HTTP status code').optional(),
  title: z.string().describe('A short, human-readable summary of the problem type. This value should not change between occurrences of the error.').optional(),
  type: z.string().url().default('about:blank').describe('A URI reference to human-readable documentation for the error.'),
})
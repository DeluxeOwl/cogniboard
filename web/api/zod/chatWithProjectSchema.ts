import { messageSchema } from './messageSchema.ts'
import { z } from 'zod'

export const chatWithProjectSchema = z.object({
  $schema: z.string().url().describe('A URL to the JSON Schema for this object.').optional(),
  messages: z.array(z.lazy(() => messageSchema)).nullable(),
})
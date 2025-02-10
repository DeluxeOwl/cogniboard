import { contentSchema } from './contentSchema.ts'
import { z } from 'zod'

export const messageSchema = z.object({
  content: z.array(z.lazy(() => contentSchema)).nullable(),
  role: z.string(),
})
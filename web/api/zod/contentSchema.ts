import { z } from 'zod'

export const contentSchema = z.object({
  text: z.string(),
  type: z.string(),
})
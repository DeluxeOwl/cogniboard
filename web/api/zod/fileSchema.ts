import { z } from 'zod'

export const fileSchema = z.object({
  id: z.string(),
  mime_type: z.string(),
  name: z.string(),
  size: z.number().int(),
  uploaded_at: z.string().datetime(),
})
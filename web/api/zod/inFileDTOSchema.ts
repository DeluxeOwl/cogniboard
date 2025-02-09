import { z } from 'zod'

export const inFileDTOSchema = z.object({
  mime_type: z.string(),
  name: z.string(),
  size: z.number().int(),
  uploaded_at: z.string().datetime(),
})
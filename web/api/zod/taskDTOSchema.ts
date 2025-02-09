import { fileDTOSchema } from './fileDTOSchema.ts'
import { z } from 'zod'

export const taskDTOSchema = z.object({
  asignee: z.string().nullable(),
  completed_at: z.string().datetime().nullable(),
  created_at: z.string().datetime(),
  description: z.string().nullable(),
  due_date: z.string().datetime().nullable(),
  files: z.array(z.lazy(() => fileDTOSchema)).nullable(),
  id: z.string(),
  status: z.string(),
  title: z.string(),
  updated_at: z.string().datetime(),
})
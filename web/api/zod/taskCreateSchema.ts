import { errorModelSchema } from './errorModelSchema.ts'
import { z } from 'zod'

/**
 * @description No Content
 */
export const taskCreate204Schema = z.any()

/**
 * @description Error
 */
export const taskCreateErrorSchema = z.lazy(() => errorModelSchema)

export const taskCreateMutationRequestSchema = z.object({
  assignee_name: z.string().describe("Task's asignee (if any)").optional(),
  description: z.string().describe("Task's description").optional(),
  due_date: z.string().datetime().describe("Task's due date (if any)").optional(),
  files: z.array(z.instanceof(File)).optional(),
  title: z.string().min(1).max(50).describe("Task's name"),
})

export const taskCreateMutationResponseSchema = z.lazy(() => taskCreate204Schema)
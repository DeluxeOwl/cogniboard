import { errorModelSchema } from './errorModelSchema.ts'
import { z } from 'zod'

export const taskEditPathParamsSchema = z.object({
  taskId: z.string(),
})

/**
 * @description No Content
 */
export const taskEdit204Schema = z.any()

/**
 * @description Error
 */
export const taskEditErrorSchema = z.lazy(() => errorModelSchema)

export const taskEditMutationRequestSchema = z.object({
  assignee_name: z.string().describe("Task's asignee (if any)").optional(),
  description: z.string().describe("Task's description").optional(),
  due_date: z.string().datetime().describe("Task's due date (if any)").optional(),
  files: z.array(z.instanceof(File)).optional(),
  title: z.string().min(1).max(50).describe("Task's name"),
})

export const taskEditMutationResponseSchema = z.lazy(() => taskEdit204Schema)
import { changeTaskStatusSchema } from './changeTaskStatusSchema.ts'
import { errorModelSchema } from './errorModelSchema.ts'
import { z } from 'zod'

export const taskChangeStatusPathParamsSchema = z.object({
  taskId: z.string(),
})

/**
 * @description No Content
 */
export const taskChangeStatus204Schema = z.any()

/**
 * @description Error
 */
export const taskChangeStatusErrorSchema = z.lazy(() => errorModelSchema)

export const taskChangeStatusMutationRequestSchema = z.lazy(() => changeTaskStatusSchema).and(z.object({ $schema: z.never() }))

export const taskChangeStatusMutationResponseSchema = z.lazy(() => taskChangeStatus204Schema)
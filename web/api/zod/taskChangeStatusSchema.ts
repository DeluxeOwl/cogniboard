import { errorModelSchema } from './errorModelSchema.ts'
import { inChangeTaskStatusDTOSchema } from './inChangeTaskStatusDTOSchema.ts'
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

export const taskChangeStatusMutationRequestSchema = z.lazy(() => inChangeTaskStatusDTOSchema).and(z.object({ $schema: z.never() }))

export const taskChangeStatusMutationResponseSchema = z.lazy(() => taskChangeStatus204Schema)
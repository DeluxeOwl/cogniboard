import { errorModelSchema } from './errorModelSchema.ts'
import { inEditTaskDTOSchema } from './inEditTaskDTOSchema.ts'
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

export const taskEditMutationRequestSchema = z.lazy(() => inEditTaskDTOSchema).and(z.object({ $schema: z.never() }))

export const taskEditMutationResponseSchema = z.lazy(() => taskEdit204Schema)
import { chatWithProjectSchema } from './chatWithProjectSchema.ts'
import { errorModelSchema } from './errorModelSchema.ts'
import { z } from 'zod'

/**
 * @description OK
 */
export const projectChat200Schema = z.any()

/**
 * @description Error
 */
export const projectChatErrorSchema = z.lazy(() => errorModelSchema)

export const projectChatMutationRequestSchema = z.lazy(() => chatWithProjectSchema).and(z.object({ $schema: z.never() }))

export const projectChatMutationResponseSchema = z.lazy(() => projectChat200Schema)
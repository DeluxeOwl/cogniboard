import { errorModelSchema } from './errorModelSchema.ts'
import { inTasksDTOSchema } from './inTasksDTOSchema.ts'
import { z } from 'zod'

/**
 * @description OK
 */
export const tasks200Schema = z.lazy(() => inTasksDTOSchema)

/**
 * @description Error
 */
export const tasksErrorSchema = z.lazy(() => errorModelSchema)

export const tasksQueryResponseSchema = z.lazy(() => tasks200Schema)
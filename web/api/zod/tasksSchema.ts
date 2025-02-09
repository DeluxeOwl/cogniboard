import { errorModelSchema } from './errorModelSchema.ts'
import { listTasksSchema } from './listTasksSchema.ts'
import { z } from 'zod'

/**
 * @description OK
 */
export const tasks200Schema = z.lazy(() => listTasksSchema)

/**
 * @description Error
 */
export const tasksErrorSchema = z.lazy(() => errorModelSchema)

export const tasksQueryResponseSchema = z.lazy(() => tasks200Schema)
import type { ErrorModel } from './ErrorModel.ts'
import type { InTasksDTO } from './InTasksDTO.ts'

/**
 * @description OK
 */
export type Tasks200 = InTasksDTO

/**
 * @description Error
 */
export type TasksError = ErrorModel

export type TasksQueryResponse = Tasks200

export type TasksQuery = {
  Response: Tasks200
  Errors: any
}
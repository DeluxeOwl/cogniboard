import type { ErrorModel } from './ErrorModel.ts'
import type { ListTasks } from './ListTasks.ts'

/**
 * @description OK
 */
export type Tasks200 = ListTasks

/**
 * @description Error
 */
export type TasksError = ErrorModel

export type TasksQueryResponse = Tasks200

export type TasksQuery = {
  Response: Tasks200
  Errors: any
}
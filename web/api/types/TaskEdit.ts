import type { ErrorModel } from './ErrorModel.ts'
import type { InEditTaskDTO } from './InEditTaskDTO.ts'

export type TaskEditPathParams = {
  /**
   * @type string
   */
  taskId: string
}

/**
 * @description No Content
 */
export type TaskEdit204 = any

/**
 * @description Error
 */
export type TaskEditError = ErrorModel

export type TaskEditMutationRequest = Omit<NonNullable<InEditTaskDTO>, '$schema'>

export type TaskEditMutationResponse = TaskEdit204

export type TaskEditMutation = {
  Response: TaskEdit204
  Request: TaskEditMutationRequest
  PathParams: TaskEditPathParams
  Errors: any
}
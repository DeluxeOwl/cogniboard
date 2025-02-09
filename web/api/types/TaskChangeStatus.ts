import type { ChangeTaskStatus } from './ChangeTaskStatus.ts'
import type { ErrorModel } from './ErrorModel.ts'

export type TaskChangeStatusPathParams = {
  /**
   * @type string
   */
  taskId: string
}

/**
 * @description No Content
 */
export type TaskChangeStatus204 = any

/**
 * @description Error
 */
export type TaskChangeStatusError = ErrorModel

export type TaskChangeStatusMutationRequest = Omit<NonNullable<ChangeTaskStatus>, '$schema'>

export type TaskChangeStatusMutationResponse = TaskChangeStatus204

export type TaskChangeStatusMutation = {
  Response: TaskChangeStatus204
  Request: TaskChangeStatusMutationRequest
  PathParams: TaskChangeStatusPathParams
  Errors: any
}
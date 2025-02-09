import type { ErrorModel } from './ErrorModel.ts'

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

export type TaskEditMutationRequest = {
  /**
   * @description Task\'s asignee (if any)
   * @type string | undefined
   */
  assignee_name?: string
  /**
   * @description Task\'s description
   * @type string | undefined
   */
  description?: string
  /**
   * @description Task\'s due date (if any)
   * @type string | undefined, date-time
   */
  due_date?: string
  /**
   * @type array | undefined
   */
  files?: Blob[]
  /**
   * @description Task\'s name
   * @minLength 1
   * @maxLength 50
   * @type string
   */
  title: string
}

export type TaskEditMutationResponse = TaskEdit204

export type TaskEditMutation = {
  Response: TaskEdit204
  Request: TaskEditMutationRequest
  PathParams: TaskEditPathParams
  Errors: any
}
import type { ErrorModel } from './ErrorModel.ts'

/**
 * @description No Content
 */
export type TaskCreate204 = any

/**
 * @description Error
 */
export type TaskCreateError = ErrorModel

export type TaskCreateMutationRequest = {
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

export type TaskCreateMutationResponse = TaskCreate204

export type TaskCreateMutation = {
  Response: TaskCreate204
  Request: TaskCreateMutationRequest
  Errors: any
}
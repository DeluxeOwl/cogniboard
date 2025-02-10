import type { ChatWithProject } from './ChatWithProject.ts'
import type { ErrorModel } from './ErrorModel.ts'

/**
 * @description OK
 */
export type ProjectChat200 = any

/**
 * @description Error
 */
export type ProjectChatError = ErrorModel

export type ProjectChatMutationRequest = Omit<NonNullable<ChatWithProject>, '$schema'>

export type ProjectChatMutationResponse = ProjectChat200

export type ProjectChatMutation = {
  Response: ProjectChat200
  Request: ProjectChatMutationRequest
  Errors: any
}
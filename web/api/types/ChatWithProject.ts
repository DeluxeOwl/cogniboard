import type { Message } from './Message.ts'

export type ChatWithProject = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @type array
   */
  messages: Message[] | null
}
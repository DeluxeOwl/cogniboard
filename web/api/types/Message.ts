import type { Content } from './Content.ts'

export type Message = {
  /**
   * @type array
   */
  content: Content[] | null
  /**
   * @type string
   */
  role: string
}
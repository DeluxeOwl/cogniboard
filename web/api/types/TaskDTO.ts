import type { FileDTO } from './FileDTO.ts'

export type TaskDTO = {
  /**
   * @type string
   */
  asignee: string | null
  /**
   * @type string, date-time
   */
  completed_at: string | null
  /**
   * @type string, date-time
   */
  created_at: string
  /**
   * @type string
   */
  description: string | null
  /**
   * @type string, date-time
   */
  due_date: string | null
  /**
   * @type array
   */
  files: FileDTO[] | null
  /**
   * @type string
   */
  id: string
  /**
   * @type string
   */
  status: string
  /**
   * @type string
   */
  title: string
  /**
   * @type string, date-time
   */
  updated_at: string
}
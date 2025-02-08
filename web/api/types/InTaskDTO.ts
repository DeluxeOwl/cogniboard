export type InTaskDTO = {
  /**
   * @type string
   */
  assignee: string | null
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
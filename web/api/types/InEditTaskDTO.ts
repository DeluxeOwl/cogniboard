export type InEditTaskDTO = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @description Name of the person to assign the task to
   * @type string | undefined
   */
  assignee_name?: string
  /**
   * @description Task\'s description
   * @type string | undefined
   */
  description?: string
  /**
   * @description Task\'s due date
   * @type string | undefined, date-time
   */
  due_date?: string
  /**
   * @description Task\'s status
   * @type string | undefined
   */
  status?: string
  /**
   * @description Task\'s title
   * @maxLength 50
   * @type string | undefined
   */
  title?: string
}
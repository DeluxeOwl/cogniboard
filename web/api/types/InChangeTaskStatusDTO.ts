export type InChangeTaskStatusDTO = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @description New status for the task
   * @minLength 1
   * @type string
   */
  status: string
}
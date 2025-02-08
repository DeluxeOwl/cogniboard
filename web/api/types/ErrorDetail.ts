export type ErrorDetail = {
  /**
   * @description Where the error occurred, e.g. \'body.items[3].tags\' or \'path.thing-id\'
   * @type string | undefined
   */
  location?: string
  /**
   * @description Error message text
   * @type string | undefined
   */
  message?: string
  value?: any
}
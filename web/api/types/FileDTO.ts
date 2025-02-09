export type FileDTO = {
  /**
   * @type string
   */
  id: string
  /**
   * @type string
   */
  mime_type: string
  /**
   * @type string
   */
  name: string
  /**
   * @type integer, int64
   */
  size: number
  /**
   * @type string, date-time
   */
  uploaded_at: string
}
import type { ErrorDetail } from './ErrorDetail.ts'

export type ErrorModel = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @description A human-readable explanation specific to this occurrence of the problem.
   * @type string | undefined
   */
  detail?: string
  /**
   * @description Optional list of individual error details
   * @type array
   */
  errors?: ErrorDetail[] | null
  /**
   * @description A URI reference that identifies the specific occurrence of the problem.
   * @type string | undefined, uri
   */
  instance?: string
  /**
   * @description HTTP status code
   * @type integer | undefined, int64
   */
  status?: number
  /**
   * @description A short, human-readable summary of the problem type. This value should not change between occurrences of the error.
   * @type string | undefined
   */
  title?: string
  /**
   * @description A URI reference to human-readable documentation for the error.
   * @default "about:blank"
   * @type string | undefined, uri
   */
  type?: string
}
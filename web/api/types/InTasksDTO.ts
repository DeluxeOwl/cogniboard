import type { InTaskDTO } from './InTaskDTO.ts'

export type InTasksDTO = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @type array
   */
  tasks: InTaskDTO[] | null
}
import type { TaskDTO } from './TaskDTO.ts'

export type ListTasks = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @type array
   */
  tasks: TaskDTO[] | null
}
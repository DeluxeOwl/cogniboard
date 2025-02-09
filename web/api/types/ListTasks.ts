import type { Task } from './Task.ts'

export type ListTasks = {
  /**
   * @description A URL to the JSON Schema for this object.
   * @type string | undefined, uri
   */
  readonly $schema?: string
  /**
   * @type array
   */
  tasks: Task[] | null
}
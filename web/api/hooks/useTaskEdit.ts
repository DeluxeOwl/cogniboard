import client from '@kubb/plugin-client/clients/axios'
import type { TaskEditMutationRequest, TaskEditMutationResponse, TaskEditPathParams } from '../types/TaskEdit.ts'
import type { RequestConfig, ResponseErrorConfig } from '@kubb/plugin-client/clients/axios'
import type { UseMutationOptions } from '@tanstack/react-query'
import { useMutation } from '@tanstack/react-query'

export const taskEditMutationKey = () => [{ url: '/tasks/{taskId}/edit' }] as const

export type TaskEditMutationKey = ReturnType<typeof taskEditMutationKey>

/**
 * @summary Edit a task
 * {@link /tasks/:taskId/edit}
 */
export async function taskEdit(
  taskId: TaskEditPathParams['taskId'],
  data?: TaskEditMutationRequest,
  config: Partial<RequestConfig<TaskEditMutationRequest>> & { client?: typeof client } = {},
) {
  const { client: request = client, ...requestConfig } = config

  const res = await request<TaskEditMutationResponse, ResponseErrorConfig<Error>, TaskEditMutationRequest>({
    method: 'POST',
    url: `/tasks/${taskId}/edit`,
    baseURL: 'http://127.0.0.1:8888/v1/api',
    data,
    ...requestConfig,
  })
  return res.data
}

/**
 * @summary Edit a task
 * {@link /tasks/:taskId/edit}
 */
export function useTaskEdit(
  options: {
    mutation?: UseMutationOptions<
      TaskEditMutationResponse,
      ResponseErrorConfig<Error>,
      { taskId: TaskEditPathParams['taskId']; data?: TaskEditMutationRequest }
    >
    client?: Partial<RequestConfig<TaskEditMutationRequest>> & { client?: typeof client }
  } = {},
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? taskEditMutationKey()

  return useMutation<TaskEditMutationResponse, ResponseErrorConfig<Error>, { taskId: TaskEditPathParams['taskId']; data?: TaskEditMutationRequest }>({
    mutationFn: async ({ taskId, data }) => {
      return taskEdit(taskId, data, config)
    },
    mutationKey,
    ...mutationOptions,
  })
}
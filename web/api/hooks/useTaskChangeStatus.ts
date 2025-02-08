import client from '@kubb/plugin-client/clients/axios'
import type { TaskChangeStatusMutationRequest, TaskChangeStatusMutationResponse, TaskChangeStatusPathParams } from '../types/TaskChangeStatus.ts'
import type { RequestConfig, ResponseErrorConfig } from '@kubb/plugin-client/clients/axios'
import type { UseMutationOptions } from '@tanstack/react-query'
import { useMutation } from '@tanstack/react-query'

export const taskChangeStatusMutationKey = () => [{ url: '/tasks/{taskId}/status' }] as const

export type TaskChangeStatusMutationKey = ReturnType<typeof taskChangeStatusMutationKey>

/**
 * @summary Change task status
 * {@link /tasks/:taskId/status}
 */
export async function taskChangeStatus(
  taskId: TaskChangeStatusPathParams['taskId'],
  data: TaskChangeStatusMutationRequest,
  config: Partial<RequestConfig<TaskChangeStatusMutationRequest>> & { client?: typeof client } = {},
) {
  const { client: request = client, ...requestConfig } = config

  const res = await request<TaskChangeStatusMutationResponse, ResponseErrorConfig<Error>, TaskChangeStatusMutationRequest>({
    method: 'POST',
    url: `/tasks/${taskId}/status`,
    baseURL: 'http://127.0.0.1:8888/v1/api',
    data,
    ...requestConfig,
  })
  return res.data
}

/**
 * @summary Change task status
 * {@link /tasks/:taskId/status}
 */
export function useTaskChangeStatus(
  options: {
    mutation?: UseMutationOptions<
      TaskChangeStatusMutationResponse,
      ResponseErrorConfig<Error>,
      { taskId: TaskChangeStatusPathParams['taskId']; data: TaskChangeStatusMutationRequest }
    >
    client?: Partial<RequestConfig<TaskChangeStatusMutationRequest>> & { client?: typeof client }
  } = {},
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? taskChangeStatusMutationKey()

  return useMutation<
    TaskChangeStatusMutationResponse,
    ResponseErrorConfig<Error>,
    { taskId: TaskChangeStatusPathParams['taskId']; data: TaskChangeStatusMutationRequest }
  >({
    mutationFn: async ({ taskId, data }) => {
      return taskChangeStatus(taskId, data, config)
    },
    mutationKey,
    ...mutationOptions,
  })
}
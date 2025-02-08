import client from '@kubb/plugin-client/clients/axios'
import type { TasksQueryResponse } from '../types/Tasks.ts'
import type { RequestConfig, ResponseErrorConfig } from '@kubb/plugin-client/clients/axios'
import type { QueryKey, UseSuspenseQueryOptions, UseSuspenseQueryResult } from '@tanstack/react-query'
import { queryOptions, useSuspenseQuery } from '@tanstack/react-query'

export const tasksSuspenseQueryKey = () => [{ url: '/tasks' }] as const

export type TasksSuspenseQueryKey = ReturnType<typeof tasksSuspenseQueryKey>

/**
 * @summary Get all tasks
 * {@link /tasks}
 */
export async function tasksSuspense(config: Partial<RequestConfig> & { client?: typeof client } = {}) {
  const { client: request = client, ...requestConfig } = config

  const res = await request<TasksQueryResponse, ResponseErrorConfig<Error>, unknown>({
    method: 'GET',
    url: `/tasks`,
    baseURL: 'http://127.0.0.1:8888/v1/api',
    ...requestConfig,
  })
  return res.data
}

export function tasksSuspenseQueryOptions(config: Partial<RequestConfig> & { client?: typeof client } = {}) {
  const queryKey = tasksSuspenseQueryKey()
  return queryOptions<TasksQueryResponse, ResponseErrorConfig<Error>, TasksQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      return tasksSuspense(config)
    },
  })
}

/**
 * @summary Get all tasks
 * {@link /tasks}
 */
export function useTasksSuspense<TData = TasksQueryResponse, TQueryData = TasksQueryResponse, TQueryKey extends QueryKey = TasksSuspenseQueryKey>(
  options: {
    query?: Partial<UseSuspenseQueryOptions<TasksQueryResponse, ResponseErrorConfig<Error>, TData, TQueryKey>>
    client?: Partial<RequestConfig> & { client?: typeof client }
  } = {},
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? tasksSuspenseQueryKey()

  const query = useSuspenseQuery({
    ...(tasksSuspenseQueryOptions(config) as unknown as UseSuspenseQueryOptions),
    queryKey,
    ...(queryOptions as unknown as Omit<UseSuspenseQueryOptions, 'queryKey'>),
  }) as UseSuspenseQueryResult<TData, ResponseErrorConfig<Error>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}
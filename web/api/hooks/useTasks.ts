import client from '@kubb/plugin-client/clients/axios'
import type { TasksQueryResponse } from '../types/Tasks.ts'
import type { RequestConfig, ResponseErrorConfig } from '@kubb/plugin-client/clients/axios'
import type { QueryKey, QueryObserverOptions, UseQueryResult } from '@tanstack/react-query'
import { queryOptions, useQuery } from '@tanstack/react-query'

export const tasksQueryKey = () => [{ url: '/tasks' }] as const

export type TasksQueryKey = ReturnType<typeof tasksQueryKey>

/**
 * @summary Get all tasks
 * {@link /tasks}
 */
export async function tasks(config: Partial<RequestConfig> & { client?: typeof client } = {}) {
  const { client: request = client, ...requestConfig } = config

  const res = await request<TasksQueryResponse, ResponseErrorConfig<Error>, unknown>({
    method: 'GET',
    url: `/tasks`,
    baseURL: 'http://127.0.0.1:8888/v1/api',
    ...requestConfig,
  })
  return res.data
}

export function tasksQueryOptions(config: Partial<RequestConfig> & { client?: typeof client } = {}) {
  const queryKey = tasksQueryKey()
  return queryOptions<TasksQueryResponse, ResponseErrorConfig<Error>, TasksQueryResponse, typeof queryKey>({
    queryKey,
    queryFn: async ({ signal }) => {
      config.signal = signal
      return tasks(config)
    },
  })
}

/**
 * @summary Get all tasks
 * {@link /tasks}
 */
export function useTasks<TData = TasksQueryResponse, TQueryData = TasksQueryResponse, TQueryKey extends QueryKey = TasksQueryKey>(
  options: {
    query?: Partial<QueryObserverOptions<TasksQueryResponse, ResponseErrorConfig<Error>, TData, TQueryData, TQueryKey>>
    client?: Partial<RequestConfig> & { client?: typeof client }
  } = {},
) {
  const { query: queryOptions, client: config = {} } = options ?? {}
  const queryKey = queryOptions?.queryKey ?? tasksQueryKey()

  const query = useQuery({
    ...(tasksQueryOptions(config) as unknown as QueryObserverOptions),
    queryKey,
    ...(queryOptions as unknown as Omit<QueryObserverOptions, 'queryKey'>),
  }) as UseQueryResult<TData, ResponseErrorConfig<Error>> & { queryKey: TQueryKey }

  query.queryKey = queryKey as TQueryKey

  return query
}
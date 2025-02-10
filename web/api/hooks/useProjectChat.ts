import client from '@kubb/plugin-client/clients/axios'
import type { ProjectChatMutationRequest, ProjectChatMutationResponse } from '../types/ProjectChat.ts'
import type { RequestConfig, ResponseErrorConfig } from '@kubb/plugin-client/clients/axios'
import type { UseMutationOptions } from '@tanstack/react-query'
import { useMutation } from '@tanstack/react-query'

export const projectChatMutationKey = () => [{ url: '/chat' }] as const

export type ProjectChatMutationKey = ReturnType<typeof projectChatMutationKey>

/**
 * @summary Chat about your project
 * {@link /chat}
 */
export async function projectChat(
  data: ProjectChatMutationRequest,
  config: Partial<RequestConfig<ProjectChatMutationRequest>> & { client?: typeof client } = {},
) {
  const { client: request = client, ...requestConfig } = config

  const res = await request<ProjectChatMutationResponse, ResponseErrorConfig<Error>, ProjectChatMutationRequest>({
    method: 'POST',
    url: `/chat`,
    baseURL: 'http://127.0.0.1:8888/v1/api',
    data,
    ...requestConfig,
  })
  return res.data
}

/**
 * @summary Chat about your project
 * {@link /chat}
 */
export function useProjectChat(
  options: {
    mutation?: UseMutationOptions<ProjectChatMutationResponse, ResponseErrorConfig<Error>, { data: ProjectChatMutationRequest }>
    client?: Partial<RequestConfig<ProjectChatMutationRequest>> & { client?: typeof client }
  } = {},
) {
  const { mutation: mutationOptions, client: config = {} } = options ?? {}
  const mutationKey = mutationOptions?.mutationKey ?? projectChatMutationKey()

  return useMutation<ProjectChatMutationResponse, ResponseErrorConfig<Error>, { data: ProjectChatMutationRequest }>({
    mutationFn: async ({ data }) => {
      return projectChat(data, config)
    },
    mutationKey,
    ...mutationOptions,
  })
}
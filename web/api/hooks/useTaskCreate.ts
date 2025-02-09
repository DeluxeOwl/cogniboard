import client from "@kubb/plugin-client/clients/axios";
import type { TaskCreateMutationRequest, TaskCreateMutationResponse } from "../types/TaskCreate.ts";
import type { RequestConfig, ResponseErrorConfig } from "@kubb/plugin-client/clients/axios";
import type { UseMutationOptions } from "@tanstack/react-query";
import { useMutation } from "@tanstack/react-query";

export const taskCreateMutationKey = () => [{ url: "/tasks/create" }] as const;

export type TaskCreateMutationKey = ReturnType<typeof taskCreateMutationKey>;

/**
 * @summary Create a task
 * {@link /tasks/create}
 */
export async function taskCreate(
	data: TaskCreateMutationRequest,
	config: Partial<RequestConfig<TaskCreateMutationRequest>> & { client?: typeof client } = {}
) {
	const { client: request = client, ...requestConfig } = config;

	const formData = new FormData();
	if (data) {
		Object.keys(data).forEach((key) => {
			const value = data[key as keyof typeof data];
			if (typeof key === "string" && (typeof value === "string" || value instanceof Blob)) {
				formData.append(key, value);
			}
		});
	}
	const res = await request<
		TaskCreateMutationResponse,
		ResponseErrorConfig<Error>,
		TaskCreateMutationRequest
	>({
		method: "POST",
		url: `/tasks/create`,
		baseURL: "http://127.0.0.1:8888/v1/api",
		data: formData,
		headers: { "Content-Type": "multipart/form-data", ...requestConfig.headers },
		...requestConfig,
	});
	return res.data;
}

/**
 * @summary Create a task
 * {@link /tasks/create}
 */
export function useTaskCreate(
	options: {
		mutation?: UseMutationOptions<
			TaskCreateMutationResponse,
			ResponseErrorConfig<Error>,
			{ data: TaskCreateMutationRequest }
		>;
		client?: Partial<RequestConfig<TaskCreateMutationRequest>> & { client?: typeof client };
	} = {}
) {
	const { mutation: mutationOptions, client: config = {} } = options ?? {};
	const mutationKey = mutationOptions?.mutationKey ?? taskCreateMutationKey();

	return useMutation<
		TaskCreateMutationResponse,
		ResponseErrorConfig<Error>,
		{ data: TaskCreateMutationRequest }
	>({
		mutationFn: async ({ data }) => {
			return taskCreate(data, config);
		},
		mutationKey,
		...mutationOptions,
	});
}

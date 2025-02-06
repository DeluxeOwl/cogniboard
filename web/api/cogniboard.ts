/**
 * Generated by orval v7.5.0 🍺
 * Do not edit manually.
 * CogniBoard
 * OpenAPI spec version: 0.0.1
 */
import { useMutation, useQuery } from "@tanstack/react-query";
import type {
	DataTag,
	DefinedInitialDataOptions,
	DefinedUseQueryResult,
	MutationFunction,
	QueryFunction,
	QueryKey,
	UndefinedInitialDataOptions,
	UseMutationOptions,
	UseMutationResult,
	UseQueryOptions,
	UseQueryResult,
} from "@tanstack/react-query";
import type {
	AssignTaskDTO,
	ChangeTaskStatusDTO,
	CreateTaskDTO,
	ErrorModel,
	Tasks,
} from "./cogniboard.schemas";

// https://stackoverflow.com/questions/49579094/typescript-conditional-types-filter-out-readonly-properties-pick-only-requir/49579497#49579497
type IfEquals<X, Y, A = X, B = never> = (<T>() => T extends X ? 1 : 2) extends <
	T,
>() => T extends Y ? 1 : 2
	? A
	: B;

type WritableKeys<T> = {
	[P in keyof T]-?: IfEquals<
		{ [Q in P]: T[P] },
		{ -readonly [Q in P]: T[P] },
		P
	>;
}[keyof T];

type UnionToIntersection<U> = (U extends any ? (k: U) => void : never) extends (
	k: infer I,
) => void
	? I
	: never;
type DistributeReadOnlyOverUnions<T> = T extends any ? NonReadonly<T> : never;

type Writable<T> = Pick<T, WritableKeys<T>>;
type NonReadonly<T> = [T] extends [UnionToIntersection<T>]
	? {
			[P in keyof Writable<T>]: T[P] extends object
				? NonReadonly<NonNullable<T[P]>>
				: T[P];
		}
	: DistributeReadOnlyOverUnions<T>;

/**
 * @summary Get all tasks
 */
export type getTasksResponse = {
	data: Tasks | ErrorModel;
	status: number;
	headers: Headers;
};

export const getGetTasksUrl = () => {
	return `http://127.0.0.1:8888/v1/api/tasks`;
};

export const getTasks = async (
	options?: RequestInit,
): Promise<getTasksResponse> => {
	const res = await fetch(getGetTasksUrl(), {
		...options,
		method: "GET",
	});

	const body = [204, 205, 304].includes(res.status) ? null : await res.text();
	const data: getTasksResponse["data"] = body ? JSON.parse(body) : {};

	return { data, status: res.status, headers: res.headers } as getTasksResponse;
};

export const getGetTasksQueryKey = () => {
	return [`http://127.0.0.1:8888/v1/api/tasks`] as const;
};

export const getGetTasksQueryOptions = <
	TData = Awaited<ReturnType<typeof getTasks>>,
	TError = ErrorModel,
>(options?: {
	query?: Partial<
		UseQueryOptions<Awaited<ReturnType<typeof getTasks>>, TError, TData>
	>;
	fetch?: RequestInit;
}) => {
	const { query: queryOptions, fetch: fetchOptions } = options ?? {};

	const queryKey = queryOptions?.queryKey ?? getGetTasksQueryKey();

	const queryFn: QueryFunction<Awaited<ReturnType<typeof getTasks>>> = ({
		signal,
	}) => getTasks({ signal, ...fetchOptions });

	return { queryKey, queryFn, ...queryOptions } as UseQueryOptions<
		Awaited<ReturnType<typeof getTasks>>,
		TError,
		TData
	> & { queryKey: DataTag<QueryKey, TData, TError> };
};

export type GetTasksQueryResult = NonNullable<
	Awaited<ReturnType<typeof getTasks>>
>;
export type GetTasksQueryError = ErrorModel;

export function useGetTasks<
	TData = Awaited<ReturnType<typeof getTasks>>,
	TError = ErrorModel,
>(options: {
	query: Partial<
		UseQueryOptions<Awaited<ReturnType<typeof getTasks>>, TError, TData>
	> &
		Pick<
			DefinedInitialDataOptions<
				Awaited<ReturnType<typeof getTasks>>,
				TError,
				Awaited<ReturnType<typeof getTasks>>
			>,
			"initialData"
		>;
	fetch?: RequestInit;
}): DefinedUseQueryResult<TData, TError> & {
	queryKey: DataTag<QueryKey, TData, TError>;
};
export function useGetTasks<
	TData = Awaited<ReturnType<typeof getTasks>>,
	TError = ErrorModel,
>(options?: {
	query?: Partial<
		UseQueryOptions<Awaited<ReturnType<typeof getTasks>>, TError, TData>
	> &
		Pick<
			UndefinedInitialDataOptions<
				Awaited<ReturnType<typeof getTasks>>,
				TError,
				Awaited<ReturnType<typeof getTasks>>
			>,
			"initialData"
		>;
	fetch?: RequestInit;
}): UseQueryResult<TData, TError> & {
	queryKey: DataTag<QueryKey, TData, TError>;
};
export function useGetTasks<
	TData = Awaited<ReturnType<typeof getTasks>>,
	TError = ErrorModel,
>(options?: {
	query?: Partial<
		UseQueryOptions<Awaited<ReturnType<typeof getTasks>>, TError, TData>
	>;
	fetch?: RequestInit;
}): UseQueryResult<TData, TError> & {
	queryKey: DataTag<QueryKey, TData, TError>;
};
/**
 * @summary Get all tasks
 */

export function useGetTasks<
	TData = Awaited<ReturnType<typeof getTasks>>,
	TError = ErrorModel,
>(options?: {
	query?: Partial<
		UseQueryOptions<Awaited<ReturnType<typeof getTasks>>, TError, TData>
	>;
	fetch?: RequestInit;
}): UseQueryResult<TData, TError> & {
	queryKey: DataTag<QueryKey, TData, TError>;
} {
	const queryOptions = getGetTasksQueryOptions(options);

	const query = useQuery(queryOptions) as UseQueryResult<TData, TError> & {
		queryKey: DataTag<QueryKey, TData, TError>;
	};

	query.queryKey = queryOptions.queryKey;

	return query;
}

/**
 * @summary Create a task
 */
export type createTaskResponse = {
	data: void | ErrorModel;
	status: number;
	headers: Headers;
};

export const getCreateTaskUrl = () => {
	return `http://127.0.0.1:8888/v1/api/tasks/create`;
};

export const createTask = async (
	createTaskDto: NonReadonly<CreateTaskDTO>,
	options?: RequestInit,
): Promise<createTaskResponse> => {
	const res = await fetch(getCreateTaskUrl(), {
		...options,
		method: "POST",
		headers: { "Content-Type": "application/json", ...options?.headers },
		body: JSON.stringify(createTaskDto),
	});

	const body = [204, 205, 304].includes(res.status) ? null : await res.text();
	const data: createTaskResponse["data"] = body ? JSON.parse(body) : {};

	return {
		data,
		status: res.status,
		headers: res.headers,
	} as createTaskResponse;
};

export const getCreateTaskMutationOptions = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof createTask>>,
		TError,
		{ data: NonReadonly<CreateTaskDTO> },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationOptions<
	Awaited<ReturnType<typeof createTask>>,
	TError,
	{ data: NonReadonly<CreateTaskDTO> },
	TContext
> => {
	const mutationKey = ["createTask"];
	const { mutation: mutationOptions, fetch: fetchOptions } = options
		? options.mutation &&
			"mutationKey" in options.mutation &&
			options.mutation.mutationKey
			? options
			: { ...options, mutation: { ...options.mutation, mutationKey } }
		: { mutation: { mutationKey }, fetch: undefined };

	const mutationFn: MutationFunction<
		Awaited<ReturnType<typeof createTask>>,
		{ data: NonReadonly<CreateTaskDTO> }
	> = (props) => {
		const { data } = props ?? {};

		return createTask(data, fetchOptions);
	};

	return { mutationFn, ...mutationOptions };
};

export type CreateTaskMutationResult = NonNullable<
	Awaited<ReturnType<typeof createTask>>
>;
export type CreateTaskMutationBody = NonReadonly<CreateTaskDTO>;
export type CreateTaskMutationError = ErrorModel;

/**
 * @summary Create a task
 */
export const useCreateTask = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof createTask>>,
		TError,
		{ data: NonReadonly<CreateTaskDTO> },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationResult<
	Awaited<ReturnType<typeof createTask>>,
	TError,
	{ data: NonReadonly<CreateTaskDTO> },
	TContext
> => {
	const mutationOptions = getCreateTaskMutationOptions(options);

	return useMutation(mutationOptions);
};

/**
 * @summary Assign a task to someone
 */
export type assignTaskResponse = {
	data: void | ErrorModel;
	status: number;
	headers: Headers;
};

export const getAssignTaskUrl = (taskId: string) => {
	return `http://127.0.0.1:8888/v1/api/tasks/${taskId}/assign`;
};

export const assignTask = async (
	taskId: string,
	assignTaskDto: NonReadonly<AssignTaskDTO>,
	options?: RequestInit,
): Promise<assignTaskResponse> => {
	const res = await fetch(getAssignTaskUrl(taskId), {
		...options,
		method: "POST",
		headers: { "Content-Type": "application/json", ...options?.headers },
		body: JSON.stringify(assignTaskDto),
	});

	const body = [204, 205, 304].includes(res.status) ? null : await res.text();
	const data: assignTaskResponse["data"] = body ? JSON.parse(body) : {};

	return {
		data,
		status: res.status,
		headers: res.headers,
	} as assignTaskResponse;
};

export const getAssignTaskMutationOptions = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof assignTask>>,
		TError,
		{ taskId: string; data: NonReadonly<AssignTaskDTO> },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationOptions<
	Awaited<ReturnType<typeof assignTask>>,
	TError,
	{ taskId: string; data: NonReadonly<AssignTaskDTO> },
	TContext
> => {
	const mutationKey = ["assignTask"];
	const { mutation: mutationOptions, fetch: fetchOptions } = options
		? options.mutation &&
			"mutationKey" in options.mutation &&
			options.mutation.mutationKey
			? options
			: { ...options, mutation: { ...options.mutation, mutationKey } }
		: { mutation: { mutationKey }, fetch: undefined };

	const mutationFn: MutationFunction<
		Awaited<ReturnType<typeof assignTask>>,
		{ taskId: string; data: NonReadonly<AssignTaskDTO> }
	> = (props) => {
		const { taskId, data } = props ?? {};

		return assignTask(taskId, data, fetchOptions);
	};

	return { mutationFn, ...mutationOptions };
};

export type AssignTaskMutationResult = NonNullable<
	Awaited<ReturnType<typeof assignTask>>
>;
export type AssignTaskMutationBody = NonReadonly<AssignTaskDTO>;
export type AssignTaskMutationError = ErrorModel;

/**
 * @summary Assign a task to someone
 */
export const useAssignTask = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof assignTask>>,
		TError,
		{ taskId: string; data: NonReadonly<AssignTaskDTO> },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationResult<
	Awaited<ReturnType<typeof assignTask>>,
	TError,
	{ taskId: string; data: NonReadonly<AssignTaskDTO> },
	TContext
> => {
	const mutationOptions = getAssignTaskMutationOptions(options);

	return useMutation(mutationOptions);
};

/**
 * @summary Change task status
 */
export type changeTaskStatusResponse = {
	data: void | ErrorModel;
	status: number;
	headers: Headers;
};

export const getChangeTaskStatusUrl = (taskId: string) => {
	return `http://127.0.0.1:8888/v1/api/tasks/${taskId}/status`;
};

export const changeTaskStatus = async (
	taskId: string,
	changeTaskStatusDto: NonReadonly<ChangeTaskStatusDTO>,
	options?: RequestInit,
): Promise<changeTaskStatusResponse> => {
	const res = await fetch(getChangeTaskStatusUrl(taskId), {
		...options,
		method: "POST",
		headers: { "Content-Type": "application/json", ...options?.headers },
		body: JSON.stringify(changeTaskStatusDto),
	});

	const body = [204, 205, 304].includes(res.status) ? null : await res.text();
	const data: changeTaskStatusResponse["data"] = body ? JSON.parse(body) : {};

	return {
		data,
		status: res.status,
		headers: res.headers,
	} as changeTaskStatusResponse;
};

export const getChangeTaskStatusMutationOptions = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof changeTaskStatus>>,
		TError,
		{ taskId: string; data: NonReadonly<ChangeTaskStatusDTO> },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationOptions<
	Awaited<ReturnType<typeof changeTaskStatus>>,
	TError,
	{ taskId: string; data: NonReadonly<ChangeTaskStatusDTO> },
	TContext
> => {
	const mutationKey = ["changeTaskStatus"];
	const { mutation: mutationOptions, fetch: fetchOptions } = options
		? options.mutation &&
			"mutationKey" in options.mutation &&
			options.mutation.mutationKey
			? options
			: { ...options, mutation: { ...options.mutation, mutationKey } }
		: { mutation: { mutationKey }, fetch: undefined };

	const mutationFn: MutationFunction<
		Awaited<ReturnType<typeof changeTaskStatus>>,
		{ taskId: string; data: NonReadonly<ChangeTaskStatusDTO> }
	> = (props) => {
		const { taskId, data } = props ?? {};

		return changeTaskStatus(taskId, data, fetchOptions);
	};

	return { mutationFn, ...mutationOptions };
};

export type ChangeTaskStatusMutationResult = NonNullable<
	Awaited<ReturnType<typeof changeTaskStatus>>
>;
export type ChangeTaskStatusMutationBody = NonReadonly<ChangeTaskStatusDTO>;
export type ChangeTaskStatusMutationError = ErrorModel;

/**
 * @summary Change task status
 */
export const useChangeTaskStatus = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof changeTaskStatus>>,
		TError,
		{ taskId: string; data: NonReadonly<ChangeTaskStatusDTO> },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationResult<
	Awaited<ReturnType<typeof changeTaskStatus>>,
	TError,
	{ taskId: string; data: NonReadonly<ChangeTaskStatusDTO> },
	TContext
> => {
	const mutationOptions = getChangeTaskStatusMutationOptions(options);

	return useMutation(mutationOptions);
};

/**
 * @summary Unassign a task
 */
export type unassignTaskResponse = {
	data: void | ErrorModel;
	status: number;
	headers: Headers;
};

export const getUnassignTaskUrl = (taskId: string) => {
	return `http://127.0.0.1:8888/v1/api/tasks/${taskId}/unassign`;
};

export const unassignTask = async (
	taskId: string,
	options?: RequestInit,
): Promise<unassignTaskResponse> => {
	const res = await fetch(getUnassignTaskUrl(taskId), {
		...options,
		method: "POST",
	});

	const body = [204, 205, 304].includes(res.status) ? null : await res.text();
	const data: unassignTaskResponse["data"] = body ? JSON.parse(body) : {};

	return {
		data,
		status: res.status,
		headers: res.headers,
	} as unassignTaskResponse;
};

export const getUnassignTaskMutationOptions = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof unassignTask>>,
		TError,
		{ taskId: string },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationOptions<
	Awaited<ReturnType<typeof unassignTask>>,
	TError,
	{ taskId: string },
	TContext
> => {
	const mutationKey = ["unassignTask"];
	const { mutation: mutationOptions, fetch: fetchOptions } = options
		? options.mutation &&
			"mutationKey" in options.mutation &&
			options.mutation.mutationKey
			? options
			: { ...options, mutation: { ...options.mutation, mutationKey } }
		: { mutation: { mutationKey }, fetch: undefined };

	const mutationFn: MutationFunction<
		Awaited<ReturnType<typeof unassignTask>>,
		{ taskId: string }
	> = (props) => {
		const { taskId } = props ?? {};

		return unassignTask(taskId, fetchOptions);
	};

	return { mutationFn, ...mutationOptions };
};

export type UnassignTaskMutationResult = NonNullable<
	Awaited<ReturnType<typeof unassignTask>>
>;

export type UnassignTaskMutationError = ErrorModel;

/**
 * @summary Unassign a task
 */
export const useUnassignTask = <
	TError = ErrorModel,
	TContext = unknown,
>(options?: {
	mutation?: UseMutationOptions<
		Awaited<ReturnType<typeof unassignTask>>,
		TError,
		{ taskId: string },
		TContext
	>;
	fetch?: RequestInit;
}): UseMutationResult<
	Awaited<ReturnType<typeof unassignTask>>,
	TError,
	{ taskId: string },
	TContext
> => {
	const mutationOptions = getUnassignTaskMutationOptions(options);

	return useMutation(mutationOptions);
};

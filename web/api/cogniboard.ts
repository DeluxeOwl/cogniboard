/**
 * Generated by orval v7.5.0 🍺
 * Do not edit manually.
 * CogniBoard
 * OpenAPI spec version: 0.0.1
 */
import {
  useMutation,
  useQuery
} from '@tanstack/react-query'
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
  UseQueryResult
} from '@tanstack/react-query'
import axios from 'axios'
import type {
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse
} from 'axios'

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

type UnionToIntersection<U> =
  (U extends any ? (k: U)=>void : never) extends ((k: infer I)=>void) ? I : never;
type DistributeReadOnlyOverUnions<T> = T extends any ? NonReadonly<T> : never;

type Writable<T> = Pick<T, WritableKeys<T>>;
type NonReadonly<T> = [T] extends [UnionToIntersection<T>] ? {
  [P in keyof Writable<T>]: T[P] extends object
    ? NonReadonly<NonNullable<T[P]>>
    : T[P];
} : DistributeReadOnlyOverUnions<T>;

export interface AssignTaskDTO {
  /** A URL to the JSON Schema for this object. */
  readonly $schema?: string;
  /**
   * Name of the person to assign the task to
   * @minLength 1
   */
  assignee_name: string;
}

export interface ChangeTaskStatusDTO {
  /** A URL to the JSON Schema for this object. */
  readonly $schema?: string;
  /**
   * New status for the task
   * @minLength 1
   */
  status: string;
}

export interface CreateTaskDTO {
  /** A URL to the JSON Schema for this object. */
  readonly $schema?: string;
  /** Task's asignee (if any) */
  assignee_name?: string;
  /** Task's description */
  description?: string;
  /** Task's due date (if any) */
  due_date?: string;
  /**
   * Task's name
   * @minLength 1
   * @maxLength 50
   */
  title: string;
}

export interface ErrorDetail {
  /** Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id' */
  location?: string;
  /** Error message text */
  message?: string;
  /** The value at the given location */
  value?: unknown;
}

export interface ErrorModel {
  /** A URL to the JSON Schema for this object. */
  readonly $schema?: string;
  /** A human-readable explanation specific to this occurrence of the problem. */
  detail?: string;
  /**
   * Optional list of individual error details
   * @nullable
   */
  errors?: ErrorDetail[] | null;
  /** A URI reference that identifies the specific occurrence of the problem. */
  instance?: string;
  /** HTTP status code */
  status?: number;
  /** A short, human-readable summary of the problem type. This value should not change between occurrences of the error. */
  title?: string;
  /** A URI reference to human-readable documentation for the error. */
  type?: string;
}

export interface GetTasksDTO {
  /** @nullable */
  assignee: string | null;
  /** @nullable */
  completed_at: string | null;
  created_at: string;
  /** @nullable */
  description: string | null;
  /** @nullable */
  due_date: string | null;
  id: string;
  status: string;
  title: string;
}

export interface Tasks {
  /** A URL to the JSON Schema for this object. */
  readonly $schema?: string;
  /** @nullable */
  tasks: GetTasksDTO[] | null;
}





/**
 * @summary Get all tasks
 */
export const tasks = (
     options?: AxiosRequestConfig
 ): Promise<AxiosResponse<Tasks>> => {
    
    
    return axios.get(
      `http://127.0.0.1:8888/v1/api/tasks`,options
    );
  }


export const getTasksQueryKey = () => {
    return [`http://127.0.0.1:8888/v1/api/tasks`] as const;
    }

    
export const getTasksQueryOptions = <TData = Awaited<ReturnType<typeof tasks>>, TError = AxiosError<ErrorModel>>( options?: { query?:Partial<UseQueryOptions<Awaited<ReturnType<typeof tasks>>, TError, TData>>, axios?: AxiosRequestConfig}
) => {

const {query: queryOptions, axios: axiosOptions} = options ?? {};

  const queryKey =  queryOptions?.queryKey ?? getTasksQueryKey();

  

    const queryFn: QueryFunction<Awaited<ReturnType<typeof tasks>>> = ({ signal }) => tasks({ signal, ...axiosOptions });

      

      

   return  { queryKey, queryFn, ...queryOptions} as UseQueryOptions<Awaited<ReturnType<typeof tasks>>, TError, TData> & { queryKey: DataTag<QueryKey, TData, TError> }
}

export type TasksQueryResult = NonNullable<Awaited<ReturnType<typeof tasks>>>
export type TasksQueryError = AxiosError<ErrorModel>


export function useTasks<TData = Awaited<ReturnType<typeof tasks>>, TError = AxiosError<ErrorModel>>(
  options: { query:Partial<UseQueryOptions<Awaited<ReturnType<typeof tasks>>, TError, TData>> & Pick<
        DefinedInitialDataOptions<
          Awaited<ReturnType<typeof tasks>>,
          TError,
          Awaited<ReturnType<typeof tasks>>
        > , 'initialData'
      >, axios?: AxiosRequestConfig}

  ):  DefinedUseQueryResult<TData, TError> & { queryKey: DataTag<QueryKey, TData, TError> }
export function useTasks<TData = Awaited<ReturnType<typeof tasks>>, TError = AxiosError<ErrorModel>>(
  options?: { query?:Partial<UseQueryOptions<Awaited<ReturnType<typeof tasks>>, TError, TData>> & Pick<
        UndefinedInitialDataOptions<
          Awaited<ReturnType<typeof tasks>>,
          TError,
          Awaited<ReturnType<typeof tasks>>
        > , 'initialData'
      >, axios?: AxiosRequestConfig}

  ):  UseQueryResult<TData, TError> & { queryKey: DataTag<QueryKey, TData, TError> }
export function useTasks<TData = Awaited<ReturnType<typeof tasks>>, TError = AxiosError<ErrorModel>>(
  options?: { query?:Partial<UseQueryOptions<Awaited<ReturnType<typeof tasks>>, TError, TData>>, axios?: AxiosRequestConfig}

  ):  UseQueryResult<TData, TError> & { queryKey: DataTag<QueryKey, TData, TError> }
/**
 * @summary Get all tasks
 */

export function useTasks<TData = Awaited<ReturnType<typeof tasks>>, TError = AxiosError<ErrorModel>>(
  options?: { query?:Partial<UseQueryOptions<Awaited<ReturnType<typeof tasks>>, TError, TData>>, axios?: AxiosRequestConfig}

  ):  UseQueryResult<TData, TError> & { queryKey: DataTag<QueryKey, TData, TError> } {

  const queryOptions = getTasksQueryOptions(options)

  const query = useQuery(queryOptions) as  UseQueryResult<TData, TError> & { queryKey: DataTag<QueryKey, TData, TError> };

  query.queryKey = queryOptions.queryKey ;

  return query;
}




/**
 * @summary Create a task
 */
export const createTask = (
    createTaskDTO: NonReadonly<CreateTaskDTO>, options?: AxiosRequestConfig
 ): Promise<AxiosResponse<void>> => {
    
    
    return axios.post(
      `http://127.0.0.1:8888/v1/api/tasks/create`,
      createTaskDTO,options
    );
  }



export const getCreateTaskMutationOptions = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof createTask>>, TError,{data: NonReadonly<CreateTaskDTO>}, TContext>, axios?: AxiosRequestConfig}
): UseMutationOptions<Awaited<ReturnType<typeof createTask>>, TError,{data: NonReadonly<CreateTaskDTO>}, TContext> => {
    
const mutationKey = ['createTask'];
const {mutation: mutationOptions, axios: axiosOptions} = options ?
      options.mutation && 'mutationKey' in options.mutation && options.mutation.mutationKey ?
      options
      : {...options, mutation: {...options.mutation, mutationKey}}
      : {mutation: { mutationKey, }, axios: undefined};

      


      const mutationFn: MutationFunction<Awaited<ReturnType<typeof createTask>>, {data: NonReadonly<CreateTaskDTO>}> = (props) => {
          const {data} = props ?? {};

          return  createTask(data,axiosOptions)
        }

        


  return  { mutationFn, ...mutationOptions }}

    export type CreateTaskMutationResult = NonNullable<Awaited<ReturnType<typeof createTask>>>
    export type CreateTaskMutationBody = NonReadonly<CreateTaskDTO>
    export type CreateTaskMutationError = AxiosError<ErrorModel>

    /**
 * @summary Create a task
 */
export const useCreateTask = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof createTask>>, TError,{data: NonReadonly<CreateTaskDTO>}, TContext>, axios?: AxiosRequestConfig}
): UseMutationResult<
        Awaited<ReturnType<typeof createTask>>,
        TError,
        {data: NonReadonly<CreateTaskDTO>},
        TContext
      > => {

      const mutationOptions = getCreateTaskMutationOptions(options);

      return useMutation(mutationOptions);
    }
    
/**
 * @summary Assign a task to someone
 */
export const taskAssign = (
    taskId: string,
    assignTaskDTO: NonReadonly<AssignTaskDTO>, options?: AxiosRequestConfig
 ): Promise<AxiosResponse<void>> => {
    
    
    return axios.post(
      `http://127.0.0.1:8888/v1/api/tasks/${taskId}/assign`,
      assignTaskDTO,options
    );
  }



export const getTaskAssignMutationOptions = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof taskAssign>>, TError,{taskId: string;data: NonReadonly<AssignTaskDTO>}, TContext>, axios?: AxiosRequestConfig}
): UseMutationOptions<Awaited<ReturnType<typeof taskAssign>>, TError,{taskId: string;data: NonReadonly<AssignTaskDTO>}, TContext> => {
    
const mutationKey = ['taskAssign'];
const {mutation: mutationOptions, axios: axiosOptions} = options ?
      options.mutation && 'mutationKey' in options.mutation && options.mutation.mutationKey ?
      options
      : {...options, mutation: {...options.mutation, mutationKey}}
      : {mutation: { mutationKey, }, axios: undefined};

      


      const mutationFn: MutationFunction<Awaited<ReturnType<typeof taskAssign>>, {taskId: string;data: NonReadonly<AssignTaskDTO>}> = (props) => {
          const {taskId,data} = props ?? {};

          return  taskAssign(taskId,data,axiosOptions)
        }

        


  return  { mutationFn, ...mutationOptions }}

    export type TaskAssignMutationResult = NonNullable<Awaited<ReturnType<typeof taskAssign>>>
    export type TaskAssignMutationBody = NonReadonly<AssignTaskDTO>
    export type TaskAssignMutationError = AxiosError<ErrorModel>

    /**
 * @summary Assign a task to someone
 */
export const useTaskAssign = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof taskAssign>>, TError,{taskId: string;data: NonReadonly<AssignTaskDTO>}, TContext>, axios?: AxiosRequestConfig}
): UseMutationResult<
        Awaited<ReturnType<typeof taskAssign>>,
        TError,
        {taskId: string;data: NonReadonly<AssignTaskDTO>},
        TContext
      > => {

      const mutationOptions = getTaskAssignMutationOptions(options);

      return useMutation(mutationOptions);
    }
    
/**
 * @summary Change task status
 */
export const taskChangeStatus = (
    taskId: string,
    changeTaskStatusDTO: NonReadonly<ChangeTaskStatusDTO>, options?: AxiosRequestConfig
 ): Promise<AxiosResponse<void>> => {
    
    
    return axios.post(
      `http://127.0.0.1:8888/v1/api/tasks/${taskId}/status`,
      changeTaskStatusDTO,options
    );
  }



export const getTaskChangeStatusMutationOptions = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof taskChangeStatus>>, TError,{taskId: string;data: NonReadonly<ChangeTaskStatusDTO>}, TContext>, axios?: AxiosRequestConfig}
): UseMutationOptions<Awaited<ReturnType<typeof taskChangeStatus>>, TError,{taskId: string;data: NonReadonly<ChangeTaskStatusDTO>}, TContext> => {
    
const mutationKey = ['taskChangeStatus'];
const {mutation: mutationOptions, axios: axiosOptions} = options ?
      options.mutation && 'mutationKey' in options.mutation && options.mutation.mutationKey ?
      options
      : {...options, mutation: {...options.mutation, mutationKey}}
      : {mutation: { mutationKey, }, axios: undefined};

      


      const mutationFn: MutationFunction<Awaited<ReturnType<typeof taskChangeStatus>>, {taskId: string;data: NonReadonly<ChangeTaskStatusDTO>}> = (props) => {
          const {taskId,data} = props ?? {};

          return  taskChangeStatus(taskId,data,axiosOptions)
        }

        


  return  { mutationFn, ...mutationOptions }}

    export type TaskChangeStatusMutationResult = NonNullable<Awaited<ReturnType<typeof taskChangeStatus>>>
    export type TaskChangeStatusMutationBody = NonReadonly<ChangeTaskStatusDTO>
    export type TaskChangeStatusMutationError = AxiosError<ErrorModel>

    /**
 * @summary Change task status
 */
export const useTaskChangeStatus = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof taskChangeStatus>>, TError,{taskId: string;data: NonReadonly<ChangeTaskStatusDTO>}, TContext>, axios?: AxiosRequestConfig}
): UseMutationResult<
        Awaited<ReturnType<typeof taskChangeStatus>>,
        TError,
        {taskId: string;data: NonReadonly<ChangeTaskStatusDTO>},
        TContext
      > => {

      const mutationOptions = getTaskChangeStatusMutationOptions(options);

      return useMutation(mutationOptions);
    }
    
/**
 * @summary Unassign a task
 */
export const taskUnassign = (
    taskId: string, options?: AxiosRequestConfig
 ): Promise<AxiosResponse<void>> => {
    
    
    return axios.post(
      `http://127.0.0.1:8888/v1/api/tasks/${taskId}/unassign`,undefined,options
    );
  }



export const getTaskUnassignMutationOptions = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof taskUnassign>>, TError,{taskId: string}, TContext>, axios?: AxiosRequestConfig}
): UseMutationOptions<Awaited<ReturnType<typeof taskUnassign>>, TError,{taskId: string}, TContext> => {
    
const mutationKey = ['taskUnassign'];
const {mutation: mutationOptions, axios: axiosOptions} = options ?
      options.mutation && 'mutationKey' in options.mutation && options.mutation.mutationKey ?
      options
      : {...options, mutation: {...options.mutation, mutationKey}}
      : {mutation: { mutationKey, }, axios: undefined};

      


      const mutationFn: MutationFunction<Awaited<ReturnType<typeof taskUnassign>>, {taskId: string}> = (props) => {
          const {taskId} = props ?? {};

          return  taskUnassign(taskId,axiosOptions)
        }

        


  return  { mutationFn, ...mutationOptions }}

    export type TaskUnassignMutationResult = NonNullable<Awaited<ReturnType<typeof taskUnassign>>>
    
    export type TaskUnassignMutationError = AxiosError<ErrorModel>

    /**
 * @summary Unassign a task
 */
export const useTaskUnassign = <TError = AxiosError<ErrorModel>,
    TContext = unknown>(options?: { mutation?:UseMutationOptions<Awaited<ReturnType<typeof taskUnassign>>, TError,{taskId: string}, TContext>, axios?: AxiosRequestConfig}
): UseMutationResult<
        Awaited<ReturnType<typeof taskUnassign>>,
        TError,
        {taskId: string},
        TContext
      > => {

      const mutationOptions = getTaskUnassignMutationOptions(options);

      return useMutation(mutationOptions);
    }
    

export type { TaskChangeStatusMutationKey } from './hooks/useTaskChangeStatus.ts'
export type { TaskCreateMutationKey } from './hooks/useTaskCreate.ts'
export type { TaskEditMutationKey } from './hooks/useTaskEdit.ts'
export type { TasksQueryKey } from './hooks/useTasks.ts'
export type { TasksSuspenseQueryKey } from './hooks/useTasksSuspense.ts'
export type { ErrorDetail } from './types/ErrorDetail.ts'
export type { ErrorModel } from './types/ErrorModel.ts'
export type { InChangeTaskStatusDTO } from './types/InChangeTaskStatusDTO.ts'
export type { InTaskDTO } from './types/InTaskDTO.ts'
export type { InTasksDTO } from './types/InTasksDTO.ts'
export type {
  TaskChangeStatusPathParams,
  TaskChangeStatus204,
  TaskChangeStatusError,
  TaskChangeStatusMutationRequest,
  TaskChangeStatusMutationResponse,
  TaskChangeStatusMutation,
} from './types/TaskChangeStatus.ts'
export type { TaskCreate204, TaskCreateError, TaskCreateMutationRequest, TaskCreateMutationResponse, TaskCreateMutation } from './types/TaskCreate.ts'
export type { TaskEditPathParams, TaskEdit204, TaskEditError, TaskEditMutationRequest, TaskEditMutationResponse, TaskEditMutation } from './types/TaskEdit.ts'
export type { Tasks200, TasksError, TasksQueryResponse, TasksQuery } from './types/Tasks.ts'
export { taskChangeStatusMutationKey, taskChangeStatus, useTaskChangeStatus } from './hooks/useTaskChangeStatus.ts'
export { taskCreateMutationKey, taskCreate, useTaskCreate } from './hooks/useTaskCreate.ts'
export { taskEditMutationKey, taskEdit, useTaskEdit } from './hooks/useTaskEdit.ts'
export { tasksQueryKey, tasks, tasksQueryOptions, useTasks } from './hooks/useTasks.ts'
export { tasksSuspenseQueryKey, tasksSuspense, tasksSuspenseQueryOptions, useTasksSuspense } from './hooks/useTasksSuspense.ts'
export { errorDetailSchema } from './zod/errorDetailSchema.ts'
export { errorModelSchema } from './zod/errorModelSchema.ts'
export { inChangeTaskStatusDTOSchema } from './zod/inChangeTaskStatusDTOSchema.ts'
export { inTaskDTOSchema } from './zod/inTaskDTOSchema.ts'
export { inTasksDTOSchema } from './zod/inTasksDTOSchema.ts'
export {
  taskChangeStatusPathParamsSchema,
  taskChangeStatus204Schema,
  taskChangeStatusErrorSchema,
  taskChangeStatusMutationRequestSchema,
  taskChangeStatusMutationResponseSchema,
} from './zod/taskChangeStatusSchema.ts'
export { taskCreate204Schema, taskCreateErrorSchema, taskCreateMutationRequestSchema, taskCreateMutationResponseSchema } from './zod/taskCreateSchema.ts'
export {
  taskEditPathParamsSchema,
  taskEdit204Schema,
  taskEditErrorSchema,
  taskEditMutationRequestSchema,
  taskEditMutationResponseSchema,
} from './zod/taskEditSchema.ts'
export { tasks200Schema, tasksErrorSchema, tasksQueryResponseSchema } from './zod/tasksSchema.ts'
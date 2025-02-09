export type { TaskChangeStatusMutationKey } from './hooks/useTaskChangeStatus.ts'
export type { TaskCreateMutationKey } from './hooks/useTaskCreate.ts'
export type { TaskEditMutationKey } from './hooks/useTaskEdit.ts'
export type { TasksQueryKey } from './hooks/useTasks.ts'
export type { TasksSuspenseQueryKey } from './hooks/useTasksSuspense.ts'
export type { ChangeTaskStatus } from './types/ChangeTaskStatus.ts'
export type { ErrorDetail } from './types/ErrorDetail.ts'
export type { ErrorModel } from './types/ErrorModel.ts'
export type { FileDTO } from './types/FileDTO.ts'
export type { ListTasks } from './types/ListTasks.ts'
export type {
  TaskChangeStatusPathParams,
  TaskChangeStatus204,
  TaskChangeStatusError,
  TaskChangeStatusMutationRequest,
  TaskChangeStatusMutationResponse,
  TaskChangeStatusMutation,
} from './types/TaskChangeStatus.ts'
export type { TaskCreate204, TaskCreateError, TaskCreateMutationRequest, TaskCreateMutationResponse, TaskCreateMutation } from './types/TaskCreate.ts'
export type { TaskDTO } from './types/TaskDTO.ts'
export type { TaskEditPathParams, TaskEdit204, TaskEditError, TaskEditMutationRequest, TaskEditMutationResponse, TaskEditMutation } from './types/TaskEdit.ts'
export type { Tasks200, TasksError, TasksQueryResponse, TasksQuery } from './types/Tasks.ts'
export { taskChangeStatusMutationKey, taskChangeStatus, useTaskChangeStatus } from './hooks/useTaskChangeStatus.ts'
export { taskCreateMutationKey, taskCreate, useTaskCreate } from './hooks/useTaskCreate.ts'
export { taskEditMutationKey, taskEdit, useTaskEdit } from './hooks/useTaskEdit.ts'
export { tasksQueryKey, tasks, tasksQueryOptions, useTasks } from './hooks/useTasks.ts'
export { tasksSuspenseQueryKey, tasksSuspense, tasksSuspenseQueryOptions, useTasksSuspense } from './hooks/useTasksSuspense.ts'
export { changeTaskStatusSchema } from './zod/changeTaskStatusSchema.ts'
export { errorDetailSchema } from './zod/errorDetailSchema.ts'
export { errorModelSchema } from './zod/errorModelSchema.ts'
export { fileDTOSchema } from './zod/fileDTOSchema.ts'
export { listTasksSchema } from './zod/listTasksSchema.ts'
export {
  taskChangeStatusPathParamsSchema,
  taskChangeStatus204Schema,
  taskChangeStatusErrorSchema,
  taskChangeStatusMutationRequestSchema,
  taskChangeStatusMutationResponseSchema,
} from './zod/taskChangeStatusSchema.ts'
export { taskCreate204Schema, taskCreateErrorSchema, taskCreateMutationRequestSchema, taskCreateMutationResponseSchema } from './zod/taskCreateSchema.ts'
export { taskDTOSchema } from './zod/taskDTOSchema.ts'
export {
  taskEditPathParamsSchema,
  taskEdit204Schema,
  taskEditErrorSchema,
  taskEditMutationRequestSchema,
  taskEditMutationResponseSchema,
} from './zod/taskEditSchema.ts'
export { tasks200Schema, tasksErrorSchema, tasksQueryResponseSchema } from './zod/tasksSchema.ts'
import { type Task, tasksQueryKey, useTaskChangeStatus, useTasks } from "@/api/index.ts";
import { useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
	KanbanBoard,
	KanbanCard,
	KanbanCards,
	KanbanHeader,
	KanbanProvider,
} from "./project/kanban";
import type { DragEndEvent } from "./project/kanban.ts";
import AddTaskDialog from "./project/task.add";
import EditTaskDialog from "./project/task.edit";
import { Button } from "@/components/ui/button.tsx";
import {
	Sheet,
	SheetContent,
	SheetDescription,
	SheetHeader,
	SheetTitle,
	SheetTrigger,
} from "@/components/ui/sheet.tsx";

function useHome() {
	const [selectedTask, setSelectedTask] = useState<Task | null>(null);
	const [editDialogOpen, setEditDialogOpen] = useState(false);
	const { tasks, isLoading, isError, error, handleDragEnd } = useKanbanBoard();

	const isLoadingOrNoTasks = isLoading || !tasks;
	const hasError = isError && error;

	const errorMessage = hasError ? error.response?.data.message : "";

	const handleTaskSelect = (task: Task) => {
		setSelectedTask(task);
		setEditDialogOpen(true);
	};

	// Ensure tasks is always an array
	const safeTasks = tasks ?? [];

	return {
		tasks: safeTasks,
		selectedTask,
		editDialogOpen,
		isLoadingOrNoTasks,
		hasError,
		errorMessage,
		handleDragEnd,
		handleTaskSelect,
		setEditDialogOpen,
	};
}

const Home = () => {
	const {
		tasks,
		selectedTask,
		editDialogOpen,
		isLoadingOrNoTasks,
		hasError,
		errorMessage,
		handleDragEnd,
		handleTaskSelect,
		setEditDialogOpen,
	} = useHome();

	if (isLoadingOrNoTasks) {
		return <div>Loading ...</div>;
	}

	if (hasError) {
		return <div>Error: {errorMessage}</div>;
	}

	return (
		<main className="p-2">
			<span className="flex justify-between mx-4">
				<AddTaskDialog />
				<Sheet>
					<SheetTrigger asChild>
						<Button
							className="bg-orange-600 hover:bg-orange-500"
							variant="default"
						>
							Start chatting
						</Button>
					</SheetTrigger>
					<SheetContent>
						<SheetHeader>
							<SheetTitle>Chat about your project</SheetTitle>
							<SheetDescription>
								Ask questions about tasks, uploaded documents, assignees etc.
							</SheetDescription>
						</SheetHeader>
					</SheetContent>
				</Sheet>
			</span>
			{selectedTask && (
				<EditTaskDialog
					task={selectedTask}
					open={editDialogOpen}
					onOpenChange={setEditDialogOpen}
				/>
			)}
			<KanbanProvider
				onDragEnd={handleDragEnd}
				className="p-4"
			>
				{statuses.map((status) => (
					<KanbanBoard
						key={status.name}
						id={status.name}
					>
						<KanbanHeader
							name={status.name}
							color={status.color}
						/>
						<KanbanCards>
							{tasks
								.filter((task) => task.status === status.name)
								.map((feature, index) => (
									<KanbanCard
										key={feature.updated_at}
										id={feature.id}
										name={feature.title}
										parent={status.name}
										index={index}
										onClick={() => handleTaskSelect(feature)}
									>
										<div className="flex items-start justify-between gap-2">
											<div className="flex flex-col gap-1 min-w-0 w-full">
												<p className="m-0 flex-1 font-bold text-sm truncate overflow-hidden">
													{feature.title}
												</p>
												<p className="m-0 text-muted-foreground text-xs">{feature.description}</p>
												{feature.assignee && <p>Assigned to: {feature.assignee}</p>}
											</div>
										</div>
										<p className="m-0 text-muted-foreground text-xs">
											{shortDateFormatter.format(new Date(feature.created_at))}
											{feature.completed_at
												? `- ${dateFormatter.format(new Date(feature.completed_at))}`
												: null}
										</p>
									</KanbanCard>
								))}
						</KanbanCards>
					</KanbanBoard>
				))}
			</KanbanProvider>
		</main>
	);
};

export default Home;

const dateFormatter = new Intl.DateTimeFormat("en-US", {
	month: "short",
	day: "numeric",
	year: "numeric",
});

const shortDateFormatter = new Intl.DateTimeFormat("en-US", {
	month: "short",
	day: "numeric",
});

const statuses = [
	{ id: "1", name: "pending", color: "#6B7280" },
	{ id: "2", name: "in_progress", color: "#F59E0B" },
	{ id: "3", name: "in_review", color: "red" },
	{ id: "4", name: "completed", color: "#10B981" },
];

export const assignees = ["John", "Mary", "Steve", "Laura", "Alex"];

function useKanbanBoard() {
	const queryClient = useQueryClient();

	const { data, isLoading, isError, error } = useTasks();

	const mutation = useTaskChangeStatus();

	const handleDragEnd = (event: DragEndEvent) => {
		const { active, over } = event;

		if (!over) {
			return;
		}

		const status = statuses.find((status) => status.name === over.id);

		if (!status) {
			return;
		}

		const taskId = String(active.id);
		const newStatus = status.name;

		// Get the current query data
		const queryKey = tasksQueryKey();
		type DataType = typeof data;
		const previousTasks = (queryClient.getQueryData(queryKey) as DataType)?.tasks;

		// Find the task's current status
		const currentTask = previousTasks?.find((task: any) => task.id === taskId);
		if (currentTask?.status === newStatus) {
			return; // Skip if status hasn't changed
		}

		// Optimistically update the task status
		queryClient.setQueryData(queryKey, (old: DataType) => {
			return {
				$schema: old?.$schema,
				tasks: old?.tasks?.map((task: Task) =>
					task.id === taskId ? { ...task, status: newStatus } : task
				),
			};
		});

		mutation.mutate(
			{ taskId, data: { status: newStatus } },
			{
				// Still invalidate on success to ensure we have the latest data
				onSuccess: () => {
					queryClient.invalidateQueries({ queryKey: queryKey });
				},
				// Revert to previous state on error
				onError: () => {
					queryClient.setQueryData(queryKey, previousTasks);
				},
			}
		);
	};

	return {
		tasks: data?.tasks,
		isLoading,
		isError,
		error,
		handleDragEnd,
	};
}

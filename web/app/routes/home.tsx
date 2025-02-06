import {
	getTasksQueryKey,
	useTaskChangeStatus,
	useTasks,
} from "@/api/cogniboard";
import { useQueryClient } from "@tanstack/react-query";
import AddTaskDialog from "./project/add.task";
import {
	KanbanBoard,
	KanbanCard,
	KanbanCards,
	KanbanHeader,
	KanbanProvider,
} from "./project/kanban";
import type { DragEndEvent } from "./project/kanban.ts";

const dateFormatter = new Intl.DateTimeFormat("en-US", {
	month: "short",
	day: "numeric",
	year: "numeric",
});

const shortDateFormatter = new Intl.DateTimeFormat("en-US", {
	month: "short",
	day: "numeric",
});

const exampleStatuses = [
	{ id: "1", name: "pending", color: "#6B7280" },
	{ id: "2", name: "in_progress", color: "#F59E0B" },
	{ id: "3", name: "in_review", color: "red" },
	{ id: "4", name: "completed", color: "#10B981" },
];

function useKanbanBoard() {
	const queryClient = useQueryClient();

	const { data, isLoading, isError, error } = useTasks();
	const mutation = useTaskChangeStatus();

	const handleDragEnd = (event: DragEndEvent) => {
		const { active, over } = event;

		if (!over) {
			return;
		}

		const status = exampleStatuses.find((status) => status.name === over.id);

		if (!status) {
			return;
		}

		const taskId = String(active.id);
		const newStatus = status.name;

		// Get the current query data
		const queryKey = getTasksQueryKey();
		const previousTasks = queryClient.getQueryData(queryKey);

		mutation.mutate(
			{ taskId, data: { status: newStatus } },
			{
				onSuccess: () => {
					queryClient.invalidateQueries({ queryKey: queryKey });
				},

				onError: () => {
					queryClient.setQueryData(queryKey, previousTasks);
				},
			},
		);
	};

	return {
		tasks: data?.data.tasks,
		isLoading,
		isError,
		error,
		handleDragEnd,
	};
}

const Home = () => {
	const { tasks, isLoading, isError, error, handleDragEnd } = useKanbanBoard();

	if (isLoading || !tasks) {
		return <div>Loading ...</div>;
	}

	if (isError && error) {
		return (
			<div>
				Error: {error.response?.data.title}{" "}
				{error.response?.data.errors?.map((e) => e.message).join(", ")}
			</div>
		);
	}

	return (
		<main className="p-2">
			<span className="ms-4">
				<AddTaskDialog />
			</span>
			<KanbanProvider onDragEnd={handleDragEnd} className="p-4">
				{exampleStatuses.map((status) => (
					<KanbanBoard key={status.name} id={status.name}>
						<KanbanHeader name={status.name} color={status.color} />
						<KanbanCards>
							{tasks
								.filter((task) => task.status === status.name)
								.map((feature, index) => (
									<KanbanCard
										key={feature.id}
										id={feature.id}
										name={feature.title}
										parent={status.name}
										index={index}
									>
										<div className="flex items-start justify-between gap-2">
											<div className="flex flex-col gap-1">
												<p className="m-0 flex-1 font-bold text-sm">
													{feature.title}
												</p>
												<p className="m-0 text-muted-foreground text-xs">
													{feature.description}
												</p>
												{feature.assignee && (
													<p>Assigned to: {feature.assignee}</p>
												)}
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

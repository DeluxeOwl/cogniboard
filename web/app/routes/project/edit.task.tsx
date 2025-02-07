import { getTasksQueryKey, useTaskEdit } from "@/api/cogniboard";
import { taskEditBody, taskEditBodyTitleMax } from "@/api/cogniboard.zod";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import type { z } from "zod";
import { Loader2 } from "lucide-react";
import { useId, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { TaskDTO } from "@/api/cogniboard";

interface EditTaskDialogProps {
	task: TaskDTO;
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

export default function EditTaskDialog({ task, open, onOpenChange }: EditTaskDialogProps) {
	const id = useId();

	const { form, mutation, onSubmit } = useEditTask({
		task,
		onSuccess: () => onOpenChange(false),
	});

	const {
		register,
		formState: { errors },
		watch,
	} = form;

	const title = watch("title");
	const charactersLeft = taskEditBodyTitleMax - (title?.length || 0);

	return (
		<Dialog
			open={open}
			onOpenChange={onOpenChange}
		>
			<DialogContent className="flex flex-col gap-0 overflow-y-visible p-0 sm:max-w-lg [&>button:last-child]:top-3.5">
				<DialogHeader className="contents space-y-0 text-left">
					<DialogTitle className="border-b border-border px-6 py-4 text-base">
						Edit task
					</DialogTitle>
				</DialogHeader>
				<DialogDescription className="sr-only">
					Edit task details. Only the title is required.
				</DialogDescription>
				<div className="overflow-y-auto">
					<div className="px-6 pb-6 pt-4">
						<form
							className="space-y-4"
							onSubmit={onSubmit}
						>
							<div className="flex flex-col gap-4 sm:flex-row">
								<div className="flex-1 space-y-2">
									<Label htmlFor={`${id}-task-title`}>Task title</Label>
									<Input
										id={`${id}-task-title`}
										placeholder="TODO: Fix this app"
										type="text"
										maxLength={taskEditBodyTitleMax}
										{...register("title")}
										aria-invalid={errors.title ? "true" : "false"}
									/>
									{errors.title ? (
										<p className="text-sm text-destructive">{errors.title.message}</p>
									) : (
										<p
											id={`${id}-task-title`}
											className="mt-2 text-right text-xs text-muted-foreground"
											role="status"
											aria-live="polite"
										>
											<span className="tabular-nums">{charactersLeft}</span> characters left
										</p>
									)}
								</div>
							</div>

							<div className="space-y-2">
								<Label htmlFor={`${id}-description`}>Description</Label>
								<Textarea
									id={`${id}-description`}
									placeholder="Write a few sentences about this task"
									aria-describedby={`${id}-description`}
									{...register("description")}
									aria-invalid={errors.description ? "true" : "false"}
								/>
								{errors.description && (
									<p className="text-sm text-destructive">{errors.description.message}</p>
								)}
							</div>
						</form>
						{mutation.error && (
							<div className="mt-4 rounded-md bg-destructive/15 p-3 text-sm text-destructive">
								{mutation.error.response?.data.title || "Failed to edit task. Please try again."}
							</div>
						)}
					</div>
				</div>
				<DialogFooter className="border-t border-border px-6 py-4">
					<DialogClose asChild={true}>
						<Button
							type="button"
							variant="outline"
							disabled={mutation.isPending}
						>
							Cancel
						</Button>
					</DialogClose>
					<Button
						type="submit"
						onClick={onSubmit}
						disabled={mutation.isPending}
					>
						{mutation.isPending ? (
							<>
								<Loader2 className="size-4 me-2 animate-spin" />
								Editing ...
							</>
						) : (
							"Edit task"
						)}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}

type FormData = z.infer<typeof taskEditBody>;

interface UseEditTaskProps {
	task: TaskDTO;
	onSuccess?: () => void;
}

export function useEditTask({ task, onSuccess }: UseEditTaskProps) {
	const queryClient = useQueryClient();
	const form = useForm<FormData>({
		resolver: zodResolver(taskEditBody),
		defaultValues: {
			title: task.title,
			description: task.description || "",
		},
	});

	const mutation = useTaskEdit();

	const onSubmit = form.handleSubmit((data) => {
		mutation.mutate(
			{
				taskId: task.id,
				data: {
					title: data.title,
					description: data.description || undefined,
				},
			},
			{
				onSuccess: () => {
					queryClient.invalidateQueries({ queryKey: getTasksQueryKey() });
					onSuccess?.();
					form.reset();
				},
			}
		);
	});

	return {
		form,
		mutation,
		onSubmit,
	};
}

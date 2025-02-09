import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import type { z } from "zod";

import { taskCreateMutationRequestSchema, tasksQueryKey } from "@/api";
import { convertFilesToFormDataFormat } from "@/lib/form-data";
import { useTaskCreate } from "@/api/hooks/useTaskCreate";
import { Dropzone, DropzoneContent, DropzoneEmptyState } from "@/components/ui/dropzone";
import { useQueryClient } from "@tanstack/react-query";
import { Loader2, LucidePlus } from "lucide-react";
import { useId, useState } from "react";

export default function AddTaskDialog() {
	const id = useId();
	const [open, setOpen] = useState(false);

	const { form, mutation, onSubmit } = useAddTask({
		onSuccess: () => setOpen(false),
	});

	const {
		register,
		formState: { errors },
		watch,
		setValue,
		setError,
	} = form;

	const title = watch("title");
	const files = watch("files");
	const taskCreateBodyTitleMax =
		taskCreateMutationRequestSchema.shape.title._def.checks.find((check) => check.kind === "max")
			?.value ?? 50;
	const charactersLeft = taskCreateBodyTitleMax - (title?.length || 0);

	return (
		<Dialog
			open={open}
			onOpenChange={setOpen}
		>
			<DialogTrigger asChild={true}>
				<Button>
					<LucidePlus className="size-4 me-2" />
					Add a task
				</Button>
			</DialogTrigger>
			<DialogContent className="flex flex-col gap-0 overflow-y-visible p-0 sm:max-w-lg [&>button:last-child]:top-3.5">
				<DialogHeader className="contents space-y-0 text-left">
					<DialogTitle className="border-b border-border px-6 py-4 text-base">
						Add a new task
					</DialogTitle>
				</DialogHeader>
				<DialogDescription className="sr-only">
					Create a task, only the title is required.
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
										maxLength={taskCreateBodyTitleMax}
										{...register("title")}
										aria-invalid={errors.title ? "true" : "false"}
									/>
									{errors.title ? (
										<p className="text-sm text-destructive">{errors.title.message?.toString()}</p>
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
									<p className="text-sm text-destructive">
										{errors.description.message?.toString()}
									</p>
								)}
							</div>

							<div className="space-y-2">
								<Label htmlFor={`${id}-files`}>Files</Label>
								<Dropzone
									id={`${id}-files`}
									maxSize={1024 * 1024 * 50}
									minSize={1024}
									maxFiles={10}
									accept={{ "image/*": [], "application/pdf": [".pdf"], "text/csv": [".csv"] }}
									onDrop={(acceptedFiles) => {
										setValue("files", acceptedFiles);
									}}
									src={files}
									onError={(error) => {
										setError("files", {
											message: error.message,
										});
									}}
								>
									<DropzoneEmptyState />
									<DropzoneContent />
								</Dropzone>
								{errors.files && (
									<p className="text-sm text-destructive">{errors.files.message?.toString()}</p>
								)}
							</div>
						</form>
						{mutation.error && (
							<div className="mt-4 rounded-md bg-destructive/15 p-3 text-sm text-destructive">
								{mutation.error.response?.data.message ||
									"Failed to create task. Please try again."}
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
								Adding...
							</>
						) : (
							"Add task"
						)}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}

type FormData = z.infer<typeof taskCreateMutationRequestSchema>;

interface UseAddTaskProps {
	onSuccess?: () => void;
}

export function useAddTask({ onSuccess }: UseAddTaskProps = {}) {
	const queryClient = useQueryClient();
	const form = useForm<FormData>({
		resolver: zodResolver(taskCreateMutationRequestSchema),
		defaultValues: {
			title: "Water the flowers",
			description: "",
			files: [],
		},
	});

	const mutation = useTaskCreate();

	const onSubmit = form.handleSubmit((data) => {
		const formattedData = convertFilesToFormDataFormat(data);

		mutation.mutate(
			{
				data: formattedData,
			},
			{
				onSuccess: () => {
					queryClient.invalidateQueries({ queryKey: tasksQueryKey() });
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

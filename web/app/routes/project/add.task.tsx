import { useCreateTask } from "@/api/cogniboard";
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
import { useCharacterLimit } from "@/hooks/use-character-limit";
import { Loader2, LucidePlus } from "lucide-react";
import { useId, useState } from "react";

export default function AddTaskDialog() {
	const id = useId();
	const [open, setOpen] = useState(false);

	const mutation = useCreateTask();

	const maxLength = 50;
	const {
		value,
		characterCount,
		handleChange,
		maxLength: limit,
	} = useCharacterLimit({
		maxLength,
		initialValue: "Water the flowers",
	});

	const [description, setDescription] = useState("");

	const handleCreateTask = () => {
		mutation.mutate(
			{
				data: {
					title: value,
					description: description || undefined,
				},
			},
			{
				onSuccess: () => {
					setOpen(false);
				},
			}
		);
	};

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
						<form className="space-y-4">
							<div className="flex flex-col gap-4 sm:flex-row">
								<div className="flex-1 space-y-2">
									<Label htmlFor={`${id}-task-title`}>Task title</Label>
									<Input
										id={`${id}-task-title`}
										placeholder="TODO: Fix this app"
										defaultValue={value}
										type="text"
										maxLength={maxLength}
										onChange={handleChange}
										required={true}
									/>
									<p
										id={`${id}-task-title`}
										className="mt-2 text-right text-xs text-muted-foreground"
										role="status"
										aria-live="polite"
									>
										<span className="tabular-nums">{limit - characterCount}</span> characters left
									</p>
								</div>
							</div>

							<div className="space-y-2">
								<Label htmlFor={`${id}-description`}>Description</Label>
								<Textarea
									id={`${id}-description`}
									value={description}
									onChange={(e) => setDescription(e.target.value)}
									placeholder="Write a few sentences about this task"
									aria-describedby={`${id}-description`}
								/>
							</div>
						</form>
						{mutation.error && (
							<div className="mt-4 rounded-md bg-destructive/15 p-3 text-sm text-destructive">
								{mutation.error?.detail || "Failed to create task. Please try again."}
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
						type="button"
						onClick={handleCreateTask}
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
				{mutation.error ? <p>Got an error: {JSON.stringify(mutation.error)}</p> : null}
			</DialogContent>
		</Dialog>
	);
}

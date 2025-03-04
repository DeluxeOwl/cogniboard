import { tasksQueryKey } from "@/api";
import { Thread } from "@/components/assistant-ui/thread";
import {
	AssistantRuntimeProvider,
	useLocalRuntime,
	type ChatModelAdapter,
} from "@assistant-ui/react";
import { useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { z } from "zod";

export default function Chat() {
	const queryClient = useQueryClient();
	const [adapter] = useState(() =>
		CreateOpenAIWithoutProxyAdapter({
			onRefetch: async () => {
				await queryClient.invalidateQueries({
					queryKey: tasksQueryKey(),
					type: "all",
				});
				await queryClient.refetchQueries({
					queryKey: tasksQueryKey(),
					type: "all",
				});
			},
		})
	);
	const runtime = useLocalRuntime(adapter);
	return (
		<AssistantRuntimeProvider runtime={runtime}>
			<div className="h-full pb-12">
				<Thread />
			</div>
		</AssistantRuntimeProvider>
	);
}

const DeltaSchema = z.object({
	content: z.string().nullable().optional(),
});

const ChoiceSchema = z.object({
	index: z.number(),
	delta: DeltaSchema,
	finish_reason: z.string().nullable(),
});

const ChatResponseSchema = z.object({
	id: z.string(),
	object: z.string(),
	created: z.number(),
	model: z.string(),
	choices: z.array(ChoiceSchema).optional(),
});

type CreateOpenAIWithoutProxyAdapterProps = {
	onRefetch?: () => void;
};

function CreateOpenAIWithoutProxyAdapter(
	props: CreateOpenAIWithoutProxyAdapterProps
): ChatModelAdapter {
	return {
		async *run({ messages, abortSignal, context }) {
			try {
				const response = await fetch(`${BaseURL}/v1/api/chat`, {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
						Accept: "text/event-stream",
					},
					body: JSON.stringify({
						messages: messages.map((msg) => ({
							role: msg.role,
							content: msg.content,
						})),
					}),
					signal: abortSignal,
				});

				if (!response.ok) {
					throw new Error(`HTTP error! status: ${response.status}`);
				}

				const reader = response.body?.getReader();
				if (!reader) throw new Error("No response body");

				const decoder = new TextDecoder();
				let buffer = "";
				let accumulatedText = "";

				while (true) {
					const { done, value } = await reader.read();
					if (done) break;

					const chunk = decoder.decode(value);
					buffer += chunk;

					const lines = buffer.split("\n");
					buffer = lines.pop() || "";

					for (const line of lines) {
						try {
							const parsedData = ChatResponseSchema.parse(JSON.parse(line));

							if (parsedData.choices?.[0].finish_reason !== null) {
								break;
							}

							const content = parsedData.choices?.[0]?.delta?.content;
							if (content) {
								accumulatedText += content;

								const hasRefetch = accumulatedText.endsWith("@refetch");

								if (hasRefetch) {
									props.onRefetch?.();
								}

								const displayText = hasRefetch
									? accumulatedText.slice(0, -8) // remove "@refetch"
									: accumulatedText;

								yield {
									content: [{ type: "text", text: displayText }],
								};
							}
						} catch (e) {
							console.error("Error parsing or validating line:", line, e);
						}
					}
				}
			} catch (error) {
				console.error("Error in chat request:", error);
				throw error;
			}
		},
	};
}

const AIModel = "o3-mini";
const BaseURL = "http://127.0.0.1:8888";
// Use this if you're just proxying openai calls, see proxy.go
const OpenAIProxyAdapter: ChatModelAdapter = {
	async *run({ messages, abortSignal, context }) {
		try {
			const response = await fetch(`${BaseURL}/chat/v1/chat/completions`, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
					Accept: "text/event-stream",
				},
				body: JSON.stringify({
					model: AIModel,
					stream: true,
					messages: messages.map((msg) => ({
						role: msg.role,
						content: msg.content,
					})),
				}),
				signal: abortSignal,
			});

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const reader = response.body?.getReader();
			if (!reader) throw new Error("No response body");

			const decoder = new TextDecoder();
			let buffer = "";
			let accumulatedText = ""; // Add this to accumulate the full response

			while (true) {
				const { done, value } = await reader.read();
				if (done) break;

				const chunk = decoder.decode(value);
				buffer += chunk;

				const lines = buffer.split("\n");
				buffer = lines.pop() || "";

				for (const line of lines) {
					if (line.trim() === "") continue;
					if (line.trim() === "data: [DONE]") break;

					if (line.startsWith("data: ")) {
						try {
							const jsonData = JSON.parse(line.slice(6));
							const content = jsonData.choices?.[0]?.delta?.content;
							if (content) {
								accumulatedText += content; // Accumulate the text
								yield {
									content: [{ type: "text", text: accumulatedText }],
								};
							}
						} catch (e) {
							console.error("Error parsing line:", line, e);
						}
					}
				}
			}
		} catch (error) {
			console.error("Error in chat request:", error);
			throw error;
		}
	},
};

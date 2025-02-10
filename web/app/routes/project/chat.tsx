import { Thread } from "@/components/assistant-ui/thread";
export default function Chat() {
	const runtime = useLocalRuntime(MyModelAdapter);
	return (
		<AssistantRuntimeProvider runtime={runtime}>
			<div className="h-full pb-12">
				<Thread />
			</div>
		</AssistantRuntimeProvider>
	);
}

import {
	AssistantRuntimeProvider,
	useLocalRuntime,
	type ChatModelAdapter,
} from "@assistant-ui/react";
import type { ReactNode } from "react";

const AIModel = "o3-mini";
const BaseChatURL = "http://127.0.0.1:8888/chat";
const MyModelAdapter: ChatModelAdapter = {
	async *run({ messages, abortSignal, context }) {
		try {
			const response = await fetch(`${BaseChatURL}/v1/chat/completions`, {
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

export function MyRuntimeProvider({
	children,
}: Readonly<{
	children: ReactNode;
}>) {
	const runtime = useLocalRuntime(MyModelAdapter);

	return <AssistantRuntimeProvider runtime={runtime}>{children}</AssistantRuntimeProvider>;
}

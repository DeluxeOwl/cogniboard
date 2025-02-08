import { defineConfig } from "@kubb/core";
import { pluginOas } from "@kubb/plugin-oas";
import { pluginReactQuery } from "@kubb/plugin-react-query";
import { pluginTs } from "@kubb/plugin-ts";
import { pluginZod } from "@kubb/plugin-zod";

export default defineConfig({
	root: ".",
	input: {
		path: "../openapi3.yaml",
	},
	output: {
		path: "./api",
		clean: true,
	},
	plugins: [
		pluginOas(),
		pluginTs(),
		pluginZod(),
		pluginReactQuery({
			client: {
				baseURL: "http://127.0.0.1:8888/v1/api",
			},
		}),
	],
});

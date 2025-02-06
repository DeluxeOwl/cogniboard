import { defineConfig } from "orval";

export default defineConfig({
	cogniboard: {
		input: "../openapi3.yaml",
		output: {
			mode: "split",
			client: "react-query",
			httpClient: "fetch",
			target: "api/cogniboard.ts",
			baseUrl: {
				getBaseUrlFromSpecification: true,
			},
		},
	},
	cogniboardZod: {
		input: "../openapi3.yaml",
		output: {
			mode: "split",
			client: "zod",
			target: "api/cogniboard.zod.ts",
			baseUrl: {
				getBaseUrlFromSpecification: true,
			},
		},
	},
});

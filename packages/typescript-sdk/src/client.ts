import createClient from "openapi-fetch";

import type { paths } from "./types.ts";

export function createApiClient(baseUrl: string) {
	return createClient<paths>({ baseUrl });
}

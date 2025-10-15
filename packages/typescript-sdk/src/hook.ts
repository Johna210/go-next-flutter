import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createApiClient } from "./client.ts";
import type { paths } from "./types.ts";

const API_URL = "http://localhost:8080";
const api = createApiClient(API_URL);

type UserListParams = paths["/api/v1/users"]["get"]["parameters"]["query"];

export const queryKeys = {
	users: {
		all: ["users"] as const,
		lists: () => [...queryKeys.users.all, "list"] as const,
		list: (filters: UserListParams) => [...queryKeys.users.lists(), filters] as const,
		details: () => [...queryKeys.users.all, "detail"] as const,
		detail: (id: string) => [...queryKeys.users.details(), id] as const,
	},
	// Add more resource keys...
} as const;

export function useUsers(params?: { limit?: number; offset?: number }) {
	return useQuery({
		queryKey: queryKeys.users.list(params || {}),
		queryFn: async () => {
			const { data, error } = await api.GET("/api/v1/users", {
				params: { query: params },
			});
			if (error) {
				throw new Error(error.title);
			}
			return data;
		},
	});
}

export function useUser(userId: string) {
	return useQuery({
		queryKey: queryKeys.users.detail(userId),
		queryFn: async () => {
			const { data, error } = await api.GET("/api/v1/users/{id}", {
				params: { path: { id: userId } },
			});
			if (error) {
				throw new Error(error.title);
			}
			return data;
		},
		enabled: Boolean(userId),
	});
}

export function useCreateUser() {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (userData: { email: string; name: string }) => {
			const { data, error } = await api.POST("/api/v1/users", {
				body: userData,
			});
			if (error) {
				throw new Error(error.title);
			}
			return data;
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: queryKeys.users.lists() });
		},
	});
}

export { api };

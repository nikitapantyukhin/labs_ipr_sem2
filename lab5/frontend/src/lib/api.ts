import { env } from "@/env";

const API_URL = env.NEXT_PUBLIC_API_URL.replace(/\/$/, "");

export type UserProfile = {
	id: number;
	full_name: string;
	social_network_link: string;
	phone_number: string;
	email: string;
	birth_date: string;
	role: string;
	group_name?: string;
	access_token?: string;
	refresh_token?: string;
};

export type Club = {
	id: number;
	name: string;
	description: string;
	sport_type_id: number;
	sport_type_name: string;
	teacher_id: number;
	teacher_name: string;
	total_places?: number;
	place: string;
	education_level_id: number;
	education_level_name: string;
	required_workout_per_week: number;
	attachments?: string[];
	created_at: string;
};

export type Workout = {
	id: number;
	club_id: number;
	start_date: string;
	end_date: string;
	cancelled: boolean;
	created_at: string;
	updated_at: string;
};

export type RegisterPayload = {
	full_name: string;
	social_network_link: string;
	phone_number: string;
	email: string;
	birth_date: string;
	password: string;
	group_id: number | null;
};

export class ApiError extends Error {
	status: number;

	constructor(message: string, status: number) {
		super(message);
		this.name = "ApiError";
		this.status = status;
	}
}

type RequestOptions = RequestInit & {
	token?: string;
};

async function request<T>(path: string, options: RequestOptions = {}) {
	const headers = new Headers(options.headers);
	headers.set("Accept", "application/json");

	if (!(options.body instanceof FormData) && !headers.has("Content-Type")) {
		headers.set("Content-Type", "application/json");
	}

	if (options.token) {
		headers.set("Authorization", `Bearer ${options.token}`);
	}

	const response = await fetch(`${API_URL}${path}`, {
		...options,
		headers,
	});

	if (!response.ok) {
		let message = "Не удалось выполнить запрос";

		try {
			const data = (await response.json()) as { message?: string };
			if (data.message) {
				message = data.message;
			}
		} catch {
			message = response.statusText || message;
		}

		throw new ApiError(message, response.status);
	}

	if (response.status === 204) {
		return undefined as T;
	}

	return (await response.json()) as T;
}

export function registerUser(payload: RegisterPayload) {
	return request<UserProfile>("/users/create", {
		method: "POST",
		body: JSON.stringify(payload),
	});
}

export function getCurrentUser(token: string) {
	return request<UserProfile>("/users/", { token });
}

export async function getClubs(token: string) {
	const response = await request<{ clubs: Club[] }>("/clubs/", { token });
	return response.clubs;
}

export async function getClubWorkouts(token: string, clubId: number) {
	const response = await request<{ workouts: Workout[] }>(
		`/clubs/${clubId}/workouts`,
		{ token },
	);
	return response.workouts;
}

export function createJoinRequest(token: string, clubId: number) {
	return request<{ id: number; club_id: number; user_id: number; status: string }>(
		"/club_join_requests/create",
		{
			method: "POST",
			token,
			body: JSON.stringify({ club_id: clubId }),
		},
	);
}

import { env } from "@/env";
import type { DefaultSession, NextAuthConfig } from "next-auth";
import CredentialsProvider from "next-auth/providers/credentials";
import { z } from "zod";

type BackendAuthResponse = {
	id: number;
	full_name: string;
	email: string;
	role: string;
	access_token: string;
	refresh_token: string;
};

type AppToken = {
	id?: string;
	sub?: string;
	role?: string;
	fullName?: string;
	accessToken?: string;
	refreshToken?: string;
};

declare module "next-auth" {
	interface Session extends DefaultSession {
		accessToken?: string;
		refreshToken?: string;
		user: {
			id: string;
			role?: string;
			fullName?: string;
		} & DefaultSession["user"];
	}

	interface User {
		fullName?: string;
		role?: string;
		accessToken?: string;
		refreshToken?: string;
	}
}

const credentialsSchema = z.object({
	email: z.string().email(),
	password: z.string().min(1),
});

export const authConfig = {
	session: {
		strategy: "jwt",
	},
	providers: [
		CredentialsProvider({
			name: "Email and password",
			credentials: {
				email: { label: "Email", type: "email" },
				password: { label: "Password", type: "password" },
			},
			async authorize(credentials) {
				const parsedCredentials = credentialsSchema.safeParse(credentials);
				if (!parsedCredentials.success) {
					return null;
				}

				const response = await fetch(`${env.BACKEND_API_URL}/users/login`, {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(parsedCredentials.data),
				});

				if (!response.ok) {
					return null;
				}

				const user = (await response.json()) as BackendAuthResponse;

				return {
					id: String(user.id),
					name: user.full_name,
					email: user.email,
					fullName: user.full_name,
					role: user.role,
					accessToken: user.access_token,
					refreshToken: user.refresh_token,
				};
			},
		}),
	],
	callbacks: {
		jwt: ({ token, user }) => {
			const appToken = token as AppToken;

			if (user) {
				appToken.id = user.id;
				appToken.role = user.role;
				appToken.fullName = user.fullName;
				appToken.accessToken = user.accessToken;
				appToken.refreshToken = user.refreshToken;
			}

			return appToken;
		},
		session: ({ session, token }) => {
			const appToken = token as AppToken;

			return {
				...session,
				accessToken: appToken.accessToken,
				refreshToken: appToken.refreshToken,
				user: {
					...session.user,
					id: appToken.id ?? appToken.sub ?? "",
					role: appToken.role,
					fullName: appToken.fullName,
				},
			};
		},
	},
	pages: {
		signIn: "/login",
	},
} satisfies NextAuthConfig;

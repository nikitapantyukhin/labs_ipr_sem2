"use client";

import AuthForm, { type AuthField } from "@/components/AuthForm";
import { signIn } from "next-auth/react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

const fields: AuthField[] = [
	{
		name: "email",
		label: "Email",
		type: "email",
		placeholder: "you@example.com",
		autoComplete: "email",
	},
	{
		name: "password",
		label: "Пароль",
		type: "password",
		placeholder: "Минимум 8 символов",
		autoComplete: "current-password",
	},
];

export default function LoginPage() {
	const router = useRouter();
	const [error, setError] = useState("");

	const handleLogin = async (data: Record<string, string>) => {
		setError("");

		const result = await signIn("credentials", {
			email: data.email,
			password: data.password,
			redirect: false,
		}).catch(() => null);

		if (!result?.ok) {
			setError("Неверный email или пароль");
			return;
		}

		router.push("/");
		router.refresh();
	};

	return (
		<main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
			<AuthForm
				title="Вход"
				description="Используйте аккаунт платформы спортивных секций."
				fields={fields}
				submitLabel="Войти"
				action={handleLogin}
				error={error}
			>
				<p className="text-center text-slate-600 text-sm">
					Нет аккаунта?{" "}
					<Link
						href="/register"
						className="font-medium text-emerald-700 hover:text-emerald-800"
					>
						Зарегистрироваться
					</Link>
				</p>
			</AuthForm>
		</main>
	);
}

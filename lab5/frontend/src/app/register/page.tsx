"use client";

import AuthForm, { type AuthField } from "@/components/AuthForm";
import { ApiError, registerUser } from "@/lib/api";
import { signIn } from "next-auth/react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

const fields: AuthField[] = [
	{
		name: "fullName",
		label: "ФИО",
		type: "text",
		placeholder: "Иван Петров",
		autoComplete: "name",
	},
	{
		name: "telegram",
		label: "Telegram",
		type: "text",
		placeholder: "@username",
		autoComplete: "username",
	},
	{
		name: "birthDate",
		label: "Дата рождения",
		type: "date",
		autoComplete: "bday",
	},
	{
		name: "phone",
		label: "Телефон",
		type: "tel",
		placeholder: "+7 (999) 000-00-00",
		autoComplete: "tel",
	},
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
		autoComplete: "new-password",
	},
	{
		name: "confirmPassword",
		label: "Повторите пароль",
		type: "password",
		placeholder: "Повторите пароль",
		autoComplete: "new-password",
	},
];

export default function RegisterPage() {
	const router = useRouter();
	const [error, setError] = useState("");

	const handleRegister = async (data: Record<string, string>) => {
		setError("");

		if (data.password !== data.confirmPassword) {
			setError("Пароли не совпадают");
			return;
		}

		if ((data.password ?? "").length < 8) {
			setError("Пароль должен быть не короче 8 символов");
			return;
		}

		try {
			await registerUser({
				full_name: data.fullName ?? "",
				social_network_link: normalizeTelegram(data.telegram ?? ""),
				phone_number: normalizePhone(data.phone ?? ""),
				email: data.email ?? "",
				birth_date: toApiDate(data.birthDate ?? ""),
				password: data.password ?? "",
				group_id: null,
			});

			const result = await signIn("credentials", {
				email: data.email,
				password: data.password,
				redirect: false,
			}).catch(() => null);

			if (!result?.ok) {
				router.push("/login");
				return;
			}

			router.push("/");
			router.refresh();
		} catch (caughtError) {
			if (caughtError instanceof ApiError) {
				setError(translateApiError(caughtError.message));
				return;
			}

			setError("Не удалось зарегистрироваться. Проверьте подключение к API.");
		}
	};

	return (
		<main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-10">
			<AuthForm
				title="Регистрация"
				description="Создайте аккаунт студента для просмотра секций и подачи заявок."
				fields={fields}
				submitLabel="Создать аккаунт"
				action={handleRegister}
				error={error}
			>
				<p className="text-center text-slate-600 text-sm">
					Уже есть аккаунт?{" "}
					<Link
						href="/login"
						className="font-medium text-emerald-700 hover:text-emerald-800"
					>
						Войти
					</Link>
				</p>
			</AuthForm>
		</main>
	);
}

function normalizePhone(value: string) {
	let digits = value.replace(/\D/g, "");
	if (digits.startsWith("8")) {
		digits = `7${digits.slice(1)}`;
	}
	if (!digits.startsWith("7")) {
		digits = `7${digits}`;
	}
	return `+${digits.slice(0, 11)}`;
}

function normalizeTelegram(value: string) {
	const trimmed = value.trim();
	if (!trimmed) {
		return "";
	}
	return trimmed.startsWith("@") ? trimmed : `@${trimmed}`;
}

function toApiDate(value: string) {
	if (!value) {
		return "";
	}
	return `${value}T00:00:00Z`;
}

function translateApiError(message: string) {
	if (message.includes("already exists")) {
		return "Пользователь с такими данными уже существует";
	}
	if (message.includes("Invalid") || message.includes("can't parse")) {
		return "Проверьте корректность заполненных полей";
	}
	return message;
}

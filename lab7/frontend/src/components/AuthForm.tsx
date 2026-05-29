"use client";

import { useEffect, useState, type FormEvent, type ReactNode } from "react";

export type AuthField = {
	name: string;
	label: string;
	type: string;
	placeholder?: string;
	autoComplete?: string;
};

type AuthFormProps = {
	title: string;
	description?: string;
	fields: AuthField[];
	submitLabel: string;
	action: (data: Record<string, string>) => Promise<void> | void;
	children?: ReactNode;
	error?: string;
};

export default function AuthForm({
	title,
	description,
	fields,
	submitLabel,
	action,
	children,
	error,
}: AuthFormProps) {
	const [formState, setFormState] = useState<Record<string, string>>(
		fields.reduce<Record<string, string>>((acc, field) => {
			acc[field.name] = "";
			return acc;
		}, {}),
	);
	const [passwordError, setPasswordError] = useState("");
	const [isSubmitting, setIsSubmitting] = useState(false);

	const handleChange = (name: string, value: string) => {
		let nextValue = value;

		if (name === "phone") {
			nextValue = formatPhone(value);
		}

		if (name === "telegram" && value && !value.startsWith("@")) {
			nextValue = `@${value.replace(/^@+/, "")}`;
		}

		setFormState((prev) => ({ ...prev, [name]: nextValue }));
	};

	useEffect(() => {
		if ("password" in formState && "confirmPassword" in formState) {
			if (
				formState.confirmPassword &&
				formState.password !== formState.confirmPassword
			) {
				setPasswordError("Пароли не совпадают");
			} else {
				setPasswordError("");
			}
		}
	}, [formState.password, formState.confirmPassword]);

	const handleSubmit = async (event: FormEvent) => {
		event.preventDefault();
		if (passwordError) {
			return;
		}

		setIsSubmitting(true);
		try {
			await action(formState);
		} finally {
			setIsSubmitting(false);
		}
	};

	return (
		<div className="w-full max-w-md rounded-lg border border-slate-200 bg-white p-8 shadow-sm">
			<div className="mb-6">
				<h1 className="font-semibold text-2xl text-slate-950">{title}</h1>
				{description ? (
					<p className="mt-2 text-sm text-slate-600">{description}</p>
				) : null}
			</div>
			<form onSubmit={handleSubmit} className="space-y-4">
				{fields.map((field) => (
					<div key={field.name}>
						<label
							htmlFor={field.name}
							className="mb-1 block font-medium text-slate-800 text-sm"
						>
							{field.label}
						</label>
						<input
							id={field.name}
							type={field.type}
							name={field.name}
							placeholder={field.placeholder}
							autoComplete={field.autoComplete}
							value={formState[field.name] ?? ""}
							onChange={(event) =>
								handleChange(field.name, event.target.value)
							}
							className="w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-slate-950 text-sm outline-none transition focus:border-emerald-600 focus:ring-2 focus:ring-emerald-100"
							required
						/>
						{field.name === "confirmPassword" && passwordError ? (
							<p className="mt-1 text-red-600 text-sm">{passwordError}</p>
						) : null}
					</div>
				))}

				{error ? (
					<div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-red-700 text-sm">
						{error}
					</div>
				) : null}

				<button
					type="submit"
					disabled={isSubmitting}
					className="inline-flex w-full items-center justify-center rounded-md bg-emerald-700 px-4 py-2.5 font-medium text-sm text-white transition hover:bg-emerald-800 disabled:cursor-not-allowed disabled:bg-slate-400"
				>
					{isSubmitting ? "Отправка..." : submitLabel}
				</button>
			</form>
			{children ? <div className="mt-5">{children}</div> : null}
		</div>
	);
}

function formatPhone(value: string) {
	let digits = value.replace(/\D/g, "");
	if (digits.startsWith("8")) {
		digits = `7${digits.slice(1)}`;
	}
	if (!digits.startsWith("7")) {
		digits = `7${digits}`;
	}
	digits = digits.slice(0, 11);

	if (digits.length <= 1) {
		return "+7";
	}

	const code = digits.slice(1, 4);
	const first = digits.slice(4, 7);
	const second = digits.slice(7, 9);
	const third = digits.slice(9, 11);

	let formatted = "+7";
	if (code) {
		formatted += ` (${code}`;
	}
	if (code.length === 3) {
		formatted += ")";
	}
	if (first) {
		formatted += ` ${first}`;
	}
	if (second) {
		formatted += `-${second}`;
	}
	if (third) {
		formatted += `-${third}`;
	}

	return formatted;
}

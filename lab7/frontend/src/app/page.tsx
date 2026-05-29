"use client";

import {
	createJoinRequest,
	getClubWorkouts,
	getClubs,
	getCurrentUser,
	type Club,
} from "@/lib/api";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
	Activity,
	CalendarDays,
	Dumbbell,
	LogIn,
	LogOut,
	MapPin,
	ShieldCheck,
	UserRound,
	Users,
} from "lucide-react";
import { signOut, useSession } from "next-auth/react";
import Link from "next/link";
import { useMemo, useState, type ReactNode } from "react";

export default function HomePage() {
	const { data: session, status } = useSession();
	const token = session?.accessToken ?? "";
	const [selectedClubId, setSelectedClubId] = useState<number | null>(null);
	const [notice, setNotice] = useState("");

	const profileQuery = useQuery({
		queryKey: ["profile"],
		queryFn: () => getCurrentUser(token),
		enabled: Boolean(token),
	});

	const clubsQuery = useQuery({
		queryKey: ["clubs"],
		queryFn: () => getClubs(token),
		enabled: Boolean(token),
	});

	const selectedClub = useMemo(
		() => clubsQuery.data?.find((club) => club.id === selectedClubId) ?? null,
		[clubsQuery.data, selectedClubId],
	);

	const workoutsQuery = useQuery({
		queryKey: ["club-workouts", selectedClubId],
		queryFn: () => getClubWorkouts(token, selectedClubId ?? 0),
		enabled: Boolean(token && selectedClubId),
	});

	const joinMutation = useMutation({
		mutationFn: (clubId: number) => createJoinRequest(token, clubId),
		onSuccess: () => setNotice("Заявка отправлена. Преподаватель увидит ее в списке заявок."),
		onError: () => setNotice("Не удалось отправить заявку. Возможно, она уже была создана."),
	});

	if (status === "loading") {
		return (
			<main className="flex min-h-screen items-center justify-center bg-slate-50 text-slate-700">
				Загрузка...
			</main>
		);
	}

	if (!session) {
		return <SignedOutView />;
	}

	const profile = profileQuery.data;
	const clubs = clubsQuery.data ?? [];
	const totalWorkouts = workoutsQuery.data?.length ?? 0;

	return (
		<main className="min-h-screen bg-[#f7f8f4] text-slate-950">
			<header className="border-slate-200 border-b bg-white">
				<div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-5 sm:flex-row sm:items-center sm:justify-between lg:px-8">
					<div>
						<p className="font-medium text-emerald-700 text-sm">Sport Platform</p>
						<h1 className="font-semibold text-2xl text-slate-950">
							Спортивные секции университета
						</h1>
					</div>
					<div className="flex flex-wrap items-center gap-3">
						<div className="inline-flex items-center gap-2 rounded-md border border-slate-200 bg-slate-50 px-3 py-2 text-slate-700 text-sm">
							<UserRound className="h-4 w-4" />
							{profile?.full_name ?? session.user.fullName ?? session.user.email}
						</div>
						<button
							type="button"
							onClick={() => signOut({ callbackUrl: "/login" })}
							className="inline-flex items-center gap-2 rounded-md border border-slate-300 bg-white px-3 py-2 font-medium text-slate-800 text-sm transition hover:bg-slate-100"
						>
							<LogOut className="h-4 w-4" />
							Выйти
						</button>
					</div>
				</div>
			</header>

			<section className="mx-auto grid max-w-7xl gap-4 px-4 py-6 sm:grid-cols-3 lg:px-8">
				<Metric
					icon={<ShieldCheck className="h-5 w-5" />}
					label="Роль"
					value={profile?.role ?? session.user.role ?? "Student"}
				/>
				<Metric
					icon={<Dumbbell className="h-5 w-5" />}
					label="Доступные секции"
					value={String(clubs.length)}
				/>
				<Metric
					icon={<CalendarDays className="h-5 w-5" />}
					label="Тренировки выбранной секции"
					value={selectedClub ? String(totalWorkouts) : "0"}
				/>
			</section>

			<section className="mx-auto grid max-w-7xl gap-6 px-4 pb-10 lg:grid-cols-[minmax(0,1fr)_360px] lg:px-8">
				<div>
					<div className="mb-4 flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
						<div>
							<h2 className="font-semibold text-xl">Каталог секций</h2>
							<p className="text-slate-600 text-sm">
								Выберите секцию, посмотрите расписание и отправьте заявку.
							</p>
						</div>
						{clubsQuery.isFetching ? (
							<span className="text-slate-500 text-sm">Обновление...</span>
						) : null}
					</div>

					{clubsQuery.isError ? (
						<StateMessage
							title="API недоступен"
							text="Проверьте, что бэкенд запущен на порту 8080."
						/>
					) : null}

					{!clubsQuery.isLoading && clubs.length === 0 ? (
						<StateMessage
							title="Секций пока нет"
							text="В docker-сборке включен демо-seed, поэтому после первого запуска здесь появятся примеры."
						/>
					) : null}

					<div className="grid gap-4 md:grid-cols-2">
						{clubs.map((club) => (
							<ClubCard
								key={club.id}
								club={club}
								isSelected={selectedClubId === club.id}
								isJoining={
									joinMutation.isPending && joinMutation.variables === club.id
								}
								onSelect={() => {
									setSelectedClubId(club.id);
									setNotice("");
								}}
								onJoin={() => joinMutation.mutate(club.id)}
							/>
						))}
					</div>
				</div>

				<aside className="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
					<div className="mb-4 flex items-center gap-2">
						<Activity className="h-5 w-5 text-emerald-700" />
						<h2 className="font-semibold text-lg">Расписание</h2>
					</div>

					{notice ? (
						<div className="mb-4 rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2 text-emerald-800 text-sm">
							{notice}
						</div>
					) : null}

					{selectedClub ? (
						<div>
							<p className="font-medium text-slate-950">{selectedClub.name}</p>
							<p className="mt-1 text-slate-600 text-sm">
								{selectedClub.place} · {selectedClub.teacher_name}
							</p>
						</div>
					) : (
						<p className="text-slate-600 text-sm">
							Выберите секцию в каталоге, чтобы увидеть ближайшие тренировки.
						</p>
					)}

					<div className="mt-5 space-y-3">
						{workoutsQuery.isLoading ? (
							<p className="text-slate-500 text-sm">Загрузка расписания...</p>
						) : null}
						{workoutsQuery.data?.map((workout) => (
							<div
								key={workout.id}
								className="rounded-md border border-slate-200 bg-slate-50 p-3"
							>
								<p className="font-medium text-slate-900 text-sm">
									{formatDate(workout.start_date)}
								</p>
								<p className="mt-1 text-slate-600 text-sm">
									до {formatTime(workout.end_date)}
								</p>
							</div>
						))}
						{selectedClub && !workoutsQuery.isLoading && !workoutsQuery.data?.length ? (
							<p className="text-slate-600 text-sm">Тренировки еще не добавлены.</p>
						) : null}
					</div>
				</aside>
			</section>
		</main>
	);
}

function SignedOutView() {
	return (
		<main className="min-h-screen bg-[#f7f8f4] text-slate-950">
			<section className="mx-auto flex min-h-screen max-w-5xl flex-col justify-center px-4 py-10">
				<div className="max-w-2xl">
					<p className="font-medium text-emerald-700 text-sm">Sport Platform</p>
					<h1 className="mt-3 font-semibold text-4xl tracking-normal">
						Рабочий кабинет спортивных секций
					</h1>
					<p className="mt-4 text-lg text-slate-600">
						Войдите, чтобы открыть каталог секций, расписание тренировок и заявки
						на вступление.
					</p>
				</div>

				<div className="mt-8 flex flex-wrap gap-3">
					<Link
						href="/login"
						className="inline-flex items-center gap-2 rounded-md bg-emerald-700 px-4 py-2.5 font-medium text-white transition hover:bg-emerald-800"
					>
						<LogIn className="h-4 w-4" />
						Войти
					</Link>
					<Link
						href="/register"
						className="inline-flex items-center gap-2 rounded-md border border-slate-300 bg-white px-4 py-2.5 font-medium text-slate-900 transition hover:bg-slate-100"
					>
						<Users className="h-4 w-4" />
						Регистрация
					</Link>
				</div>

			</section>
		</main>
	);
}

function Metric({
	icon,
	label,
	value,
}: {
	icon: ReactNode;
	label: string;
	value: string;
}) {
	return (
		<div className="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
			<div className="flex items-center gap-3">
				<div className="flex h-10 w-10 items-center justify-center rounded-md bg-emerald-50 text-emerald-700">
					{icon}
				</div>
				<div>
					<p className="text-slate-500 text-sm">{label}</p>
					<p className="font-semibold text-lg text-slate-950">{value}</p>
				</div>
			</div>
		</div>
	);
}

function ClubCard({
	club,
	isSelected,
	isJoining,
	onSelect,
	onJoin,
}: {
	club: Club;
	isSelected: boolean;
	isJoining: boolean;
	onSelect: () => void;
	onJoin: () => void;
}) {
	return (
		<article
			className={`rounded-lg border bg-white p-5 shadow-sm transition ${
				isSelected ? "border-emerald-600 ring-2 ring-emerald-100" : "border-slate-200"
			}`}
		>
			<div className="flex items-start justify-between gap-3">
				<div>
					<p className="font-semibold text-lg text-slate-950">{club.name}</p>
					<p className="mt-1 text-emerald-700 text-sm">{club.sport_type_name}</p>
				</div>
				<span className="rounded-md bg-sky-50 px-2 py-1 text-sky-800 text-xs">
					{club.education_level_name}
				</span>
			</div>

			<p className="mt-3 line-clamp-3 text-slate-600 text-sm">
				{club.description}
			</p>

			<div className="mt-4 space-y-2 text-slate-700 text-sm">
				<p className="flex items-center gap-2">
					<MapPin className="h-4 w-4 text-slate-400" />
					{club.place}
				</p>
				<p className="flex items-center gap-2">
					<UserRound className="h-4 w-4 text-slate-400" />
					{club.teacher_name}
				</p>
				<p className="flex items-center gap-2">
					<CalendarDays className="h-4 w-4 text-slate-400" />
					{club.required_workout_per_week} тренировки в неделю · мест:{" "}
					{club.total_places ?? "без лимита"}
				</p>
			</div>

			<div className="mt-5 flex flex-wrap gap-2">
				<button
					type="button"
					onClick={onSelect}
					className="inline-flex items-center justify-center rounded-md border border-slate-300 bg-white px-3 py-2 font-medium text-slate-800 text-sm transition hover:bg-slate-100"
				>
					Расписание
				</button>
				<button
					type="button"
					onClick={onJoin}
					disabled={isJoining}
					className="inline-flex items-center justify-center rounded-md bg-emerald-700 px-3 py-2 font-medium text-sm text-white transition hover:bg-emerald-800 disabled:cursor-not-allowed disabled:bg-slate-400"
				>
					{isJoining ? "Отправка..." : "Подать заявку"}
				</button>
			</div>
		</article>
	);
}

function StateMessage({ title, text }: { title: string; text: string }) {
	return (
		<div className="mb-4 rounded-lg border border-slate-200 bg-white p-5 text-slate-700 shadow-sm">
			<p className="font-medium text-slate-950">{title}</p>
			<p className="mt-1 text-sm">{text}</p>
		</div>
	);
}

function formatDate(value: string) {
	return new Intl.DateTimeFormat("ru-RU", {
		day: "2-digit",
		month: "long",
		year: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	}).format(new Date(value));
}

function formatTime(value: string) {
	return new Intl.DateTimeFormat("ru-RU", {
		hour: "2-digit",
		minute: "2-digit",
	}).format(new Date(value));
}

export const dynamic = "force-dynamic";

export function GET() {
	const uptime = typeof process.uptime === "function" ? process.uptime() : 0;
	const body = [
		"# HELP sport_platform_frontend_up Frontend availability marker.",
		"# TYPE sport_platform_frontend_up gauge",
		"sport_platform_frontend_up 1",
		"# HELP sport_platform_frontend_uptime_seconds Frontend process uptime in seconds.",
		"# TYPE sport_platform_frontend_uptime_seconds gauge",
		`sport_platform_frontend_uptime_seconds ${uptime.toFixed(0)}`,
		"",
	].join("\n");

	return new Response(body, {
		headers: {
			"Content-Type": "text/plain; version=0.0.4; charset=utf-8",
		},
	});
}

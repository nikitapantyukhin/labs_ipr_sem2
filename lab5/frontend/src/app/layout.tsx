import "@/styles/globals.css";

import Providers from "@/app/providers";
import type { Metadata } from "next";
import { Geist } from "next/font/google";

export const metadata: Metadata = {
	title: "Sport Platform",
	description: "Платформа для спортивных секций, тренировок и заявок",
	icons: [{ rel: "icon", url: "/favicon.ico" }],
};

const geist = Geist({
	subsets: ["latin"],
	variable: "--font-geist-sans",
});

export default function RootLayout({
	children,
}: Readonly<{ children: React.ReactNode }>) {
	return (
		<html lang="ru" className={`${geist.variable}`}>
			<body>
				<Providers>{children}</Providers>
			</body>
		</html>
	);
}

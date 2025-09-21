"use client";

import { useEffect, useState } from "react";
import { getToken, isTokenValid, logout } from "@/lib/auth";
import { Button } from "@/components/atoms/Button";

export default function Home() {
	const [authenticated, setAuthenticated] = useState(false);
	useEffect(() => {
		const token = getToken();
		const ok = isTokenValid(token);
		setAuthenticated(ok);
		if (!ok) {
			window.location.href = "/login";
		}
	}, []);

	if (!authenticated) return null;

	return (
		<main
			style={{
				display: "grid",
				placeItems: "center",
				minHeight: "100vh",
				gap: 16,
			}}
		>
			Bienvenido
			<Button onClick={logout} variant="ghost" style={{ width: 160 }}>
				Cerrar sesión
			</Button>
		</main>
	);
}

export function getToken(): string | null {
	if (typeof window === "undefined") return null;
	return localStorage.getItem("auth_token");
}

export function parseJwt(token: string): { exp?: number } | null {
	try {
		const payload = token.split(".")[1];
		const json = JSON.parse(
			atob(payload.replace(/-/g, "+").replace(/_/g, "/"))
		);
		return json || null;
	} catch {
		return null;
	}
}

export function isTokenValid(token: string | null): boolean {
	if (!token) return false;
	const payload = parseJwt(token);
	if (!payload || !payload.exp) return false;
	const now = Math.floor(Date.now() / 1000);
	return payload.exp > now;
}

export function logout() {
	if (typeof window === "undefined") return;
	// Best-effort revoke on backend
	const base = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
	const token = getToken();
	if (token) {
		fetch(`${base}/logout`, {
			method: "POST",
			headers: { Authorization: `Bearer ${token}` },
		}).catch(() => {});
	}
	localStorage.removeItem("auth_token");
	window.location.href = "/login";
}

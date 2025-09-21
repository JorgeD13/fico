export const API_URL =
	process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function apiPost<T>(path: string, body: unknown): Promise<T> {
	const res = await fetch(`${API_URL}${path}`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(body),
	});
	if (!res.ok) throw new Error(await res.text());
	return (await res.json()) as T;
}

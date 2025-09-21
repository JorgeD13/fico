"use client";

import { useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@/components/atoms/Input";
import { Button } from "@/components/atoms/Button";
import { LabeledField } from "@/components/molecules/LabeledField";

export function LoginForm({
	onSuccess,
}: {
	onSuccess?: (token: string) => void;
}) {
	const router = useRouter();
	const emailRef = useRef<HTMLInputElement>(null);
	const passRef = useRef<HTMLInputElement>(null);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError(null);
		const email = emailRef.current?.value?.trim();
		const password = passRef.current?.value || "";
		if (!email || !password) {
			setError("Ingresa email y contraseña");
			return;
		}
		setLoading(true);
		try {
			const base =
				process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
			const res = await fetch(`${base}/login`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email, password }),
			});
			if (!res.ok) throw new Error("Credenciales inválidas");
			const data = (await res.json()) as { token: string };
			localStorage.setItem("auth_token", data.token);
			onSuccess?.(data.token);
			router.push("/");
		} catch (err: unknown) {
			setError(err instanceof Error ? err.message : "Error al iniciar sesión");
		} finally {
			setLoading(false);
		}
	};

	return (
		<form
			onSubmit={handleSubmit}
			style={{
				background: "#fff",
				padding: 16,
				borderRadius: 12,
				boxShadow: "0 2px 10px rgba(0,0,0,.06)",
			}}
		>
			<h1 style={{ fontSize: 20, fontWeight: 700, marginBottom: 16 }}>
				Iniciar sesión
			</h1>
			<LabeledField
				label="Email"
				input={
					<Input
						ref={emailRef}
						type="email"
						placeholder="demo@example.com"
					/>
				}
			/>
			<LabeledField
				label="Contraseña"
				input={
					<Input
						ref={passRef}
						type="password"
						placeholder="••••••••"
					/>
				}
			/>
			{error ? (
				<div
					style={{ marginBottom: 8, color: "#b91c1c", fontSize: 12 }}
				>
					{error}
				</div>
			) : null}
			<Button type="submit" disabled={loading}>
				{loading ? "Ingresando..." : "Ingresar"}
			</Button>
		</form>
	);
}

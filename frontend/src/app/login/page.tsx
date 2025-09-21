"use client";

import { useState } from "react";
import { LoginForm } from "@/components/organisms/LoginForm";

export default function LoginPage() {
	const [token, setToken] = useState<string | null>(null);

	return (
		<main
			style={{
				display: "grid",
				placeItems: "center",
				minHeight: "100vh",
				padding: 16,
			}}
		>
			<div style={{ width: "100%", maxWidth: 360 }}>
				<LoginForm onSuccess={setToken} />
				{token && (
					<pre
						style={{
							marginTop: 12,
							fontSize: 12,
							whiteSpace: "pre-wrap",
							wordBreak: "break-all",
						}}
					>
						{token}
					</pre>
				)}
			</div>
		</main>
	);
}

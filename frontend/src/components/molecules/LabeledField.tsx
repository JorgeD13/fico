"use client";

import { ReactNode } from "react";

export function LabeledField({
	label,
	input,
	error,
}: {
	label: string;
	input: ReactNode;
	error?: string;
}) {
	return (
		<label style={{ display: "block", marginBottom: 12 }}>
			<div style={{ marginBottom: 6, fontSize: 12, color: "#374151" }}>
				{label}
			</div>
			{input}
			{error ? (
				<div style={{ marginTop: 6, color: "#b91c1c", fontSize: 12 }}>
					{error}
				</div>
			) : null}
		</label>
	);
}

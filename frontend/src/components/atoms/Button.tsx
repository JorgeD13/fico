"use client";

import { ButtonHTMLAttributes } from "react";

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
	variant?: "primary" | "ghost";
};

export function Button({ variant = "primary", style, ...rest }: Props) {
	const base = {
		width: "100%",
		padding: "10px 12px",
		borderRadius: 8,
		cursor: "pointer",
		border: "1px solid transparent",
		fontWeight: 600,
	} as const;

	const variants = {
		primary: { background: "#111827", color: "#fff" },
		ghost: {
			background: "transparent",
			color: "#111827",
			border: "1px solid #111827",
		},
	} as const;

	return (
		<button
			{...rest}
			style={{
				...base,
				...(variants[variant] as any),
				...(style as any),
			}}
		/>
	);
}

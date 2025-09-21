"use client";

import { InputHTMLAttributes, forwardRef } from "react";

export const Input = forwardRef<
	HTMLInputElement,
	InputHTMLAttributes<HTMLInputElement>
>(function Input(props, ref) {
	return (
		<input
			ref={ref}
			{...props}
			style={{
				width: "100%",
				padding: "10px 12px",
				border: "1px solid #ccc",
				borderRadius: 8,
				outline: "none",
				...props.style,
			}}
		/>
	);
});

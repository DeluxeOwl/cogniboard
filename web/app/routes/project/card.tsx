import { cn } from "@/lib/utils";
import type * as React from "react";

type DivElemRef =
	| React.RefObject<HTMLDivElement>
	| ((element: HTMLElement | null) => void);

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
	ref?: DivElemRef;
}

function Card({ className, ref, ...props }: CardProps) {
	return (
		<div
			ref={ref}
			className={cn(
				"rounded-xl border bg-card text-card-foreground shadow",
				className,
			)}
			{...props}
		/>
	);
}
Card.displayName = "Card";

interface CardHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
	ref?: DivElemRef;
}

function CardHeader({ className, ref, ...props }: CardHeaderProps) {
	return (
		<div
			ref={ref}
			className={cn("flex flex-col space-y-1.5 p-6", className)}
			{...props}
		/>
	);
}
CardHeader.displayName = "CardHeader";

interface CardTitleProps extends React.HTMLAttributes<HTMLHeadingElement> {
	ref?:
		| React.RefObject<HTMLParagraphElement>
		| ((element: HTMLElement | null) => void);
}

function CardTitle({ className, ref, ...props }: CardTitleProps) {
	return (
		<h3
			ref={ref}
			className={cn("font-semibold leading-none tracking-tight", className)}
			{...props}
		/>
	);
}
CardTitle.displayName = "CardTitle";

interface CardDescriptionProps
	extends React.HTMLAttributes<HTMLParagraphElement> {
	ref?:
		| React.RefObject<HTMLParagraphElement>
		| ((element: HTMLElement | null) => void);
}

function CardDescription({ className, ref, ...props }: CardDescriptionProps) {
	return (
		<p
			ref={ref}
			className={cn("text-sm text-muted-foreground", className)}
			{...props}
		/>
	);
}
CardDescription.displayName = "CardDescription";

interface CardContentProps extends React.HTMLAttributes<HTMLDivElement> {
	ref?: DivElemRef;
}

function CardContent({ className, ref, ...props }: CardContentProps) {
	return <div ref={ref} className={cn("p-6 pt-0", className)} {...props} />;
}
CardContent.displayName = "CardContent";

interface CardFooterProps extends React.HTMLAttributes<HTMLDivElement> {
	ref?: DivElemRef;
}

function CardFooter({ className, ref, ...props }: CardFooterProps) {
	return (
		<div
			ref={ref}
			className={cn("flex items-center p-6 pt-0", className)}
			{...props}
		/>
	);
}
CardFooter.displayName = "CardFooter";

export {
	Card,
	CardHeader,
	CardFooter,
	CardTitle,
	CardDescription,
	CardContent,
};

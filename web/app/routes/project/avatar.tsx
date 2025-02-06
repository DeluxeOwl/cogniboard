import * as AvatarPrimitive from "@radix-ui/react-avatar";
import type * as React from "react";

import { cn } from "@/lib/utils";

function Avatar({
	className,
	ref,
	...props
}: React.ComponentProps<typeof AvatarPrimitive.Root>) {
	return (
		<AvatarPrimitive.Root
			ref={ref}
			className={cn(
				"relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full",
				className,
			)}
			{...props}
		/>
	);
}
Avatar.displayName = AvatarPrimitive.Root.displayName;

function AvatarImage({
	className,
	ref,
	...props
}: React.ComponentProps<typeof AvatarPrimitive.Image>) {
	return (
		<AvatarPrimitive.Image
			ref={ref}
			className={cn("aspect-square h-full w-full", className)}
			{...props}
		/>
	);
}
AvatarImage.displayName = AvatarPrimitive.Image.displayName;

function AvatarFallback({
	className,
	ref,
	...props
}: React.ComponentProps<typeof AvatarPrimitive.Fallback>) {
	return (
		<AvatarPrimitive.Fallback
			ref={ref}
			className={cn(
				"flex h-full w-full items-center justify-center rounded-full bg-muted",
				className,
			)}
			{...props}
		/>
	);
}
AvatarFallback.displayName = AvatarPrimitive.Fallback.displayName;

export { Avatar, AvatarImage, AvatarFallback };

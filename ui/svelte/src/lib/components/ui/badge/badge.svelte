<script>
	import { cn } from "$lib/utils";
	let { class: className = "", variant = "default", children, ...rest } = $props();

	const variants = {
		default: "border-transparent bg-primary text-primary-foreground",
		secondary: "border-transparent bg-secondary text-secondary-foreground",
		destructive: "border-transparent bg-destructive text-white",
		outline: "text-foreground"
	};
</script>

<span class={cn("focus:ring-ring/50 inline-flex w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-md border px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-colors focus:ring-[3px]", variants[variant], className)} {...rest}>
	{@render children?.()}
</span>

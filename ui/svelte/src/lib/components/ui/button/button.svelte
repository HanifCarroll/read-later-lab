<script context="module">
	import { cn } from "$lib/utils";

	const variants = {
		default: "bg-primary text-primary-foreground shadow-xs hover:bg-primary/90",
		destructive: "bg-destructive text-white shadow-xs hover:bg-destructive/90",
		outline: "border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground",
		secondary: "bg-secondary text-secondary-foreground shadow-xs hover:bg-secondary/80",
		ghost: "hover:bg-accent hover:text-accent-foreground",
		link: "text-primary underline-offset-4 hover:underline"
	};

	const sizes = {
		default: "h-9 px-4 py-2",
		sm: "h-8 rounded-md gap-1.5 px-3",
		lg: "h-10 rounded-md px-6",
		icon: "size-9"
	};

	export function buttonVariants({ variant = "default", size = "default", class: className = "" } = {}) {
		return cn(
			"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors disabled:pointer-events-none disabled:opacity-50 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]",
			variants[variant],
			sizes[size],
			className
		);
	}
</script>

<script>
	let { class: className = "", variant = "default", size = "default", type = "button", children, ...rest } = $props();
</script>

<button class={buttonVariants({ variant, size, class: className })} {type} {...rest}>
	{@render children?.()}
</button>

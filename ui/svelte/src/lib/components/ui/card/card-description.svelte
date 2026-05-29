<script>
	import { cn } from "$lib/utils";
	let { class: className = "", children, ...rest } = $props();
</script>

<p class={cn("text-muted-foreground text-sm", className)} {...rest}>
	{@render children?.()}
</p>

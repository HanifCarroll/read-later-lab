<script>
	import { cn } from "$lib/utils";
	let { class: className = "", children, ...rest } = $props();
</script>

<h3 class={cn("leading-none font-semibold", className)} {...rest}>
	{@render children?.()}
</h3>

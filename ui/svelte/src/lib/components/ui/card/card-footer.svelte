<script>
	import { cn } from "$lib/utils";
	let { class: className = "", children, ...rest } = $props();
</script>

<div class={cn("flex items-center px-6", className)} {...rest}>
	{@render children?.()}
</div>

<script>
	import { cn } from "$lib/utils";
	let { class: className = "", children, ...rest } = $props();
</script>

<div class={cn("px-6", className)} {...rest}>
	{@render children?.()}
</div>

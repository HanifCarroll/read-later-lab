<script>
	import { cn } from "$lib/utils";
	let { class: className = "", children, ...rest } = $props();
</script>

<div class={cn("grid auto-rows-min grid-rows-[auto_auto] items-start gap-1.5 px-6", className)} {...rest}>
	{@render children?.()}
</div>

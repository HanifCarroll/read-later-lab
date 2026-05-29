<script>
	import { cn } from "$lib/utils";
	let className = "";
	export { className as class };
</script>

<div class={cn("grid auto-rows-min grid-rows-[auto_auto] items-start gap-1.5 px-6", className)} {...$$restProps}>
	<slot />
</div>

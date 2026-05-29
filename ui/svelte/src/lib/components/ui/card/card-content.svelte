<script>
	import { cn } from "$lib/utils";
	let className = "";
	export { className as class };
</script>

<div class={cn("px-6", className)} {...$$restProps}>
	<slot />
</div>

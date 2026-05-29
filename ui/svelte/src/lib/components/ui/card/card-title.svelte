<script>
	import { cn } from "$lib/utils";
	let className = "";
	export { className as class };
</script>

<h3 class={cn("leading-none font-semibold", className)} {...$$restProps}>
	<slot />
</h3>

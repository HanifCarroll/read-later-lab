<script>
	import { cn } from "$lib/utils";
	let className = "";
	export { className as class };
</script>

<div class={cn("flex items-center px-6", className)} {...$$restProps}>
	<slot />
</div>

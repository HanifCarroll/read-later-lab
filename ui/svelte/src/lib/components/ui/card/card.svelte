<script>
	import { cn } from "$lib/utils";
	let className = "";
	export { className as class };
</script>

<div class={cn("bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm", className)} {...$$restProps}>
	<slot />
</div>

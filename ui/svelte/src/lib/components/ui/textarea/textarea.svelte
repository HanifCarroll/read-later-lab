<script>
	import { cn } from "$lib/utils";
	let className = "";
	export { className as class };
</script>

<textarea
	class={cn("border-input placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 flex min-h-16 w-full rounded-md border bg-transparent px-3 py-2 text-base shadow-xs transition-colors outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm", className)}
	{...$$restProps}
>
	<slot />
</textarea>

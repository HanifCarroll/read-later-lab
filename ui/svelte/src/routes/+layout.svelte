<script>
	import { onDestroy, onMount } from "svelte";
	import { QueryClient, setQueryClientContext } from "@tanstack/svelte-query";
	import "../styles.css";

	let { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 15_000,
			},
		},
	});

	setQueryClientContext(queryClient);
	onMount(() => queryClient.mount());
	onDestroy(() => queryClient.unmount());
</script>

{@render children()}

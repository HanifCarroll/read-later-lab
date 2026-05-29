<script>
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import "../styles.css";

	let { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 15_000,
			},
		},
	});
</script>

<QueryClientProvider client={queryClient}>
	{@render children()}
</QueryClientProvider>

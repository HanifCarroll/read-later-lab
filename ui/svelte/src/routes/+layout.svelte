<script>
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import "../styles.css";

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 15_000,
			},
		},
	});
</script>

<QueryClientProvider client={queryClient}>
	<slot />
</QueryClientProvider>

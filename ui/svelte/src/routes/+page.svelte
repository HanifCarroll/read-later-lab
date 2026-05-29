<script>
	import { createMutation, createQuery, useQueryClient } from "@tanstack/svelte-query";
	import { Badge } from "$lib/components/ui/badge";
	import { Button } from "$lib/components/ui/button";
	import * as Card from "$lib/components/ui/card";
	import { Input } from "$lib/components/ui/input";
	import { Textarea } from "$lib/components/ui/textarea";
	import { createItem, listItems, updateStatus } from "$lib/api.js";

	const queryClient = useQueryClient();
	const itemsQuery = createQuery({
		queryKey: ["items"],
		queryFn: listItems,
	});
	const createItemMutation = createMutation({
		mutationFn: createItem,
		onSuccess: () => {
			url = "";
			savedReason = "";
			queryClient.invalidateQueries({ queryKey: ["items"] });
		},
	});
	const updateStatusMutation = createMutation({
		mutationFn: ({ id, status }) => updateStatus(id, status),
		onSuccess: () => queryClient.invalidateQueries({ queryKey: ["items"] }),
	});

	let url = "";
	let savedReason = "";

	function save() {
		createItemMutation.mutate({ url, savedReason });
	}

	function setStatus(item, status) {
		updateStatusMutation.mutate({ id: item.id, status });
	}

	function modeLabel(mode) {
		return String(mode || "skim").replaceAll("_", " ");
	}
</script>

<svelte:head>
	<title>Read Later Lab</title>
</svelte:head>

<div class="min-h-screen bg-background">
	<header class="border-b bg-card">
		<div class="mx-auto flex max-w-6xl flex-col gap-2 px-6 py-8 sm:flex-row sm:items-center sm:justify-between">
			<div>
				<p class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">SvelteKit SPA + Go API</p>
				<h1 class="mt-2 text-3xl font-bold tracking-tight">Read Later Lab</h1>
			</div>
			<Badge variant="secondary">Postgres + OpenAI</Badge>
		</div>
	</header>

	<main class="mx-auto grid max-w-6xl gap-8 px-6 py-8">
		<Card.Card>
			<Card.CardHeader>
				<Card.CardTitle>Save an article</Card.CardTitle>
				<Card.CardDescription>Capture why you saved it, then generate a quick read/skip triage card.</Card.CardDescription>
			</Card.CardHeader>
			<Card.CardContent>
				<form class="grid gap-4" on:submit|preventDefault={save}>
					<label class="grid gap-2 text-sm font-medium">
						URL
						<Input required type="url" bind:value={url} placeholder="https://example.com/post" />
					</label>
					<label class="grid gap-2 text-sm font-medium">
						Why did you save this?
						<Textarea required bind:value={savedReason} placeholder="I saved this because..." class="min-h-24" />
					</label>
					<div class="flex items-center gap-3">
						<Button type="submit" disabled={$createItemMutation.isPending}>{$createItemMutation.isPending ? "Analyzing…" : "Analyze and save"}</Button>
						{#if $createItemMutation.error}<p class="text-sm font-medium text-destructive">{$createItemMutation.error.message}</p>{/if}
					</div>
				</form>
			</Card.CardContent>
		</Card.Card>

		<section class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
			<div>
				<h2 class="text-2xl font-semibold tracking-tight">Inbox triage</h2>
				<p class="text-sm text-muted-foreground">{$itemsQuery.isLoading ? "Loading…" : `${$itemsQuery.data?.length ?? 0} saved articles`}</p>
			</div>
			<Button variant="outline" on:click={() => queryClient.invalidateQueries({ queryKey: ["items"] })}>Refresh</Button>
		</section>

		{#if $itemsQuery.error}
			<p class="rounded-lg border border-destructive/30 bg-destructive/10 p-4 text-sm font-medium text-destructive">{$itemsQuery.error.message}</p>
		{/if}

		<section class="grid gap-4">
			{#each $itemsQuery.data ?? [] as item}
				<Card.Card>
					<Card.CardHeader>
						<div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
							<div class="grid gap-2">
								<div class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
									<span>{item.site}</span>
									<span>·</span>
									<span>{modeLabel(item.recommendedMode)}</span>
									<span>·</span>
									<span>{item.readTimeMinutes} min</span>
								</div>
								<Card.CardTitle><a class="hover:underline" href={item.url} target="_blank" rel="noreferrer">{item.title}</a></Card.CardTitle>
							</div>
							<Badge variant="outline">{modeLabel(item.status)}</Badge>
						</div>
					</Card.CardHeader>
					<Card.CardContent class="grid gap-4">
						<p class="text-base leading-7">{item.summary}</p>
						<p class="rounded-lg bg-muted px-4 py-3 text-sm"><strong>Saved because:</strong> {item.savedReason}</p>
						<details class="rounded-lg border p-4">
							<summary class="cursor-pointer font-medium">Read/skip guidance</summary>
							<div class="mt-4 grid gap-4 lg:grid-cols-3">
								<div>
									<h4 class="font-semibold">Best sections</h4>
									<ul class="mt-2 list-disc space-y-1 pl-5 text-sm text-muted-foreground">{#each item.bestSections as section}<li>{section}</li>{/each}</ul>
								</div>
								<div>
									<h4 class="font-semibold">Read this if</h4>
									<ul class="mt-2 list-disc space-y-1 pl-5 text-sm text-muted-foreground">{#each item.readIf as reason}<li>{reason}</li>{/each}</ul>
								</div>
								<div>
									<h4 class="font-semibold">Skip this if</h4>
									<ul class="mt-2 list-disc space-y-1 pl-5 text-sm text-muted-foreground">{#each item.skipIf as reason}<li>{reason}</li>{/each}</ul>
								</div>
							</div>
						</details>
					</Card.CardContent>
					<Card.CardFooter class="flex flex-wrap gap-2">
						<Button size="sm" variant="secondary" disabled={$updateStatusMutation.isPending} on:click={() => setStatus(item, "read_soon")}>Read soon</Button>
						<Button size="sm" variant="secondary" disabled={$updateStatusMutation.isPending} on:click={() => setStatus(item, "skim_later")}>Skim later</Button>
						<Button size="sm" variant="secondary" disabled={$updateStatusMutation.isPending} on:click={() => setStatus(item, "reference")}>Reference</Button>
						<Button size="sm" variant="outline" disabled={$updateStatusMutation.isPending} on:click={() => setStatus(item, "skipped")}>Skip</Button>
						<Button size="sm" variant="outline" disabled={$updateStatusMutation.isPending} on:click={() => setStatus(item, "archived")}>Archive</Button>
					</Card.CardFooter>
				</Card.Card>
			{/each}
		</section>
	</main>
</div>

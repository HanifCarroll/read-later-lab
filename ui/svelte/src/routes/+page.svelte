<script>
  import { onMount } from 'svelte';
  import { createItem, listItems, updateStatus } from '$lib/api.js';

  let items = [];
  let url = '';
  let savedReason = '';
  let loading = true;
  let saving = false;
  let error = '';

  onMount(async () => {
    await refresh();
  });

  async function refresh() {
    loading = true;
    error = '';
    try {
      items = await listItems();
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  async function save() {
    saving = true;
    error = '';
    try {
      const item = await createItem({ url, savedReason });
      items = [item, ...items];
      url = '';
      savedReason = '';
    } catch (err) {
      error = err.message;
    } finally {
      saving = false;
    }
  }

  async function setStatus(item, status) {
    const updated = await updateStatus(item.id, status);
    items = items.map((current) => (current.id === item.id ? updated : current));
  }
</script>

<header class="topbar">
  <div>
    <p class="eyebrow">SvelteKit SPA + Go JSON API</p>
    <h1>Read Later Lab</h1>
  </div>
  <a href="/h">Try HTMX version</a>
</header>

<main>
  <section class="panel">
    <h2>Save an article</h2>
    <form class="save-form" on:submit|preventDefault={save}>
      <label>URL <input required type="url" bind:value={url} placeholder="https://example.com/post" /></label>
      <label>Why did you save this? <textarea required bind:value={savedReason} placeholder="I saved this because..."></textarea></label>
      <button disabled={saving}>{saving ? 'Analyzing…' : 'Analyze and save'}</button>
    </form>
    {#if error}<p class="error">{error}</p>{/if}
  </section>

  <section class="list-head">
    <h2>Inbox triage</h2>
    <p>{loading ? 'Loading…' : `${items.length} saved articles`}</p>
  </section>

  <section class="grid">
    {#each items as item}
      <article class="card">
        <div class="card-top">
          <div>
            <p class="eyebrow">{item.site} · {item.recommendedMode} · {item.readTimeMinutes} min</p>
            <h3><a href={item.url} target="_blank" rel="noreferrer">{item.title}</a></h3>
          </div>
          <span class={`status ${item.status}`}>{item.status}</span>
        </div>
        <p class="summary">{item.summary}</p>
        <p><strong>Saved because:</strong> {item.savedReason}</p>
        <details>
          <summary>Read/skip guidance</summary>
          <h4>Best sections to read</h4>
          <ul>{#each item.bestSections as section}<li>{section}</li>{/each}</ul>
          <div class="columns">
            <div>
              <h4>You should read this if</h4>
              <ul>{#each item.readIf as reason}<li>{reason}</li>{/each}</ul>
            </div>
            <div>
              <h4>You should skip this if</h4>
              <ul>{#each item.skipIf as reason}<li>{reason}</li>{/each}</ul>
            </div>
          </div>
        </details>
        <div class="actions">
          <button on:click={() => setStatus(item, 'read_soon')}>Read soon</button>
          <button on:click={() => setStatus(item, 'skim_later')}>Skim later</button>
          <button on:click={() => setStatus(item, 'reference')}>Reference</button>
          <button on:click={() => setStatus(item, 'skipped')}>Skip</button>
          <button on:click={() => setStatus(item, 'archived')}>Archive</button>
        </div>
      </article>
    {/each}
  </section>
</main>

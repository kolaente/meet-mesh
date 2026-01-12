<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, LinkCard } from '$lib/components/dashboard';
	import { Button, Spinner } from '$lib/components/ui';

	type Link = components['schemas']['Link'];

	let links = $state<Link[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const { data } = await api.GET('/links');
			if (data) {
				links = data;
			}
		} finally {
			loading = false;
		}
	});
</script>

<DashboardHeader title="Links">
	{#snippet actions()}
		<Button variant="primary" onclick={() => (window.location.href = '/links/new')}>
			{#snippet children()}New Link{/snippet}
		</Button>
	{/snippet}
</DashboardHeader>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if links.length === 0}
	<div class="text-center py-12">
		<svg
			class="mx-auto h-12 w-12 text-gray-400"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
			/>
		</svg>
		<h3 class="mt-2 text-sm font-medium text-gray-900">No links</h3>
		<p class="mt-1 text-sm text-gray-500">Get started by creating a new scheduling link.</p>
		<div class="mt-6">
			<Button variant="primary" onclick={() => (window.location.href = '/links/new')}>
				{#snippet children()}Create your first link{/snippet}
			</Button>
		</div>
	</div>
{:else}
	<!-- 1 column mobile, 2 columns tablet (sm:), 3 columns desktop (lg:) -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
		{#each links as link (link.id)}
			<LinkCard {link} />
		{/each}
	</div>
{/if}

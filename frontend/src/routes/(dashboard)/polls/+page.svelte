<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader } from '$lib/components/dashboard';
	import { Button, Spinner, Card, Badge } from '$lib/components/ui';

	type Poll = components['schemas']['Poll'];

	let polls = $state<Poll[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const { data } = await api.GET('/polls');
			if (data) {
				polls = data;
			}
		} finally {
			loading = false;
		}
	});
</script>

<DashboardHeader title="Polls">
	{#snippet actions()}
		<Button variant="primary" onclick={() => goto('/polls/new')}>
			{#snippet children()}New Poll{/snippet}
		</Button>
	{/snippet}
</DashboardHeader>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if polls.length === 0}
	<div class="text-center py-12">
		<svg
			class="mx-auto h-12 w-12 text-text-muted"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
			/>
		</svg>
		<h3 class="mt-2 text-sm font-medium text-text-primary">No polls</h3>
		<p class="mt-1 text-sm text-text-secondary">Get started by creating a poll for group scheduling.</p>
		<div class="mt-6">
			<Button variant="primary" onclick={() => goto('/polls/new')}>
				{#snippet children()}Create your first poll{/snippet}
			</Button>
		</div>
	</div>
{:else}
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
		{#each polls as poll (poll.id)}
			<a href="/polls/{poll.id}" class="block">
				<Card class="hover:shadow-md transition-shadow cursor-pointer active:shadow-sm">
					{#snippet children()}
						<div class="space-y-2 sm:space-y-3">
							<div class="flex items-start justify-between gap-2">
								<h3 class="font-medium text-text-primary truncate text-sm sm:text-base">{poll.name}</h3>
								<Badge variant={poll.status === 1 ? 'active' : 'cancelled'} size="sm" />
							</div>

							<div class="text-xs sm:text-sm text-text-secondary truncate">
								<span class="font-mono">/p/poll/{poll.slug}</span>
							</div>
						</div>
					{/snippet}
				</Card>
			</a>
		{/each}
	</div>
{/if}

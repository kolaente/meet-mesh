<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, VoteRow, PollOptionManager } from '$lib/components/dashboard';
	import { Button, Card, Badge, Spinner } from '$lib/components/ui';

	type Poll = components['schemas']['Poll'];
	type PollOption = components['schemas']['PollOption'];
	type Vote = components['schemas']['Vote'];

	let poll = $state<Poll | null>(null);
	let options = $state<PollOption[]>([]);
	let votes = $state<Vote[]>([]);
	let loading = $state(true);
	let deleting = $state(false);
	let copySuccess = $state(false);

	const pollId = $derived(Number($page.params.id));

	const statusVariant = $derived<'active' | 'cancelled'>(poll?.status === 1 ? 'active' : 'cancelled');

	const publicUrl = $derived(poll ? `${window.location.origin}/p/poll/${poll.slug}` : '');

	onMount(async () => {
		try {
			const { data: pollData } = await api.GET('/polls/{id}', {
				params: { path: { id: pollId } }
			});

			if (!pollData) {
				goto('/polls');
				return;
			}

			poll = pollData;

			// Fetch options
			const { data: optionsData } = await api.GET('/polls/{id}/options', {
				params: { path: { id: pollId } }
			});
			if (optionsData) {
				options = optionsData;
			}

			// Fetch votes
			const { data: votesData } = await api.GET('/polls/{id}/votes', {
				params: { path: { id: pollId } }
			});
			if (votesData) {
				votes = votesData;
			}
		} finally {
			loading = false;
		}
	});

	async function copyLink() {
		await navigator.clipboard.writeText(publicUrl);
		copySuccess = true;
		setTimeout(() => {
			copySuccess = false;
		}, 2000);
	}

	async function handleDelete() {
		if (!confirm('Are you sure you want to delete this poll? This action cannot be undone.')) {
			return;
		}

		deleting = true;
		try {
			await api.DELETE('/polls/{id}', {
				params: { path: { id: pollId } }
			});
			goto('/polls');
		} finally {
			deleting = false;
		}
	}
</script>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if poll}
	<DashboardHeader title={poll.name}>
		{#snippet actions()}
			<Button variant="secondary" onclick={copyLink}>
				{#snippet children()}
					{#if copySuccess}
						Copied!
					{:else}
						Copy Link
					{/if}
				{/snippet}
			</Button>
			<Button variant="secondary" onclick={() => goto(`/polls/${pollId}/edit`)}>
				{#snippet children()}Edit{/snippet}
			</Button>
			<Button variant="danger" onclick={handleDelete} loading={deleting}>
				{#snippet children()}Delete{/snippet}
			</Button>
		{/snippet}
	</DashboardHeader>

	<div class="space-y-6">
		<!-- Poll Info -->
		<Card>
			{#snippet children()}
				<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
					<div>
						<p class="text-sm font-medium text-gray-500">Type</p>
						<span class="inline-flex items-center px-2.5 py-1 text-sm font-medium rounded-full bg-indigo-50 text-indigo-700 border border-indigo-200">
							Poll
						</span>
					</div>

					<div>
						<p class="text-sm font-medium text-gray-500">Status</p>
						<Badge variant={statusVariant} size="md" />
					</div>

					<div class="col-span-2">
						<p class="text-sm font-medium text-gray-500">Public URL</p>
						<a
							href={publicUrl}
							target="_blank"
							rel="noopener noreferrer"
							class="text-indigo-600 hover:text-indigo-700 font-mono text-sm"
						>
							{publicUrl}
						</a>
					</div>

					{#if poll?.description}
						<div class="col-span-2 md:col-span-4">
							<p class="text-sm font-medium text-gray-500">Description</p>
							<p class="text-gray-700">{poll?.description}</p>
						</div>
					{/if}
				</div>
			{/snippet}
		</Card>

		<!-- Poll Options -->
		<PollOptionManager {pollId} bind:options />

		<!-- Votes -->
		<Card>
			{#snippet header()}
				<h2 class="text-lg font-medium text-gray-900">
					Votes ({votes.length})
				</h2>
			{/snippet}

			{#snippet children()}
				{#if votes.length === 0}
					<p class="text-gray-500 text-center py-4">No votes yet</p>
				{:else}
					<div class="divide-y divide-gray-100">
						{#each votes as vote (vote.id)}
							<VoteRow {vote} slots={options} />
						{/each}
					</div>
				{/if}
			{/snippet}
		</Card>
	</div>
{/if}

<script lang="ts">
	import { page } from '$app/state';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { Card, Spinner } from '$lib/components/ui';

	type VoteTally = components['schemas']['VoteTally'];
	type Vote = components['schemas']['Vote'];

	let tally = $state<VoteTally[]>([]);
	let votes = $state<Vote[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const slug = $derived(page.params.slug ?? '');

	$effect(() => {
		if (slug) {
			loadResults();
		}
	});

	async function loadResults() {
		if (!slug) return;

		loading = true;
		error = null;

		const { data, error: apiError } = await api.GET('/p/poll/{slug}/results', {
			params: { path: { slug } }
		});

		if (apiError) {
			error = 'Results are not available for this poll.';
			loading = false;
			return;
		}

		if (data) {
			tally = data.tally || [];
			votes = data.votes || [];
		}
		loading = false;
	}

	function getTotalVotes(t: VoteTally): number {
		return t.yes_count + t.no_count + t.maybe_count;
	}
</script>

<div class="max-w-2xl mx-auto">
	<h1 class="text-2xl font-bold text-gray-900 mb-8 text-center">Poll Results</h1>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<Spinner size="lg" />
		</div>
	{:else if error}
		<Card class="text-center py-8">
			<div class="text-red-500 mb-4">
				<svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
				</svg>
			</div>
			<p class="text-gray-600">{error}</p>
		</Card>
	{:else}
		<Card>
			{#snippet header()}
				<h2 class="text-lg font-medium text-gray-900">Vote Summary</h2>
			{/snippet}

			{#snippet children()}
				{#if tally.length === 0}
					<p class="text-gray-500 text-center py-4">No votes yet</p>
				{:else}
					<div class="space-y-4">
						{#each tally as t (t.option_id)}
							{@const total = getTotalVotes(t)}
							<div class="p-4 bg-gray-50 rounded-lg">
								<div class="flex items-center justify-between mb-2">
									<span class="font-medium text-gray-900">Option {t.option_id}</span>
									<span class="text-sm text-gray-500">{total} votes</span>
								</div>
								<div class="flex gap-2">
									<div class="flex-1">
										<div class="flex items-center justify-between text-sm mb-1">
											<span class="text-green-700">Yes</span>
											<span class="text-green-700">{t.yes_count}</span>
										</div>
										<div class="h-2 bg-gray-200 rounded-full overflow-hidden">
											<div
												class="h-full bg-green-500 rounded-full"
												style="width: {total > 0 ? (t.yes_count / total) * 100 : 0}%"
											></div>
										</div>
									</div>
									<div class="flex-1">
										<div class="flex items-center justify-between text-sm mb-1">
											<span class="text-amber-700">Maybe</span>
											<span class="text-amber-700">{t.maybe_count}</span>
										</div>
										<div class="h-2 bg-gray-200 rounded-full overflow-hidden">
											<div
												class="h-full bg-amber-500 rounded-full"
												style="width: {total > 0 ? (t.maybe_count / total) * 100 : 0}%"
											></div>
										</div>
									</div>
									<div class="flex-1">
										<div class="flex items-center justify-between text-sm mb-1">
											<span class="text-red-700">No</span>
											<span class="text-red-700">{t.no_count}</span>
										</div>
										<div class="h-2 bg-gray-200 rounded-full overflow-hidden">
											<div
												class="h-full bg-red-500 rounded-full"
												style="width: {total > 0 ? (t.no_count / total) * 100 : 0}%"
											></div>
										</div>
									</div>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			{/snippet}
		</Card>
	{/if}
</div>

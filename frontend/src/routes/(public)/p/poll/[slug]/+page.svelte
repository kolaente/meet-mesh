<script lang="ts">
	import { page } from '$app/state';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { Spinner, Card } from '$lib/components/ui';
	import { PollPage } from '$lib/components/poll';

	type PollOption = components['schemas']['PollOption'];
	type CustomField = components['schemas']['CustomField'];

	interface PublicPoll {
		name: string;
		description?: string;
		custom_fields?: CustomField[];
		options: PollOption[];
		show_results?: boolean;
		require_email?: boolean;
		organizer_name?: string;
		organizer_avatar_url?: string;
	}

	let poll = $state<PublicPoll | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const slug = $derived(page.params.slug ?? '');

	$effect(() => {
		if (slug) {
			loadPoll();
		}
	});

	async function loadPoll() {
		if (!slug) return;

		loading = true;
		error = null;

		const { data, error: apiError } = await api.GET('/p/poll/{slug}', {
			params: { path: { slug } }
		});

		if (apiError) {
			error = 'This poll is not available.';
			loading = false;
			return;
		}

		poll = data ?? null;
		loading = false;
	}
</script>

<svelte:head>
	<title>{poll?.name || 'Vote'} | Meet Mesh</title>
</svelte:head>

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
		<h1 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Page Not Found</h1>
		<p class="text-gray-600 dark:text-gray-400">{error}</p>
	</Card>
{:else if poll}
	<PollPage
		link={{
			type: 2,
			name: poll.name,
			description: poll.description,
			custom_fields: poll.custom_fields,
			slots: poll.options,
			show_results: poll.show_results,
			require_email: poll.require_email,
			organizer_name: poll.organizer_name,
			organizer_avatar_url: poll.organizer_avatar_url
		}}
		{slug}
	/>
{/if}

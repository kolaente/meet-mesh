<script lang="ts">
	import { fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import Card from '../ui/Card.svelte';
	import VoteCard from './VoteCard.svelte';
	import VoterForm from './VoterForm.svelte';

	type Slot = components['schemas']['Slot'];
	type CustomField = components['schemas']['CustomField'];
	type VoteResponse = components['schemas']['VoteResponse'];

	interface PublicLink {
		type: 1 | 2;
		name: string;
		description?: string;
		custom_fields?: CustomField[];
		slots: Slot[];
		show_results?: boolean;
		require_email?: boolean;
	}

	interface Props {
		link: PublicLink;
		slug: string;
	}

	let { link, slug }: Props = $props();

	// State for tracking votes per slot (slotId -> VoteResponse)
	let votes = $state<Record<number, VoteResponse | undefined>>({});

	// UI state
	let submitting = $state(false);
	let error = $state<string | undefined>();
	let submitted = $state(false);

	// Count how many votes the user has made
	let voteCount = $derived(
		Object.values(votes).filter((v) => v !== undefined).length
	);

	// Check if user has voted on all slots (not required, but nice feedback)
	let allVoted = $derived(
		link.slots.every((slot) => votes[slot.id] !== undefined)
	);

	// Handle vote change from VoteCard
	function handleVote(slotId: number, vote: VoteResponse | undefined) {
		votes[slotId] = vote;
	}

	// Handle form submission
	async function handleSubmit(data: { name?: string; email?: string }) {
		// Validate that at least one vote is cast
		if (voteCount === 0) {
			error = 'Please cast at least one vote before submitting.';
			return;
		}

		submitting = true;
		error = undefined;

		try {
			// Build responses object (slotId as string key -> VoteResponse)
			const responses: Record<string, VoteResponse> = {};
			for (const [slotId, vote] of Object.entries(votes)) {
				if (vote !== undefined) {
					responses[slotId] = vote;
				}
			}

			const { data: responseData, error: responseError } = await api.POST('/p/poll/{slug}/vote', {
				params: { path: { slug } },
				body: {
					guest_name: data.name,
					guest_email: data.email,
					responses
				}
			});

			if (responseError) {
				error = (responseError as { message?: string }).message || 'Failed to submit vote';
				submitting = false;
				return;
			}

			// Success - show confirmation or redirect to results
			submitted = true;

			// If show_results is enabled, redirect to results page after a moment
			if (link.show_results) {
				setTimeout(() => {
					goto(`/p/poll/${slug}/results`);
				}, 2000);
			}
		} catch (err) {
			error = 'An unexpected error occurred. Please try again.';
			submitting = false;
		}
	}
</script>

<div class="max-w-2xl mx-auto">
	<!-- Header -->
	<div class="text-center mb-8">
		<h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{link.name}</h1>
		{#if link.description}
			<p class="mt-2 text-gray-600 dark:text-gray-400">{link.description}</p>
		{/if}
	</div>

	<!-- Error message -->
	{#if error}
		<div
			class="mb-6 p-4 bg-red-50 border border-red-200 rounded-brutalist-md text-red-700"
			transition:fade
		>
			{error}
		</div>
	{/if}

	{#if submitted}
		<!-- Success state -->
		<Card>
			<div class="text-center py-8">
				<div class="w-16 h-16 mx-auto mb-4 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center">
					<svg class="w-8 h-8 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
					</svg>
				</div>
				<h2 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Vote Submitted!</h2>
				<p class="text-gray-600 dark:text-gray-400">Thank you for participating in this poll.</p>
				{#if link.show_results}
					<p class="text-sm text-gray-500 dark:text-gray-400 mt-4">Redirecting to results...</p>
				{/if}
			</div>
		</Card>
	{:else if link.slots.length === 0}
		<!-- No slots available -->
		<Card>
			<div class="text-center py-8">
				<p class="text-gray-600 dark:text-gray-400">No options available for this poll.</p>
				<p class="text-sm text-gray-500 dark:text-gray-400 mt-2">Please check back later.</p>
			</div>
		</Card>
	{:else}
		<!-- Voting interface -->
		<div class="space-y-4">
			<!-- Instructions -->
			<p class="text-sm text-gray-600 dark:text-gray-400">
				Select your availability for each option below. You can choose Yes, No, or Maybe.
			</p>

			<!-- Vote progress -->
			<div class="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
				<span>
					{voteCount} of {link.slots.length} options voted
				</span>
				{#if allVoted}
					<span class="text-green-600 dark:text-green-400 font-medium">All options voted</span>
				{/if}
			</div>

			<!-- Poll cards container with consistent spacing -->
			<div class="flex flex-col gap-4">
				<!-- Vote cards -->
				{#each link.slots as slot (slot.id)}
					<VoteCard
						option={slot}
						bind:vote={votes[slot.id]}
						onVote={handleVote}
					/>
				{/each}

				<!-- Voter form -->
				<Card>
					<VoterForm
						onSubmit={handleSubmit}
						loading={submitting}
						requireEmail={link.require_email}
					/>
				</Card>
			</div>
		</div>
	{/if}
</div>

<script lang="ts">
	import type { components } from '$lib/api/types';
	import Card from '../ui/Card.svelte';
	import VoteButtons from './VoteButtons.svelte';

	type Slot = components['schemas']['Slot'];
	type VoteResponse = components['schemas']['VoteResponse'];

	interface Props {
		option: Slot;
		vote?: VoteResponse;
		onVote?: (slotId: number, vote: VoteResponse | undefined) => void;
		class?: string;
	}

	let {
		option,
		vote = $bindable(),
		onVote,
		class: className = ''
	}: Props = $props();

	// Format date nicely
	let formattedDate = $derived.by(() => {
		const date = new Date(option.start_time);
		return date.toLocaleDateString('en-US', {
			weekday: 'long',
			month: 'long',
			day: 'numeric'
		});
	});

	// Format time range
	let formattedTime = $derived.by(() => {
		const start = new Date(option.start_time);
		const end = new Date(option.end_time);

		// SlotType: 1=time, 2=full_day, 3=multi_day
		if (option.type === 2) {
			// Full day - no time display
			return 'All day';
		} else if (option.type === 3) {
			// Multi-day - show date range
			const endDate = end.toLocaleDateString('en-US', {
				weekday: 'long',
				month: 'long',
				day: 'numeric'
			});
			return `to ${endDate}`;
		} else {
			// Time slot - show time range
			const startTime = start.toLocaleTimeString('en-US', {
				hour: 'numeric',
				minute: '2-digit'
			});
			const endTime = end.toLocaleTimeString('en-US', {
				hour: 'numeric',
				minute: '2-digit'
			});
			return `${startTime} - ${endTime}`;
		}
	});

	function handleVoteChange(newVote: VoteResponse | undefined) {
		vote = newVote;
		onVote?.(option.id, newVote);
	}
</script>

<Card class={className}>
	<div class="flex flex-col gap-3 sm:gap-4">
		<!-- Date/Time info -->
		<div class="flex items-start gap-3">
			<div class="flex-shrink-0 text-2xl" aria-hidden="true">
				{#if option.type === 3}
					<!-- Multi-day icon -->
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
					</svg>
				{:else}
					<!-- Calendar icon -->
					<svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
					</svg>
				{/if}
			</div>
			<div class="flex-1 min-w-0">
				<p class="text-sm sm:text-base font-medium text-gray-900">{formattedDate}</p>
				<p class="text-xs sm:text-sm text-gray-500">{formattedTime}</p>
			</div>
		</div>

		<!-- Vote buttons - full width on mobile -->
		<div class="flex justify-start w-full">
			<VoteButtons bind:value={vote} onChange={handleVoteChange} />
		</div>
	</div>
</Card>

<script lang="ts">
	import type { components } from '$lib/api/types';
	import { formatTime } from '$lib/utils/dates';

	type Vote = components['schemas']['Vote'];
	type Slot = components['schemas']['Slot'];
	type PollOption = components['schemas']['PollOption'];
	type VoteResponse = components['schemas']['VoteResponse'];

	interface Props {
		vote: Vote;
		slots: (Slot | PollOption)[];
	}

	let { vote, slots }: Props = $props();

	const responseLabels: Record<VoteResponse, string> = {
		1: 'Yes',
		2: 'No',
		3: 'Maybe'
	};

	const responseColors: Record<VoteResponse, string> = {
		1: 'bg-green-100 text-green-700',
		2: 'bg-red-100 text-red-700',
		3: 'bg-amber-100 text-amber-700'
	};

	const formatSlotLabel = (slot: Slot) => {
		const start = new Date(slot.start_time);
		const dateStr = start.toLocaleDateString(undefined, {
			month: 'short',
			day: 'numeric'
		});
		return `${dateStr}, ${formatTime(start)}`;
	};
</script>

<div class="flex items-center py-3 border-b border-gray-100 dark:border-neutral-700 last:border-b-0">
	<div class="w-40 flex-shrink-0">
		<p class="font-medium text-gray-900 dark:text-gray-100 truncate">
			{vote.guest_name || 'Anonymous'}
		</p>
		{#if vote.guest_email}
			<p class="text-sm text-gray-500 dark:text-gray-400 truncate">{vote.guest_email}</p>
		{/if}
	</div>

	<div class="flex gap-2 flex-wrap flex-1">
		{#each slots as slot (slot.id)}
			{@const response = vote.responses[String(slot.id)]}
			{#if response}
				<div class="flex items-center gap-1 text-xs">
					<span class="text-gray-500 dark:text-gray-400">{formatSlotLabel(slot)}:</span>
					<span class="px-1.5 py-0.5 rounded {responseColors[response]}">
						{responseLabels[response]}
					</span>
				</div>
			{/if}
		{/each}
	</div>
</div>

<script lang="ts">
	import type { components } from '$lib/api/types';

	type Vote = components['schemas']['Vote'];
	type VoteResponse = components['schemas']['VoteResponse'];

	interface Props {
		votes: Vote[];
		slotId: number;
		class?: string;
	}

	let { votes, slotId, class: className = '' }: Props = $props();

	// Count votes for this slot
	let counts = $derived.by(() => {
		let yes = 0;
		let no = 0;
		let maybe = 0;

		for (const vote of votes) {
			const response = vote.responses[String(slotId)] as VoteResponse | undefined;
			if (response === 1) yes++;
			else if (response === 2) no++;
			else if (response === 3) maybe++;
		}

		return { yes, no, maybe };
	});

	let total = $derived(counts.yes + counts.no + counts.maybe);
	let maxCount = $derived(Math.max(counts.yes, counts.no, counts.maybe, 1));

	// Calculate bar widths as percentage of max
	function getBarWidth(count: number): string {
		if (maxCount === 0) return '0%';
		return `${(count / maxCount) * 100}%`;
	}
</script>

<div class="space-y-2 {className}">
	<!-- Yes bar -->
	<div class="flex items-center gap-2">
		<span class="w-12 text-xs font-medium text-gray-600">Yes</span>
		<div class="flex-1 h-5 bg-gray-100 rounded-[var(--radius-sm)] overflow-hidden">
			<div
				class="h-full bg-indigo-500 rounded-[var(--radius-sm)] transition-all duration-300"
				style="width: {getBarWidth(counts.yes)}"
			></div>
		</div>
		<span class="w-8 text-xs font-medium text-gray-600 text-right">{counts.yes}</span>
	</div>

	<!-- No bar -->
	<div class="flex items-center gap-2">
		<span class="w-12 text-xs font-medium text-gray-600">No</span>
		<div class="flex-1 h-5 bg-gray-100 rounded-[var(--radius-sm)] overflow-hidden">
			<div
				class="h-full bg-red-400 rounded-[var(--radius-sm)] transition-all duration-300"
				style="width: {getBarWidth(counts.no)}"
			></div>
		</div>
		<span class="w-8 text-xs font-medium text-gray-600 text-right">{counts.no}</span>
	</div>

	<!-- Maybe bar -->
	<div class="flex items-center gap-2">
		<span class="w-12 text-xs font-medium text-gray-600">Maybe</span>
		<div class="flex-1 h-5 bg-gray-100 rounded-[var(--radius-sm)] overflow-hidden">
			<div
				class="h-full bg-amber-400 rounded-[var(--radius-sm)] transition-all duration-300"
				style="width: {getBarWidth(counts.maybe)}"
			></div>
		</div>
		<span class="w-8 text-xs font-medium text-gray-600 text-right">{counts.maybe}</span>
	</div>

	<!-- Total count -->
	<p class="text-xs text-gray-500 mt-2">
		{total} {total === 1 ? 'voter' : 'voters'}
	</p>
</div>

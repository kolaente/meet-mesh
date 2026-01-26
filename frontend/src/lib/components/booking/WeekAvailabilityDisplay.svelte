<script lang="ts">
	import type { components } from '$lib/api/types';
	import { getDateFormat } from '$lib/stores/dateFormat.svelte';

	type AvailabilityRule = components['schemas']['AvailabilityRule'];

	interface Props {
		rules: AvailabilityRule[];
	}

	let { rules = [] }: Props = $props();

	const dateFormat = getDateFormat();
	const days = $derived(dateFormat.getWeekDays('short'));
	// Map display order to day-of-week numbers based on week start preference
	// Sunday start: [0, 1, 2, 3, 4, 5, 6] (Sun=0, Mon=1, ..., Sat=6)
	// Monday start: [1, 2, 3, 4, 5, 6, 0] (Mon=1, Tue=2, ..., Sun=0)
	const dayIndexMap = $derived(
		dateFormat.weekStartDay === 'monday'
			? [1, 2, 3, 4, 5, 6, 0]
			: [0, 1, 2, 3, 4, 5, 6]
	);

	// Hours to display (6am to 10pm covers most use cases)
	const startHour = 6;
	const endHour = 22;
	const totalHours = endHour - startHour;

	// Hour labels to show (every 3 hours)
	const hourLabels = [6, 9, 12, 15, 18, 21];

	// Convert time string 'HH:MM' to decimal hours
	function timeToDecimal(time: string): number {
		const [hours, minutes] = time.split(':').map(Number);
		return hours + minutes / 60;
	}

	// Get blocks for a specific day
	function getBlocksForDay(dayIndex: number): { top: number; height: number }[] {
		return rules
			.filter((rule) => rule.days_of_week.includes(dayIndex))
			.map((rule) => {
				const start = Math.max(timeToDecimal(rule.start_time), startHour);
				const end = Math.min(timeToDecimal(rule.end_time), endHour);
				const top = ((start - startHour) / totalHours) * 100;
				const height = ((end - start) / totalHours) * 100;
				return { top, height };
			})
			.filter((block) => block.height > 0);
	}
</script>

<div
	class="border border-gray-200 dark:border-neutral-700 rounded-lg overflow-hidden bg-white dark:bg-neutral-900 relative"
>
	<!-- Header row with day names -->
	<div class="grid grid-cols-[40px_repeat(7,1fr)] border-b border-gray-200 dark:border-neutral-700">
		<div class="p-2"></div>
		{#each days as day}
			<div
				class="p-2 text-center text-xs font-medium text-gray-500 dark:text-gray-400 border-l border-gray-100 dark:border-neutral-800"
			>
				{day}
			</div>
		{/each}
	</div>

	<!-- Calendar body -->
	<div class="grid grid-cols-[40px_repeat(7,1fr)] relative" style="height: 200px;">
		<!-- Hour labels -->
		<div class="relative">
			{#each hourLabels as hour}
				<div
					class="absolute left-0 right-0 text-xs text-gray-400 dark:text-gray-500 text-right pr-2 -translate-y-1/2"
					style="top: {((hour - startHour) / totalHours) * 100}%"
				>
					{String(hour).padStart(2, '0')}:00
				</div>
			{/each}
		</div>

		<!-- Day columns -->
		{#each dayIndexMap as dayIndex}
			<div class="relative border-l border-gray-100 dark:border-neutral-800">
				<!-- Grid lines -->
				{#each hourLabels as hour}
					<div
						class="absolute left-0 right-0 border-t border-gray-50 dark:border-neutral-800"
						style="top: {((hour - startHour) / totalHours) * 100}%"
					></div>
				{/each}

				<!-- Availability blocks -->
				{#each getBlocksForDay(dayIndex) as block}
					<div
						class="absolute left-1 right-1 bg-indigo-500 dark:bg-indigo-600 rounded-sm opacity-80"
						style="top: {block.top}%; height: {block.height}%;"
					></div>
				{/each}
			</div>
		{/each}
	</div>

	<!-- Empty state -->
	{#if rules.length === 0}
		<div
			class="absolute inset-0 flex items-center justify-center pointer-events-none bg-gray-50/50 dark:bg-neutral-800/50"
		>
			<span class="text-sm text-gray-400 dark:text-gray-500">No availability set</span>
		</div>
	{/if}
</div>

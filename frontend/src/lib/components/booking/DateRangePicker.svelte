<script lang="ts">
	import Input from '../ui/Input.svelte';
	import { formatDateShort } from '$lib/utils/dates';

	interface Props {
		availableDates: string[];
		startDate?: string;
		endDate?: string;
		class?: string;
	}

	let {
		availableDates,
		startDate = $bindable(),
		endDate = $bindable(),
		class: className = ''
	}: Props = $props();

	// Compute min/max from available dates
	let minDate = $derived.by(() => {
		if (availableDates.length === 0) return undefined;
		const dates = availableDates.map((d) => d.split('T')[0]).sort();
		return dates[0];
	});

	let maxDate = $derived.by(() => {
		if (availableDates.length === 0) return undefined;
		const dates = availableDates.map((d) => d.split('T')[0]).sort();
		return dates[dates.length - 1];
	});

	// Validation
	let validationError = $derived.by(() => {
		if (!startDate || !endDate) return undefined;
		if (new Date(startDate) > new Date(endDate)) {
			return 'Start date must be before end date';
		}
		return undefined;
	});

	// Format for display
	function formatDateDisplay(dateStr: string | undefined): string {
		if (!dateStr) return '';
		const date = new Date(dateStr);
		return date.toLocaleDateString('en-US', {
			weekday: 'short',
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}

	let dateRangeDisplay = $derived.by(() => {
		if (!startDate || !endDate) return null;
		if (validationError) return null;
		return `${formatDateDisplay(startDate)} - ${formatDateDisplay(endDate)}`;
	});
</script>

<div class="space-y-4 {className}">
	<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
		<div class="space-y-1.5">
			<label for="start-date" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Start Date</label>
			<input
				id="start-date"
				type="date"
				bind:value={startDate}
				min={minDate}
				max={endDate || maxDate}
				class="block w-full rounded-brutalist-md border border-gray-300 dark:border-neutral-600 bg-white dark:bg-neutral-800 px-3 py-2 text-sm text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-0 focus:border-indigo-500"
			/>
		</div>

		<div class="space-y-1.5">
			<label for="end-date" class="block text-sm font-medium text-gray-700 dark:text-gray-300">End Date</label>
			<input
				id="end-date"
				type="date"
				bind:value={endDate}
				min={startDate || minDate}
				max={maxDate}
				class="block w-full rounded-brutalist-md border border-gray-300 dark:border-neutral-600 bg-white dark:bg-neutral-800 px-3 py-2 text-sm text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-0 focus:border-indigo-500"
			/>
		</div>
	</div>

	{#if validationError}
		<p class="text-sm text-red-600 dark:text-red-400">{validationError}</p>
	{/if}

	{#if dateRangeDisplay}
		<div class="bg-gray-50 dark:bg-neutral-800 rounded-brutalist-md p-3 border border-gray-200 dark:border-neutral-700">
			<p class="text-sm text-gray-600 dark:text-gray-400">
				Selected range: <span class="font-medium text-gray-900 dark:text-gray-100">{dateRangeDisplay}</span>
			</p>
		</div>
	{/if}
</div>

<script lang="ts">
	import type { components } from '$lib/api/types';
	import { getDateFormat } from '$lib/stores/dateFormat.svelte';

	type AvailabilityRule = components['schemas']['AvailabilityRule'];

	interface Props {
		rule: AvailabilityRule;
		onchange: (rule: AvailabilityRule) => void;
		onremove?: () => void;
	}

	let { rule, onchange, onremove }: Props = $props();

	const dateFormat = getDateFormat();

	// Day values ordered by week start preference
	// JS Date.getDay(): Sun=0, Mon=1, Tue=2, Wed=3, Thu=4, Fri=5, Sat=6
	const days = $derived.by(() => {
		const dayLabels = dateFormat.getWeekDays('short');
		// Map display index to JS day-of-week value based on week start
		const dayValues = dateFormat.weekStartDay === 'monday'
			? [1, 2, 3, 4, 5, 6, 0] // Mon, Tue, Wed, Thu, Fri, Sat, Sun
			: [0, 1, 2, 3, 4, 5, 6]; // Sun, Mon, Tue, Wed, Thu, Fri, Sat
		return dayLabels.map((label, i) => ({ label, value: dayValues[i] }));
	});

	// Generate time options in 15-minute intervals, 24h format (00:00 to 23:45)
	const timeOptions: string[] = [];
	for (let h = 0; h < 24; h++) {
		for (let m = 0; m < 60; m += 15) {
			timeOptions.push(`${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`);
		}
	}

	function toggleDay(dayValue: number) {
		const currentDays = rule.days_of_week;
		const newDays = currentDays.includes(dayValue)
			? currentDays.filter((d) => d !== dayValue)
			: [...currentDays, dayValue].sort((a, b) => a - b);
		onchange({ ...rule, days_of_week: newDays });
	}

	function updateTime(field: 'start_time' | 'end_time', value: string) {
		onchange({ ...rule, [field]: value });
	}
</script>

<div class="space-y-3">
	<div class="flex items-center justify-between gap-2">
		<!-- Day toggles -->
		<div class="flex gap-1 flex-wrap">
			{#each days as day}
				<button
					type="button"
					onclick={() => toggleDay(day.value)}
					class="px-2.5 py-1.5 text-xs font-medium rounded-[var(--radius-sm)] transition-colors border
						{rule.days_of_week.includes(day.value)
						? 'bg-indigo-500 text-white border-indigo-500'
						: 'bg-[var(--bg-secondary)] text-[var(--text-secondary)] border-[var(--border-color)] hover:bg-[var(--bg-tertiary)]'}"
				>
					{day.label}
				</button>
			{/each}
		</div>

		<!-- Remove button (if removable) -->
		{#if onremove}
			<button
				type="button"
				onclick={onremove}
				class="p-1 text-[var(--text-muted)] hover:text-[var(--rose)] transition-colors flex-shrink-0"
				aria-label="Remove availability rule"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		{/if}
	</div>

	<!-- Time range -->
	<div class="flex items-center gap-2">
		<select
			value={rule.start_time}
			onchange={(e) => updateTime('start_time', e.currentTarget.value)}
			class="block w-28 h-9 rounded-[var(--radius-sm)] border border-[var(--border-color)] bg-[var(--bg-secondary)] px-2 text-sm text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--sky)] focus:border-[var(--sky)]"
		>
			{#each timeOptions as time}
				<option value={time}>{time}</option>
			{/each}
		</select>
		<span class="text-sm text-[var(--text-muted)]">to</span>
		<select
			value={rule.end_time}
			onchange={(e) => updateTime('end_time', e.currentTarget.value)}
			class="block w-28 h-9 rounded-[var(--radius-sm)] border border-[var(--border-color)] bg-[var(--bg-secondary)] px-2 text-sm text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--sky)] focus:border-[var(--sky)]"
		>
			{#each timeOptions as time}
				<option value={time}>{time}</option>
			{/each}
		</select>
	</div>
</div>

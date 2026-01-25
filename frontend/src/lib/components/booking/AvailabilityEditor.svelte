<script lang="ts">
	import type { components } from '$lib/api/types';
	import WeekAvailabilityDisplay from './WeekAvailabilityDisplay.svelte';
	import AvailabilityRuleForm from './AvailabilityRuleForm.svelte';

	type AvailabilityRule = components['schemas']['AvailabilityRule'];

	interface Props {
		rules: AvailabilityRule[];
		onchange: (rules: AvailabilityRule[]) => void;
	}

	let { rules, onchange }: Props = $props();

	function updateRule(index: number, updatedRule: AvailabilityRule) {
		const newRules = rules.map((r, i) => (i === index ? updatedRule : r));
		onchange(newRules);
	}

	function removeRule(index: number) {
		const newRules = rules.filter((_, i) => i !== index);
		onchange(newRules);
	}

	function addRule() {
		const newRule: AvailabilityRule = {
			days_of_week: [1, 2, 3, 4, 5], // Mon-Fri default
			start_time: '09:00',
			end_time: '17:00'
		};
		onchange([...rules, newRule]);
	}
</script>

<div class="space-y-4">
	<span class="block text-sm font-medium text-[var(--text-secondary)]">Availability</span>

	<!-- Visual calendar display -->
	<WeekAvailabilityDisplay {rules} />

	<!-- Rule forms -->
	<div class="space-y-3">
		{#each rules as rule, index (index)}
			<div class="p-3 bg-[var(--bg-secondary)] rounded-lg">
				<AvailabilityRuleForm
					{rule}
					onchange={(r) => updateRule(index, r)}
					onremove={rules.length > 1 ? () => removeRule(index) : undefined}
				/>
			</div>
		{/each}
	</div>

	<!-- Add another rule -->
	<button
		type="button"
		onclick={addRule}
		class="text-sm text-[var(--sky)] hover:text-[var(--sky-hover)] font-medium transition-colors"
	>
		+ Add another time window
	</button>
</div>

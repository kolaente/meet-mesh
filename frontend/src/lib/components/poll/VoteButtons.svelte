<script lang="ts">
	import type { components } from '$lib/api/types';

	type VoteResponse = components['schemas']['VoteResponse'];

	interface Props {
		value?: VoteResponse;
		onChange?: (value: VoteResponse | undefined) => void;
		class?: string;
	}

	let {
		value = $bindable(),
		onChange,
		class: className = ''
	}: Props = $props();

	// VoteResponse: 1=yes, 2=no, 3=maybe
	const options: { value: VoteResponse; label: string }[] = [
		{ value: 1, label: 'Yes' },
		{ value: 2, label: 'No' },
		{ value: 3, label: 'Maybe' }
	];

	function handleClick(optionValue: VoteResponse) {
		if (value === optionValue) {
			// Toggle off if already selected
			value = undefined;
			onChange?.(undefined);
		} else {
			value = optionValue;
			onChange?.(optionValue);
		}
	}

	function getButtonClasses(optionValue: VoteResponse): string {
		const isSelected = value === optionValue;
		const base = 'px-3 sm:px-4 py-2.5 sm:py-2 min-h-[44px] text-sm font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-1 active:scale-[0.98]';

		if (isSelected) {
			switch (optionValue) {
				case 1: // Yes
					return `${base} bg-[var(--emerald)] text-white hover:opacity-90 focus:ring-[var(--emerald)]`;
				case 2: // No
					return `${base} bg-[var(--rose)] text-white hover:opacity-90 focus:ring-[var(--rose)]`;
				case 3: // Maybe
					return `${base} bg-[var(--amber)] text-white hover:opacity-90 focus:ring-[var(--amber)]`;
			}
		}

		// Unselected/ghost style
		return `${base} bg-[var(--bg-secondary)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] border-y border-[var(--border-color)] first:border-l first:rounded-l-[var(--radius-md)] last:border-r last:rounded-r-[var(--radius-md)]`;
	}
</script>

<div
	class="inline-flex w-full sm:w-auto rounded-[var(--radius-md)] shadow-sm {className}"
	role="group"
	aria-label="Vote selection"
>
	{#each options as option (option.value)}
		{@const isSelected = value === option.value}
		<button
			type="button"
			onclick={() => handleClick(option.value)}
			class="{getButtonClasses(option.value)} flex-1 sm:flex-none"
			class:border-l={!isSelected && option.value !== 1}
			class:border-[var(--border-color)]={!isSelected}
			class:border-y={!isSelected}
			class:border-r={!isSelected && option.value === 3}
			class:first:border-l={!isSelected}
			class:rounded-l-[var(--radius-md)]={option.value === 1}
			class:rounded-r-[var(--radius-md)]={option.value === 3}
			aria-pressed={isSelected}
		>
			{option.label}
		</button>
	{/each}
</div>

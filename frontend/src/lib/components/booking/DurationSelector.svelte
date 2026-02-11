<script lang="ts">
	interface Props {
		durations: number[];
		selectedDuration: number;
		class?: string;
	}

	let {
		durations,
		selectedDuration = $bindable(),
		class: className = ''
	}: Props = $props();

	function formatDuration(minutes: number): string {
		if (minutes < 60) {
			return `${minutes} min`;
		}
		const hours = Math.floor(minutes / 60);
		const remainingMins = minutes % 60;
		if (remainingMins === 0) {
			return hours === 1 ? '1 hour' : `${hours} hours`;
		}
		return `${hours}h ${remainingMins}m`;
	}

	function handleSelect(duration: number) {
		selectedDuration = duration;
	}

	function isSelected(duration: number): boolean {
		return selectedDuration === duration;
	}
</script>

<div class="flex flex-wrap gap-2 {className}">
	{#each durations as duration (duration)}
		<button
			type="button"
			onclick={() => handleSelect(duration)}
			class="px-4 py-2 text-sm font-medium rounded-full border transition-colors
				{isSelected(duration)
					? 'bg-accent-sky text-white border-accent-sky'
					: 'bg-bg-secondary text-text-primary border-border hover:border-accent-sky hover:bg-bg-tertiary'}"
		>
			{formatDuration(duration)}
		</button>
	{/each}
</div>

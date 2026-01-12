<script lang="ts">
	import { fly } from 'svelte/transition';
	import type { components } from '$lib/api/types';

	type Slot = components['schemas']['Slot'];

	interface Props {
		slots: Slot[];
		selectedSlot?: Slot;
		onSelect?: (slot: Slot) => void;
		class?: string;
	}

	let {
		slots,
		selectedSlot = $bindable(),
		onSelect,
		class: className = ''
	}: Props = $props();

	function formatTime(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleTimeString('en-US', {
			hour: 'numeric',
			minute: '2-digit',
			hour12: true
		});
	}

	function formatTimeRange(slot: Slot): string {
		return `${formatTime(slot.start_time)} - ${formatTime(slot.end_time)}`;
	}

	function handleSelect(slot: Slot) {
		selectedSlot = slot;
		onSelect?.(slot);
	}

	function isSelected(slot: Slot): boolean {
		return selectedSlot?.id === slot.id;
	}
</script>

<div class="space-y-2 {className}">
	{#each slots as slot, index (slot.id)}
		<button
			type="button"
			onclick={() => handleSelect(slot)}
			in:fly={{ x: 50, duration: 200, delay: index * 50 }}
			class="w-full px-4 py-3 min-h-[48px] text-left rounded-[var(--radius-md)] border transition-colors active:scale-[0.98]
				{isSelected(slot)
					? 'bg-indigo-600 text-white border-indigo-600'
					: 'bg-white text-gray-900 border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 active:bg-indigo-100'}"
		>
			<span class="font-medium text-sm sm:text-base">{formatTimeRange(slot)}</span>
		</button>
	{/each}

	{#if slots.length === 0}
		<p class="text-center text-gray-500 py-4 text-sm sm:text-base">No available time slots for this day.</p>
	{/if}
</div>

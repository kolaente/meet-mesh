<script lang="ts">
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { Button, Card } from '$lib/components/ui';

	type Slot = components['schemas']['Slot'];
	type SlotType = components['schemas']['SlotType'];

	interface Props {
		linkId: number;
		slots: Slot[];
	}

	let { linkId, slots = $bindable([]) }: Props = $props();

	let showAddForm = $state(false);
	let adding = $state(false);
	let deletingId = $state<number | null>(null);
	let error = $state('');

	// Form state
	let slotType = $state<SlotType>(1);
	let date = $state('');
	let startTime = $state('09:00');
	let endTime = $state('10:00');

	const slotTypeOptions = [
		{ value: 1, label: 'Time Slot' },
		{ value: 2, label: 'Full Day' },
		{ value: 3, label: 'Multi Day' }
	];

	function formatSlot(slot: Slot): string {
		const start = new Date(slot.start_time);
		const end = new Date(slot.end_time);

		if (slot.type === 2) {
			// Full day
			return start.toLocaleDateString(undefined, {
				weekday: 'long',
				month: 'short',
				day: 'numeric'
			}) + ' (Full Day)';
		} else if (slot.type === 3) {
			// Multi day
			return start.toLocaleDateString(undefined, { month: 'short', day: 'numeric' }) +
				' - ' +
				end.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
		} else {
			// Time slot
			return start.toLocaleDateString(undefined, {
				weekday: 'long',
				month: 'short',
				day: 'numeric'
			}) + ' - ' +
				start.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' }) +
				' to ' +
				end.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' });
		}
	}

	function resetForm() {
		slotType = 1;
		date = '';
		startTime = '09:00';
		endTime = '10:00';
		error = '';
	}

	async function handleAdd() {
		if (!date) {
			error = 'Please select a date';
			return;
		}

		error = '';
		adding = true;

		try {
			let startDateTime: string;
			let endDateTime: string;

			if (slotType === 2) {
				// Full day: use midnight to midnight
				startDateTime = `${date}T00:00:00Z`;
				endDateTime = `${date}T23:59:59Z`;
			} else if (slotType === 3) {
				// Multi day: use start and end dates
				startDateTime = `${date}T00:00:00Z`;
				endDateTime = `${endTime}T23:59:59Z`; // endTime is used as end date here
			} else {
				// Time slot
				startDateTime = `${date}T${startTime}:00Z`;
				endDateTime = `${date}T${endTime}:00Z`;
			}

			const { data, error: apiError } = await api.POST('/links/{id}/slots', {
				params: { path: { id: linkId } },
				body: {
					type: slotType,
					start_time: startDateTime,
					end_time: endDateTime
				}
			});

			if (apiError) {
				error = 'Failed to add slot';
				return;
			}

			if (data) {
				slots = [...slots, data];
				showAddForm = false;
				resetForm();
			}
		} catch {
			error = 'An unexpected error occurred';
		} finally {
			adding = false;
		}
	}

	async function handleDelete(slotId: number) {
		deletingId = slotId;

		try {
			await api.DELETE('/links/{id}/slots/{slotId}', {
				params: { path: { id: linkId, slotId } }
			});

			slots = slots.filter((s) => s.id !== slotId);
		} catch {
			error = 'Failed to delete slot';
		} finally {
			deletingId = null;
		}
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-medium text-gray-900">
				Time Slots ({slots.length})
			</h2>
			{#if !showAddForm}
				<Button variant="secondary" size="sm" onclick={() => (showAddForm = true)}>
					{#snippet children()}Add Slot{/snippet}
				</Button>
			{/if}
		</div>
	{/snippet}

	{#snippet children()}
		{#if error}
			<div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
				{error}
			</div>
		{/if}

		{#if showAddForm}
			<div class="mb-4 p-4 bg-gray-50 rounded-lg space-y-4">
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
					<div>
						<label for="slot-type" class="block text-sm font-medium text-gray-700 mb-1">
							Type
						</label>
						<select
							id="slot-type"
							bind:value={slotType}
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
						>
							{#each slotTypeOptions as option (option.value)}
								<option value={option.value}>{option.label}</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="slot-date" class="block text-sm font-medium text-gray-700 mb-1">
							{slotType === 3 ? 'Start Date' : 'Date'}
						</label>
						<input
							id="slot-date"
							type="date"
							bind:value={date}
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
						/>
					</div>
				</div>

				{#if slotType === 1}
					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="start-time" class="block text-sm font-medium text-gray-700 mb-1">
								Start Time
							</label>
							<input
								id="start-time"
								type="time"
								bind:value={startTime}
								class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
							/>
						</div>
						<div>
							<label for="end-time" class="block text-sm font-medium text-gray-700 mb-1">
								End Time
							</label>
							<input
								id="end-time"
								type="time"
								bind:value={endTime}
								class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
							/>
						</div>
					</div>
				{:else if slotType === 3}
					<div>
						<label for="end-date" class="block text-sm font-medium text-gray-700 mb-1">
							End Date
						</label>
						<input
							id="end-date"
							type="date"
							bind:value={endTime}
							class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
						/>
					</div>
				{/if}

				<div class="flex gap-3 pt-2">
					<Button variant="secondary" size="sm" onclick={() => { showAddForm = false; resetForm(); }}>
						{#snippet children()}Cancel{/snippet}
					</Button>
					<Button size="sm" loading={adding} onclick={handleAdd}>
						{#snippet children()}Add Slot{/snippet}
					</Button>
				</div>
			</div>
		{/if}

		{#if slots.length === 0}
			<p class="text-gray-500 text-center py-4">
				No time slots yet. Add slots for people to vote on.
			</p>
		{:else}
			<div class="divide-y divide-gray-100">
				{#each slots as slot (slot.id)}
					<div class="flex items-center justify-between py-3">
						<div class="flex items-center gap-3">
							<div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center">
								<svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
							</div>
							<span class="text-gray-900">{formatSlot(slot)}</span>
						</div>
						<button
							type="button"
							onclick={() => handleDelete(slot.id)}
							disabled={deletingId === slot.id}
							class="p-1 text-gray-400 hover:text-red-500 disabled:opacity-50"
							aria-label="Delete slot"
						>
							{#if deletingId === slot.id}
								<svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
							{:else}
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
								</svg>
							{/if}
						</button>
					</div>
				{/each}
			</div>
		{/if}
	{/snippet}
</Card>

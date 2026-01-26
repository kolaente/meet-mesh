<script lang="ts">
	import { fly, fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import Spinner from '../ui/Spinner.svelte';
	import Button from '../ui/Button.svelte';
	import Card from '../ui/Card.svelte';
	import DateCalendar from './DateCalendar.svelte';
	import TimeSlotList from './TimeSlotList.svelte';
	import DayPicker from './DayPicker.svelte';
	import DateRangePicker from './DateRangePicker.svelte';
	import BookingForm from './BookingForm.svelte';

	type Slot = components['schemas']['Slot'];
	type CustomField = components['schemas']['CustomField'];
	type SlotType = components['schemas']['SlotType'];

	interface PublicLink {
		type: 1 | 2;
		name: string;
		description?: string;
		custom_fields?: CustomField[];
		slots: Slot[];
		show_results?: boolean;
		require_email?: boolean;
	}

	interface Props {
		link: PublicLink;
		slug: string;
	}

	let { link, slug }: Props = $props();

	// Step management: 'date' | 'time' | 'form' | 'submitting' | 'confirmed'
	type Step = 'date' | 'time' | 'form' | 'submitting' | 'confirmed';
	let step = $state<Step>('date');

	// Selection state
	let selectedDate = $state<string | undefined>();
	let selectedSlot = $state<Slot | undefined>();
	let startDate = $state<string | undefined>();
	let endDate = $state<string | undefined>();

	// UI state
	let loading = $state(false);
	let error = $state<string | undefined>();
	let submitting = $state(false);

	// Determine slot type from first slot
	let slotType = $derived<SlotType>(link.slots[0]?.type ?? 1);

	// Extract unique available dates from slots
	let availableDates = $derived(
		[...new Set(link.slots.map((s) => s.start_time.split('T')[0]))]
	);

	// Filter slots for selected date (for time slot type)
	let slotsForSelectedDate = $derived.by(() => {
		if (!selectedDate) return [];
		return link.slots.filter((s) => s.start_time.startsWith(selectedDate!));
	});

	// Handle date selection
	function handleDateSelect(date: string) {
		selectedDate = date;
		if (slotType === 1) {
			// Time slot type - move to time selection
			step = 'time';
		} else if (slotType === 2) {
			// Full day type - find the slot for this date and go to form
			selectedSlot = link.slots.find((s) => s.start_time.startsWith(date));
			if (selectedSlot) {
				step = 'form';
			}
		}
	}

	// Handle time slot selection
	function handleSlotSelect(slot: Slot) {
		selectedSlot = slot;
		step = 'form';
	}

	// Handle multi-day date range complete
	function handleDateRangeComplete() {
		if (startDate && endDate) {
			// Find the slot that matches this range (or create a virtual one)
			selectedSlot = link.slots.find(
				(s) =>
					s.start_time.startsWith(startDate!) &&
					s.end_time.startsWith(endDate!)
			) || link.slots[0];
			step = 'form';
		}
	}

	// Handle form submission
	async function handleSubmit(data: { email: string; name?: string; customFields: Record<string, string> }) {
		if (!selectedSlot) return;

		submitting = true;
		error = undefined;

		try {
			const response = await api.POST('/p/booking/{slug}/book', {
				params: { path: { slug } },
				body: {
					start_time: selectedSlot.start_time,
					end_time: selectedSlot.end_time,
					guest_email: data.email,
					guest_name: data.name,
					custom_fields: Object.keys(data.customFields).length > 0 ? data.customFields : undefined
				}
			});

			if (response.error) {
				error = response.error.message || 'Failed to create booking';
				submitting = false;
				return;
			}

			// Redirect to confirmed page
			goto(`/p/booking/${slug}/confirmed`);
		} catch (err) {
			error = 'An unexpected error occurred. Please try again.';
			submitting = false;
		}
	}

	// Navigation helpers
	function goBack() {
		if (step === 'time') {
			step = 'date';
			selectedSlot = undefined;
		} else if (step === 'form') {
			if (slotType === 1) {
				step = 'time';
			} else {
				step = 'date';
			}
			selectedSlot = undefined;
		}
	}

	// Format slot for display
	function formatSlotSummary(slot: Slot): string {
		const start = new Date(slot.start_time);
		const end = new Date(slot.end_time);

		if (slotType === 2) {
			// Full day
			return start.toLocaleDateString('en-US', {
				weekday: 'long',
				month: 'long',
				day: 'numeric',
				year: 'numeric'
			});
		} else if (slotType === 3) {
			// Multi-day
			return `${start.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${end.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}`;
		} else {
			// Time slot
			const dateStr = start.toLocaleDateString('en-US', {
				weekday: 'long',
				month: 'long',
				day: 'numeric'
			});
			const timeStr = `${start.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })} - ${end.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' })}`;
			return `${dateStr} at ${timeStr}`;
		}
	}
</script>

<div class="max-w-2xl mx-auto">
	<!-- Header -->
	<div class="text-center mb-8">
		<h1 class="text-2xl font-bold text-[var(--text-primary)]">{link.name}</h1>
		{#if link.description}
			<p class="mt-2 text-[var(--text-secondary)]">{link.description}</p>
		{/if}
	</div>

	<!-- Step indicator -->
	<div class="flex items-center justify-center gap-2 mb-8">
		{#each ['date', 'time', 'form'] as s, i}
			{@const steps = slotType === 1 ? ['date', 'time', 'form'] : ['date', 'form']}
			{@const stepIndex = steps.indexOf(s)}
			{@const currentIndex = steps.indexOf(step)}
			{@const shouldShow = s !== 'time' || slotType === 1}
			{@const isCurrent = step === s}
			{@const isCompleted = stepIndex !== -1 && stepIndex < currentIndex}
			{@const isUpcoming = stepIndex !== -1 && stepIndex > currentIndex}
			{#if shouldShow}
				<div class="flex items-center gap-2">
					<div
						class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium border
							{isCurrent ? 'bg-[var(--sky)] text-white border-[var(--sky)]' : isCompleted ? 'bg-[var(--emerald)] text-white border-[var(--emerald)]' : 'bg-[var(--bg-tertiary)] text-[var(--text-muted)] border-[var(--border-color)]'}"
					>
						{#if isCompleted}
							<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
								<path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
							</svg>
						{:else}
							{stepIndex + 1}
						{/if}
					</div>
					{#if i < 2 && (s !== 'time' || slotType === 1)}
						<div class="w-8 h-px bg-[var(--border-color)]"></div>
					{/if}
				</div>
			{/if}
		{/each}
	</div>

	<!-- Error message -->
	{#if error}
		<div
			class="mb-6 p-4 bg-red-50 border border-red-200 rounded-[var(--radius-md)] text-red-700"
			transition:fade
		>
			{error}
		</div>
	{/if}

	<!-- Content based on slot type and step -->
	<div class="booking-sections">
	<Card>
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<Spinner size="lg" />
			</div>
		{:else if slotType === 1}
			<!-- Time slot flow -->
			{#if step === 'date'}
				<div in:fly={{ x: -50, duration: 200 }}>
					<h2 class="text-lg font-semibold text-[var(--text-primary)] mb-4">Select a Date</h2>
					<DateCalendar
						{availableDates}
						bind:selectedDate
						class="mx-auto max-w-sm"
					/>
					{#if selectedDate}
						<div class="mt-4 flex justify-end">
							<Button onclick={() => handleDateSelect(selectedDate!)}>
								Continue
							</Button>
						</div>
					{/if}
				</div>
			{:else if step === 'time'}
				<div in:fly={{ x: 50, duration: 200 }}>
					<div class="flex items-center gap-2 mb-4">
						<button
							type="button"
							onclick={goBack}
							class="p-2 text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-secondary)] rounded-[var(--radius-md)] transition-colors"
							aria-label="Go back"
						>
							<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
							</svg>
						</button>
						<h2 class="text-lg font-semibold text-[var(--text-primary)]">Select a Time</h2>
					</div>
					<p class="text-sm text-[var(--text-secondary)] mb-4">
						{new Date(selectedDate + 'T00:00:00').toLocaleDateString('en-US', {
							weekday: 'long',
							month: 'long',
							day: 'numeric'
						})}
					</p>
					<TimeSlotList
						slots={slotsForSelectedDate}
						bind:selectedSlot
						onSelect={handleSlotSelect}
					/>
				</div>
			{:else if step === 'form'}
				<div in:fly={{ x: 50, duration: 200 }}>
					<div class="flex items-center gap-2 mb-4">
						<button
							type="button"
							onclick={goBack}
							class="p-2 text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-secondary)] rounded-[var(--radius-md)] transition-colors"
							aria-label="Go back"
						>
							<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
							</svg>
						</button>
						<h2 class="text-lg font-semibold text-[var(--text-primary)]">Your Details</h2>
					</div>
					{#if selectedSlot}
						<div class="mb-6 p-3 bg-indigo-50 dark:bg-indigo-900/30 rounded-[var(--radius-md)] border border-indigo-100 dark:border-indigo-800">
							<p class="text-sm text-indigo-800 dark:text-indigo-200 font-medium">
								{formatSlotSummary(selectedSlot)}
							</p>
						</div>
					{/if}
					<BookingForm
						customFields={link.custom_fields}
						onSubmit={handleSubmit}
						loading={submitting}
					/>
				</div>
			{/if}
		{:else if slotType === 2}
			<!-- Full day flow -->
			{#if step === 'date'}
				<div in:fly={{ x: -50, duration: 200 }}>
					<h2 class="text-lg font-semibold text-[var(--text-primary)] mb-4">Select a Date</h2>
					<DayPicker
						{availableDates}
						bind:selectedDate
					/>
					{#if selectedDate}
						<div class="mt-4 flex justify-end">
							<Button onclick={() => handleDateSelect(selectedDate!)}>
								Continue
							</Button>
						</div>
					{/if}
				</div>
			{:else if step === 'form'}
				<div in:fly={{ x: 50, duration: 200 }}>
					<div class="flex items-center gap-2 mb-4">
						<button
							type="button"
							onclick={goBack}
							class="p-2 text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-secondary)] rounded-[var(--radius-md)] transition-colors"
							aria-label="Go back"
						>
							<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
							</svg>
						</button>
						<h2 class="text-lg font-semibold text-[var(--text-primary)]">Your Details</h2>
					</div>
					{#if selectedSlot}
						<div class="mb-6 p-3 bg-indigo-50 dark:bg-indigo-900/30 rounded-[var(--radius-md)] border border-indigo-100 dark:border-indigo-800">
							<p class="text-sm text-indigo-800 dark:text-indigo-200 font-medium">
								{formatSlotSummary(selectedSlot)}
							</p>
						</div>
					{/if}
					<BookingForm
						customFields={link.custom_fields}
						onSubmit={handleSubmit}
						loading={submitting}
					/>
				</div>
			{/if}
		{:else if slotType === 3}
			<!-- Multi-day flow -->
			{#if step === 'date'}
				<div in:fly={{ x: -50, duration: 200 }}>
					<h2 class="text-lg font-semibold text-[var(--text-primary)] mb-4">Select Date Range</h2>
					<DateRangePicker
						{availableDates}
						bind:startDate
						bind:endDate
					/>
					{#if startDate && endDate && new Date(startDate) <= new Date(endDate)}
						<div class="mt-4 flex justify-end">
							<Button onclick={handleDateRangeComplete}>
								Continue
							</Button>
						</div>
					{/if}
				</div>
			{:else if step === 'form'}
				<div in:fly={{ x: 50, duration: 200 }}>
					<div class="flex items-center gap-2 mb-4">
						<button
							type="button"
							onclick={goBack}
							class="p-2 text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-secondary)] rounded-[var(--radius-md)] transition-colors"
							aria-label="Go back"
						>
							<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
							</svg>
						</button>
						<h2 class="text-lg font-semibold text-[var(--text-primary)]">Your Details</h2>
					</div>
					{#if selectedSlot}
						<div class="mb-6 p-3 bg-indigo-50 dark:bg-indigo-900/30 rounded-[var(--radius-md)] border border-indigo-100 dark:border-indigo-800">
							<p class="text-sm text-indigo-800 dark:text-indigo-200 font-medium">
								{formatSlotSummary(selectedSlot)}
							</p>
						</div>
					{/if}
					<BookingForm
						customFields={link.custom_fields}
						onSubmit={handleSubmit}
						loading={submitting}
					/>
				</div>
			{/if}
		{/if}
	</Card>

	<!-- No slots available -->
	{#if link.slots.length === 0}
		<Card>
			<div class="text-center py-8">
				<p class="text-[var(--text-secondary)]">No available time slots at the moment.</p>
				<p class="text-sm text-[var(--text-muted)] mt-2">Please check back later.</p>
			</div>
		</Card>
	{/if}
	</div>
</div>

<style>
	.booking-sections {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}
</style>

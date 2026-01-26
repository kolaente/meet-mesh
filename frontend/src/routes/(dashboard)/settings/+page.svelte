<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, CalendarConnectionCard, AddCalendarDialog } from '$lib/components/dashboard';
	import { Card, Button, Spinner, Select } from '$lib/components/ui';
	import { getDateFormat } from '$lib/stores/dateFormat.svelte';

	type CalendarConnection = components['schemas']['CalendarConnection'];

	let calendars = $state<CalendarConnection[]>([]);
	let loading = $state(true);
	let addDialogOpen = $state(false);

	const dateFormat = getDateFormat();

	// Reactive bindings for select components
	let timeFormatValue = $derived(dateFormat.timeFormat);
	let weekStartDayValue = $derived(dateFormat.weekStartDay);

	const timeFormatOptions = [
		{ value: '12h', label: '12-hour (2:00 PM)' },
		{ value: '24h', label: '24-hour (14:00)' }
	];

	const weekStartOptions = [
		{ value: 'sunday', label: 'Sunday' },
		{ value: 'monday', label: 'Monday' }
	];

	onMount(async () => {
		await loadCalendars();
	});

	async function loadCalendars() {
		loading = true;
		try {
			const { data } = await api.GET('/calendars');
			if (data) {
				calendars = data;
			}
		} finally {
			loading = false;
		}
	}

	async function handleCalendarAdded() {
		await loadCalendars();
		addDialogOpen = false;
	}

	async function handleDeleteCalendar(id: number) {
		try {
			await api.DELETE('/calendars/{id}', {
				params: { path: { id } }
			});
			calendars = calendars.filter((c) => c.id !== id);
		} catch (e) {
			// Handle error silently for now
		}
	}
</script>

<DashboardHeader title="Settings">
	{#snippet actions()}
		<Button variant="primary" onclick={() => (addDialogOpen = true)}>
			{#snippet children()}Add Calendar{/snippet}
		</Button>
	{/snippet}
</DashboardHeader>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else}
	<div class="space-y-6">
		<!-- Calendar Connections Section -->
		<section>
			<div class="mb-4">
				<h2 class="text-lg font-medium text-[var(--text-primary)]">Calendar Connections</h2>
				<p class="text-sm text-[var(--text-secondary)]">Connect your CalDAV calendars to check availability</p>
			</div>

			<!-- Calendar List -->
			{#if calendars.length === 0}
				<Card>
					{#snippet children()}
						<div class="text-center py-8">
							<svg
								class="mx-auto h-12 w-12 text-[var(--text-muted)]"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
								/>
							</svg>
							<h3 class="mt-2 text-sm font-medium text-[var(--text-primary)]">No calendars connected</h3>
							<p class="mt-1 text-sm text-[var(--text-secondary)]">
								Connect a CalDAV calendar to enable availability checking.
							</p>
							<div class="mt-6">
								<Button variant="primary" onclick={() => (addDialogOpen = true)}>
									{#snippet children()}Connect your first calendar{/snippet}
								</Button>
							</div>
						</div>
					{/snippet}
				</Card>
			{:else if calendars.length > 0}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each calendars as calendar (calendar.id)}
						<CalendarConnectionCard {calendar} onDelete={handleDeleteCalendar} />
					{/each}
				</div>
			{/if}
		</section>

		<!-- Date & Time Format Section -->
		<section>
			<div class="mb-4">
				<h2 class="text-lg font-medium text-[var(--text-primary)]">Date & Time Format</h2>
				<p class="text-sm text-[var(--text-secondary)]">Customize how dates and times are displayed</p>
			</div>

			<Card>
				{#snippet children()}
					<div class="space-y-6">
						<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
							<div>
								<p class="font-medium text-[var(--text-primary)]">Time Format</p>
								<p class="text-sm text-[var(--text-secondary)]">Choose 12-hour or 24-hour time display</p>
							</div>
							<div class="w-full sm:w-48">
								<Select
									name="timeFormat"
									options={timeFormatOptions}
									value={timeFormatValue}
									onchange={(value) => dateFormat.setTimeFormat(value as '12h' | '24h')}
								/>
							</div>
						</div>

						<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
							<div>
								<p class="font-medium text-[var(--text-primary)]">Week Starts On</p>
								<p class="text-sm text-[var(--text-secondary)]">Choose which day starts your week</p>
							</div>
							<div class="w-full sm:w-48">
								<Select
									name="weekStartDay"
									options={weekStartOptions}
									value={weekStartDayValue}
									onchange={(value) => dateFormat.setWeekStartDay(value as 'sunday' | 'monday')}
								/>
							</div>
						</div>

						<div class="pt-2 border-t border-[var(--border-color)]">
							<Button variant="secondary" onclick={() => dateFormat.reset()}>
								{#snippet children()}Reset to Browser Defaults{/snippet}
							</Button>
						</div>
					</div>
				{/snippet}
			</Card>
		</section>
	</div>
{/if}

<AddCalendarDialog
	bind:open={addDialogOpen}
	onSuccess={handleCalendarAdded}
	onClose={() => addDialogOpen = false}
/>

<script lang="ts">
	import type { components } from '$lib/api/types';
	import { Dialog, Button, Input, Spinner, Checkbox } from '$lib/components/ui';
	import { api } from '$lib/api/client';

	type DiscoveredCalendar = components['schemas']['DiscoveredCalendar'];

	type Step = 'credentials' | 'discovering' | 'selection' | 'saving';

	interface Props {
		open?: boolean;
		onSuccess: () => void;
		onClose: () => void;
	}

	let { open = $bindable(false), onSuccess, onClose }: Props = $props();

	// Form state
	let serverUrl = $state('');
	let username = $state('');
	let password = $state('');

	// Flow state
	let step = $state<Step>('credentials');
	let error = $state('');

	// Discovery results
	let discoveredCalendars = $state<DiscoveredCalendar[]>([]);
	let selectedCalendarUrls = $state<Set<string>>(new Set());
	let writeCalendarUrl = $state<string | undefined>(undefined);

	function resetForm() {
		serverUrl = '';
		username = '';
		password = '';
		step = 'credentials';
		error = '';
		discoveredCalendars = [];
		selectedCalendarUrls = new Set();
		writeCalendarUrl = undefined;
	}

	function handleClose() {
		resetForm();
		onClose();
	}

	async function handleDiscover() {
		error = '';
		step = 'discovering';

		const { data, error: apiError } = await api.POST('/calendars/discover', {
			body: {
				server_url: serverUrl,
				username,
				password
			}
		});

		if (apiError) {
			error = 'Failed to connect to calendar server';
			step = 'credentials';
			return;
		}

		if (!data.success) {
			error = data.error || 'Failed to discover calendars';
			step = 'credentials';
			return;
		}

		if (data.calendars.length === 0) {
			error = 'No calendars found on this server';
			step = 'credentials';
			return;
		}

		discoveredCalendars = data.calendars;
		// Pre-select all calendars for availability
		selectedCalendarUrls = new Set(data.calendars.map((c) => c.url));
		step = 'selection';
	}

	function toggleCalendarSelection(url: string) {
		const newSet = new Set(selectedCalendarUrls);
		if (newSet.has(url)) {
			newSet.delete(url);
			// If we're deselecting the write calendar, clear it
			if (writeCalendarUrl === url) {
				writeCalendarUrl = undefined;
			}
		} else {
			newSet.add(url);
		}
		selectedCalendarUrls = newSet;
	}

	function setWriteCalendar(url: string) {
		// Toggle: if already selected, deselect; otherwise select
		if (writeCalendarUrl === url) {
			writeCalendarUrl = undefined;
		} else {
			writeCalendarUrl = url;
			// Ensure the write calendar is also in the availability selection
			if (!selectedCalendarUrls.has(url)) {
				const newSet = new Set(selectedCalendarUrls);
				newSet.add(url);
				selectedCalendarUrls = newSet;
			}
		}
	}

	async function handleSave() {
		if (selectedCalendarUrls.size === 0) {
			error = 'Please select at least one calendar for availability';
			return;
		}

		error = '';
		step = 'saving';

		const { error: apiError } = await api.POST('/calendars', {
			body: {
				server_url: serverUrl,
				username,
				password,
				calendar_urls: Array.from(selectedCalendarUrls),
				write_url: writeCalendarUrl
			}
		});

		if (apiError) {
			error = 'Failed to save calendar connection';
			step = 'selection';
			return;
		}

		resetForm();
		onSuccess();
	}

	function goBack() {
		if (step === 'selection') {
			step = 'credentials';
		}
	}

	// Reactive title based on step
	let dialogTitle = $derived(
		step === 'credentials'
			? 'Add Calendar Connection'
			: step === 'discovering'
				? 'Discovering Calendars'
				: step === 'selection'
					? 'Select Calendars'
					: 'Saving Connection'
	);
</script>

<Dialog bind:open title={dialogTitle}>
	{#snippet children()}
		{#if step === 'credentials'}
			<div class="space-y-4">
				<p class="text-text-secondary text-sm">
					Enter your CalDAV server credentials to connect your calendars.
				</p>

				{#if error}
					<div
						class="rounded-lg border-2 border-red-500 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm"
					>
						{error}
					</div>
				{/if}

				<Input
					name="server_url"
					label="Server URL"
					type="url"
					placeholder="https://caldav.example.com"
					bind:value={serverUrl}
					required
				/>

				<Input
					name="username"
					label="Username"
					type="text"
					placeholder="your-username"
					bind:value={username}
					required
				/>

				<Input
					name="password"
					label="Password"
					type="password"
					placeholder="your-password"
					bind:value={password}
					required
				/>
			</div>
		{:else if step === 'discovering'}
			<div class="flex flex-col items-center justify-center py-8">
				<Spinner size="lg" />
				<p class="mt-4 text-text-secondary">Discovering calendars...</p>
			</div>
		{:else if step === 'selection'}
			<div class="space-y-4">
				<p class="text-text-secondary text-sm">
					Select which calendars to use for availability checking, and optionally choose one for
					writing new events.
				</p>

				{#if error}
					<div
						class="rounded-lg border-2 border-red-500 bg-red-50 dark:bg-red-900/20 px-4 py-3 text-red-700 dark:text-red-400 text-sm"
					>
						{error}
					</div>
				{/if}

				<div class="max-h-64 overflow-y-auto rounded-lg border-2 border-border">
					{#each discoveredCalendars as calendar}
						<div
							class="flex items-center justify-between px-4 py-3 border-b border-border last:border-b-0"
						>
							<div class="flex items-center gap-3 flex-1 min-w-0">
								<Checkbox
									checked={selectedCalendarUrls.has(calendar.url)}
									onchange={() => toggleCalendarSelection(calendar.url)}
								/>
								<div class="min-w-0">
									<p class="font-medium text-text-primary truncate">
										{calendar.name}
									</p>
									{#if calendar.description}
										<p class="text-xs text-text-secondary truncate">
											{calendar.description}
										</p>
									{/if}
								</div>
							</div>
							<button
								type="button"
								class="ml-2 px-2 py-1 text-xs rounded-md transition-colors {writeCalendarUrl ===
								calendar.url
									? 'bg-accent-sky text-white'
									: 'bg-bg-tertiary text-text-secondary hover:bg-bg-secondary'}"
								onclick={() => setWriteCalendar(calendar.url)}
								title={writeCalendarUrl === calendar.url
									? 'Click to remove as write calendar'
									: 'Set as write calendar'}
							>
								{writeCalendarUrl === calendar.url ? 'Write' : 'Set write'}
							</button>
						</div>
					{/each}
				</div>

				<div class="text-xs text-text-muted space-y-1">
					<p>
						<strong>Availability:</strong> Selected calendars will be checked for busy times when generating
						available slots.
					</p>
					<p>
						<strong>Write calendar:</strong> New bookings will be added to this calendar. Leave unset
						to skip writing events.
					</p>
				</div>
			</div>
		{:else if step === 'saving'}
			<div class="flex flex-col items-center justify-center py-8">
				<Spinner size="lg" />
				<p class="mt-4 text-text-secondary">Saving calendar connection...</p>
			</div>
		{/if}
	{/snippet}

	{#snippet footer()}
		{#if step === 'credentials'}
			<Button variant="secondary" onclick={handleClose}>Cancel</Button>
			<Button
				variant="primary"
				onclick={handleDiscover}
				disabled={!serverUrl || !username || !password}
			>
				{#snippet children()}
					Connect
				{/snippet}
			</Button>
		{:else if step === 'selection'}
			<Button variant="secondary" onclick={goBack}>Back</Button>
			<Button variant="primary" onclick={handleSave} disabled={selectedCalendarUrls.size === 0}>
				{#snippet children()}
					Save Connection
				{/snippet}
			</Button>
		{/if}
	{/snippet}
</Dialog>

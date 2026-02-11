<script lang="ts">
	import type { components } from '$lib/api/types';
	import { Card, Button, Badge, Dialog, Spinner } from '$lib/components/ui';
	import { api } from '$lib/api/client';

	type CalendarConnection = components['schemas']['CalendarConnection'];
	type CalendarTestResult = components['schemas']['CalendarTestResult'];

	interface Props {
		calendar: CalendarConnection;
		onDelete: (id: number) => void;
	}

	let { calendar, onDelete }: Props = $props();

	let deleting = $state(false);
	let testing = $state(false);
	let testDialogOpen = $state(false);
	let testResult = $state<CalendarTestResult | null>(null);

	async function handleDelete() {
		deleting = true;
		onDelete(calendar.id);
	}

	async function handleTest() {
		testing = true;
		testDialogOpen = true;
		testResult = null;

		const { data, error } = await api.POST('/calendars/{id}/test', {
			params: { path: { id: calendar.id } }
		});

		if (error) {
			testResult = { success: false, error: 'Failed to test calendar connection', events: [] };
		} else {
			testResult = data;
		}

		testing = false;
	}

	function formatDateTime(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleString(undefined, {
			weekday: 'short',
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		});
	}

	// Extract the hostname from the server URL for display
	function getProviderName(url: string): string {
		try {
			const hostname = new URL(url).hostname;
			return hostname;
		} catch {
			return url;
		}
	}
</script>

<Card>
	{#snippet children()}
		<div class="space-y-3">
			<div class="flex items-start justify-between gap-2">
				<div class="min-w-0 flex-1">
					<h3 class="font-medium text-text-primary truncate">{getProviderName(calendar.server_url)}</h3>
					<p class="text-sm text-text-secondary truncate">{calendar.username}</p>
				</div>
				<div class="flex items-center gap-2">
					<Badge variant="caldav" size="xs" />
				</div>
			</div>

			<div class="text-sm text-text-secondary">
				<span class="font-mono text-xs truncate block">{calendar.server_url}</span>
			</div>

			{#if calendar.calendar_urls && calendar.calendar_urls.length > 0}
				<div class="text-sm text-text-secondary">
					<span class="font-medium">{calendar.calendar_urls.length}</span>
					calendar{calendar.calendar_urls.length === 1 ? '' : 's'} synced
				</div>
			{/if}
			{#if calendar.write_url}
				<div class="flex items-center gap-1 text-sm text-green-600 dark:text-green-400">
					<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
					</svg>
					<span>Write enabled</span>
				</div>
			{/if}

			<div class="flex justify-end gap-2 pt-2">
				<Button
					variant="secondary"
					size="sm"
					onclick={handleTest}
				>
					{#snippet children()}
						Test
					{/snippet}
				</Button>
				<Button
					variant="danger"
					size="sm"
					onclick={handleDelete}
					disabled={deleting}
					loading={deleting}
				>
					{#snippet children()}
						{deleting ? 'Removing...' : 'Remove'}
					{/snippet}
				</Button>
			</div>
		</div>
	{/snippet}
</Card>

<Dialog bind:open={testDialogOpen} title="Calendar Test Results">
	{#snippet children()}
		{#if testing}
			<div class="flex flex-col items-center justify-center py-8">
				<Spinner size="lg" />
				<p class="mt-4 text-text-secondary">Testing calendar connection...</p>
			</div>
		{:else if testResult}
			{#if testResult.success}
				<div class="space-y-4">
					<div class="flex items-center gap-2 text-green-600 dark:text-green-400">
						<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						<span class="font-medium">Connection successful</span>
					</div>
					<p class="text-text-secondary">
						Found {testResult.events.length} event{testResult.events.length === 1 ? '' : 's'} in the next 7 days
					</p>
					{#if testResult.events.length > 0}
						<div class="max-h-64 overflow-y-auto rounded-lg border border-gray-200 dark:border-neutral-700">
							<ul class="divide-y divide-gray-200 dark:divide-neutral-700">
								{#each testResult.events as event}
									<li class="px-4 py-3">
										<p class="font-medium text-text-primary">{event.title}</p>
										<p class="text-sm text-text-secondary">
											{formatDateTime(event.start)} - {formatDateTime(event.end)}
										</p>
									</li>
								{/each}
							</ul>
						</div>
					{/if}
				</div>
			{:else}
				<div class="space-y-4">
					<div class="flex items-center gap-2 text-red-600 dark:text-red-400">
						<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						<span class="font-medium">Connection failed</span>
					</div>
					<p class="text-text-secondary">{testResult.error}</p>
				</div>
			{/if}
		{/if}
	{/snippet}
	{#snippet footer()}
		<Button variant="secondary" onclick={() => testDialogOpen = false}>Close</Button>
	{/snippet}
</Dialog>

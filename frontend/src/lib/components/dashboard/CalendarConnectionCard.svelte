<script lang="ts">
	import type { components } from '$lib/api/types';
	import { Card, Button } from '$lib/components/ui';

	type CalendarConnection = components['schemas']['CalendarConnection'];

	interface Props {
		calendar: CalendarConnection;
		onDelete: (id: number) => void;
	}

	let { calendar, onDelete }: Props = $props();

	let deleting = $state(false);

	async function handleDelete() {
		deleting = true;
		onDelete(calendar.id);
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
					<h3 class="font-medium text-[var(--text-primary)] truncate">{getProviderName(calendar.server_url)}</h3>
					<p class="text-sm text-[var(--text-secondary)] truncate">{calendar.username}</p>
				</div>
				<div class="flex items-center gap-2">
					<span class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full bg-[var(--emerald)] text-white border border-[var(--border-color)]">
						CalDAV
					</span>
				</div>
			</div>

			<div class="text-sm text-[var(--text-secondary)]">
				<span class="font-mono text-xs truncate block">{calendar.server_url}</span>
			</div>

			<div class="flex justify-end pt-2">
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

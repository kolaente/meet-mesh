<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader } from '$lib/components/dashboard';
	import { Button, Spinner, Card, Badge } from '$lib/components/ui';

	type BookingLink = components['schemas']['BookingLink'];

	let bookingLinks = $state<BookingLink[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const { data } = await api.GET('/booking-links');
			if (data) {
				bookingLinks = data;
			}
		} finally {
			loading = false;
		}
	});
</script>

<DashboardHeader title="Booking Links">
	{#snippet actions()}
		<Button variant="primary" onclick={() => (window.location.href = '/booking-links/new')}>
			{#snippet children()}New Booking Link{/snippet}
		</Button>
	{/snippet}
</DashboardHeader>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if bookingLinks.length === 0}
	<div class="text-center py-12">
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
		<h3 class="mt-2 text-sm font-medium text-[var(--text-primary)]">No booking links</h3>
		<p class="mt-1 text-sm text-[var(--text-secondary)]">Get started by creating a booking link for 1:1 scheduling.</p>
		<div class="mt-6">
			<Button variant="primary" onclick={() => (window.location.href = '/booking-links/new')}>
				{#snippet children()}Create your first booking link{/snippet}
			</Button>
		</div>
	</div>
{:else}
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
		{#each bookingLinks as link (link.id)}
			<a href="/booking-links/{link.id}" class="block">
				<Card class="hover:shadow-md transition-shadow cursor-pointer active:shadow-sm">
					{#snippet children()}
						<div class="space-y-2 sm:space-y-3">
							<div class="flex items-start justify-between gap-2">
								<h3 class="font-medium text-[var(--text-primary)] truncate text-sm sm:text-base">{link.name}</h3>
								<Badge variant={link.status === 1 ? 'active' : 'cancelled'} size="sm" />
							</div>

							<div class="text-xs sm:text-sm text-[var(--text-secondary)] truncate">
								<span class="font-mono">/p/booking/{link.slug}</span>
							</div>
						</div>
					{/snippet}
				</Card>
			</a>
		{/each}
	</div>
{/if}

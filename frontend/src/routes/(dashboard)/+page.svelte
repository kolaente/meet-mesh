<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, StatsCard, BookingRow } from '$lib/components/dashboard';
	import { Card, Button, Spinner } from '$lib/components/ui';

	type Link = components['schemas']['Link'];
	type Booking = components['schemas']['Booking'];

	let links = $state<Link[]>([]);
	let recentBookings = $state<Booking[]>([]);
	let loading = $state(true);

	const totalLinks = $derived(links.length);
	const pendingBookings = $derived(recentBookings.filter((b) => b.status === 1).length);
	const totalBookings = $derived(recentBookings.length);

	onMount(async () => {
		try {
			const [linksRes, ...bookingPromises] = await Promise.all([
				api.GET('/links'),
				// We'll fetch bookings from links after we have them
			]);

			if (linksRes.data) {
				links = linksRes.data;

				// Fetch bookings for booking-type links (type === 1)
				const bookingLinks = links.filter((l) => l.type === 1).slice(0, 5);
				const bookingResults = await Promise.all(
					bookingLinks.map((l) =>
						api.GET('/links/{id}/bookings', { params: { path: { id: l.id } } })
					)
				);

				const allBookings: Booking[] = [];
				for (const res of bookingResults) {
					if (res.data) {
						allBookings.push(...res.data);
					}
				}

				// Sort by created_at descending and take recent
				recentBookings = allBookings
					.sort((a, b) => {
						const dateA = a.created_at ? new Date(a.created_at).getTime() : 0;
						const dateB = b.created_at ? new Date(b.created_at).getTime() : 0;
						return dateB - dateA;
					})
					.slice(0, 10);
			}
		} finally {
			loading = false;
		}
	});

	function handleBookingUpdate(updated: Booking) {
		recentBookings = recentBookings.map((b) => (b.id === updated.id ? updated : b));
	}
</script>

<DashboardHeader title="Dashboard">
	{#snippet actions()}
		<Button variant="primary" onclick={() => (window.location.href = '/links/new')}>
			{#snippet children()}New Link{/snippet}
		</Button>
	{/snippet}
</DashboardHeader>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else}
	<div class="space-y-4 sm:space-y-6">
		<!-- Stats - stack on mobile, 3 columns on tablet and up -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
			<StatsCard label="Total Links" value={totalLinks}>
				{#snippet icon()}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
						/>
					</svg>
				{/snippet}
			</StatsCard>

			<StatsCard label="Pending Bookings" value={pendingBookings}>
				{#snippet icon()}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
				{/snippet}
			</StatsCard>

			<StatsCard label="Total Bookings" value={totalBookings}>
				{#snippet icon()}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
						/>
					</svg>
				{/snippet}
			</StatsCard>
		</div>

		<!-- Recent Bookings -->
		<Card>
			{#snippet header()}
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-medium text-gray-900">Recent Bookings</h2>
					<a href="/links" class="text-sm text-indigo-600 hover:text-indigo-700">View all links</a>
				</div>
			{/snippet}

			{#snippet children()}
				{#if recentBookings.length === 0}
					<p class="text-gray-500 text-center py-4">No bookings yet</p>
				{:else}
					<div class="divide-y divide-gray-100">
						{#each recentBookings as booking (booking.id)}
							<BookingRow {booking} onUpdate={handleBookingUpdate} />
						{/each}
					</div>
				{/if}
			{/snippet}
		</Card>
	</div>
{/if}

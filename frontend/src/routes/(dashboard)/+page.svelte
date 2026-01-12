<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, StatsCard, BookingRow } from '$lib/components/dashboard';
	import { Card, Button, Spinner } from '$lib/components/ui';

	type BookingLink = components['schemas']['BookingLink'];
	type Poll = components['schemas']['Poll'];
	type Booking = components['schemas']['Booking'];

	let bookingLinks = $state<BookingLink[]>([]);
	let polls = $state<Poll[]>([]);
	let recentBookings = $state<Booking[]>([]);
	let loading = $state(true);

	const totalLinks = $derived(bookingLinks.length + polls.length);
	const pendingBookings = $derived(recentBookings.filter((b) => b.status === 1).length);
	const totalBookings = $derived(recentBookings.length);

	onMount(async () => {
		try {
			const [bookingLinksRes, pollsRes] = await Promise.all([
				api.GET('/booking-links'),
				api.GET('/polls')
			]);

			if (bookingLinksRes.data) {
				bookingLinks = bookingLinksRes.data;

				// Fetch bookings for booking links
				const bookingResults = await Promise.all(
					bookingLinks.slice(0, 5).map((l) =>
						api.GET('/booking-links/{id}/bookings', { params: { path: { id: l.id } } })
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

			if (pollsRes.data) {
				polls = pollsRes.data;
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
		<Button variant="secondary" onclick={() => (window.location.href = '/booking-links/new')}>
			{#snippet children()}New Booking Link{/snippet}
		</Button>
		<Button variant="primary" onclick={() => (window.location.href = '/polls/new')}>
			{#snippet children()}New Poll{/snippet}
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
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
			<StatsCard label="Booking Links" value={bookingLinks.length}>
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

			<StatsCard label="Polls" value={polls.length}>
				{#snippet icon()}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
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
					<a href="/booking-links" class="text-sm text-indigo-600 hover:text-indigo-700">View all booking links</a>
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

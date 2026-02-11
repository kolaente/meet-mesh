<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import {
		DashboardHeader,
		StatsCard,
		BookingRow,
		QuickActions
	} from '$lib/components/dashboard';
	import { Card, Button, Spinner } from '$lib/components/ui';

	type BookingLink = components['schemas']['BookingLink'];
	type Poll = components['schemas']['Poll'];
	type Booking = components['schemas']['Booking'];

	let bookingLinks = $state<BookingLink[]>([]);
	let polls = $state<Poll[]>([]);
	let recentBookings = $state<Booking[]>([]);
	let loading = $state(true);

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

	// Quick actions configuration
	const quickActions = [
		{
			href: '/booking-links/new',
			title: 'Create Booking Link',
			description: 'Set up 1:1 scheduling',
			icon: plusIcon
		},
		{
			href: '/polls/new',
			title: 'Start New Poll',
			description: 'Coordinate group availability',
			icon: usersIcon
		},
		{
			href: '/settings',
			title: 'Sync Calendar',
			description: 'Update availability',
			icon: syncIcon
		}
	];
</script>

{#snippet plusIcon()}
	<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path stroke-linecap="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
	</svg>
{/snippet}

{#snippet usersIcon()}
	<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path
			stroke-linecap="round"
			stroke-width="2"
			d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"
		/>
	</svg>
{/snippet}

{#snippet syncIcon()}
	<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
		<path
			stroke-linecap="round"
			stroke-width="2"
			d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
		/>
	</svg>
{/snippet}

<DashboardHeader title="Dashboard" subtitle="Welcome back! Here's what's happening.">
	{#snippet actions()}
		<Button variant="secondary" onclick={() => goto('/booking-links/new')}>
			{#snippet children()}
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
					/>
				</svg>
				New Link
			{/snippet}
		</Button>
		<Button variant="primary" onclick={() => goto('/polls/new')}>
			{#snippet children()}
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
				</svg>
				New Poll
			{/snippet}
		</Button>
	{/snippet}
</DashboardHeader>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else}
	<div class="space-y-6">
		<!-- Stats Grid: 4 columns on desktop, 2 on tablet/mobile -->
		<div class="stats-grid">
			<StatsCard label="Booking Links" value={bookingLinks.length} color="sky">
				{#snippet icon()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
						/>
					</svg>
				{/snippet}
			</StatsCard>

			<StatsCard label="Active Polls" value={polls.length} color="violet">
				{#snippet icon()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-width="2"
							d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
						/>
					</svg>
				{/snippet}
			</StatsCard>

			<StatsCard
				label="Pending Approval"
				value={pendingBookings}
				color="amber"
				trend={pendingBookings > 0 ? { value: `${pendingBookings} new`, type: 'alert' } : undefined}
			>
				{#snippet icon()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-width="2"
							d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
				{/snippet}
			</StatsCard>

			<StatsCard label="Total Bookings" value={totalBookings} color="emerald">
				{#snippet icon()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-width="2"
							d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
						/>
					</svg>
				{/snippet}
			</StatsCard>
		</div>

		<!-- Two-column content grid -->
		<div class="content-grid">
			<!-- Left: Recent Bookings -->
			<Card>
				{#snippet header()}
					<div class="flex items-center justify-between w-full">
						<span class="card-title">Recent Bookings</span>
						{#if pendingBookings > 0}
							<span class="card-badge">{pendingBookings} pending</span>
						{/if}
					</div>
				{/snippet}

				{#snippet children()}
					{#if recentBookings.length === 0}
						<p class="text-[var(--text-muted)] text-center py-4">No bookings yet</p>
					{:else}
						<div class="booking-list">
							{#each recentBookings as booking (booking.id)}
								<BookingRow {booking} onUpdate={handleBookingUpdate} />
							{/each}
						</div>
					{/if}
				{/snippet}
			</Card>

			<!-- Right column: stacked cards -->
			<div class="right-column">
				<!-- Quick Actions -->
				<Card>
					{#snippet header()}
						<span class="card-title">Quick Actions</span>
					{/snippet}

					{#snippet children()}
					<div class="component-fullbleed">
						<QuickActions actions={quickActions} />
					</div>
				{/snippet}
				</Card>
			</div>
		</div>
	</div>
{/if}

<style>
	/* Stats Grid */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 1rem;
	}

	/* Content Grid: 1.6fr 1fr like the prototype */
	.content-grid {
		display: grid;
		grid-template-columns: 1.6fr 1fr;
		gap: 1.5rem;
	}

	/* Right column stacked cards */
	.right-column {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	/* Card header typography */
	.card-title {
		font-size: 0.9rem;
		font-weight: 700;
		color: var(--text-primary);
	}

	.card-badge {
		font-size: 0.7rem;
		font-weight: 700;
		padding: 0.25rem 0.6rem;
		border-radius: var(--radius);
		background: var(--amber);
		color: white;
		border: 1px solid var(--border-color);
	}

	/* Remove default Card padding for components with their own padding */
	.booking-list,
	.component-fullbleed {
		margin: -1rem -1.25rem;
	}

	/* Responsive: tablet */
	@media (max-width: 1100px) {
		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.content-grid {
			grid-template-columns: 1fr;
		}
	}

	/* Responsive: mobile */
	@media (max-width: 600px) {
		.stats-grid {
			grid-template-columns: 1fr 1fr;
		}
	}
</style>

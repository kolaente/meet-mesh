<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, BookingRow } from '$lib/components/dashboard';
	import { Button, Card, Badge, Spinner } from '$lib/components/ui';

	type BookingLink = components['schemas']['BookingLink'];
	type Booking = components['schemas']['Booking'];

	let link = $state<BookingLink | null>(null);
	let bookings = $state<Booking[]>([]);
	let loading = $state(true);
	let deleting = $state(false);
	let copySuccess = $state(false);

	const linkId = $derived(Number($page.params.id));

	const statusVariant = $derived<'active' | 'cancelled'>(link?.status === 1 ? 'active' : 'cancelled');

	const publicUrl = $derived(link ? `${window.location.origin}/p/booking/${link.slug}` : '');

	onMount(async () => {
		try {
			const { data: linkData } = await api.GET('/booking-links/{id}', {
				params: { path: { id: linkId } }
			});

			if (!linkData) {
				goto('/booking-links');
				return;
			}

			link = linkData;

			// Fetch bookings
			const { data: bookingsData } = await api.GET('/booking-links/{id}/bookings', {
				params: { path: { id: linkId } }
			});
			if (bookingsData) {
				bookings = bookingsData;
			}
		} finally {
			loading = false;
		}
	});

	async function copyLink() {
		await navigator.clipboard.writeText(publicUrl);
		copySuccess = true;
		setTimeout(() => {
			copySuccess = false;
		}, 2000);
	}

	async function handleDelete() {
		if (!confirm('Are you sure you want to delete this booking link? This action cannot be undone.')) {
			return;
		}

		deleting = true;
		try {
			await api.DELETE('/booking-links/{id}', {
				params: { path: { id: linkId } }
			});
			goto('/booking-links');
		} finally {
			deleting = false;
		}
	}

	function handleBookingUpdate(updated: Booking) {
		bookings = bookings.map((b) => (b.id === updated.id ? updated : b));
	}
</script>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if link}
	<DashboardHeader title={link.name}>
		{#snippet actions()}
			<Button variant="secondary" onclick={copyLink}>
				{#snippet children()}
					{#if copySuccess}
						Copied!
					{:else}
						Copy Link
					{/if}
				{/snippet}
			</Button>
			<Button variant="secondary" onclick={() => goto(`/booking-links/${linkId}/edit`)}>
				{#snippet children()}Edit{/snippet}
			</Button>
			<Button variant="danger" onclick={handleDelete} loading={deleting}>
				{#snippet children()}Delete{/snippet}
			</Button>
		{/snippet}
	</DashboardHeader>

	<div class="space-y-6">
		<!-- Link Info -->
		<Card>
			{#snippet children()}
				<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
					<div>
						<p class="text-sm font-medium text-gray-500">Type</p>
						<span class="inline-flex items-center px-2.5 py-1 text-sm font-medium rounded-full bg-green-50 text-green-700 border border-green-200">
							Booking
						</span>
					</div>

					<div>
						<p class="text-sm font-medium text-gray-500">Status</p>
						<Badge variant={statusVariant} size="md" />
					</div>

					<div class="col-span-2">
						<p class="text-sm font-medium text-gray-500">Public URL</p>
						<a
							href={publicUrl}
							target="_blank"
							rel="noopener noreferrer"
							class="text-indigo-600 hover:text-indigo-700 font-mono text-sm"
						>
							{publicUrl}
						</a>
					</div>

					{#if link?.description}
						<div class="col-span-2 md:col-span-4">
							<p class="text-sm font-medium text-gray-500">Description</p>
							<p class="text-gray-700">{link?.description}</p>
						</div>
					{/if}
				</div>
			{/snippet}
		</Card>

		<!-- Bookings -->
		<Card>
			{#snippet header()}
				<h2 class="text-lg font-medium text-gray-900">
					Bookings ({bookings.length})
				</h2>
			{/snippet}

			{#snippet children()}
				{#if bookings.length === 0}
					<p class="text-gray-500 text-center py-4">No bookings yet</p>
				{:else}
					<div class="divide-y divide-gray-100">
						{#each bookings as booking (booking.id)}
							<BookingRow {booking} onUpdate={handleBookingUpdate} />
						{/each}
					</div>
				{/if}
			{/snippet}
		</Card>
	</div>
{/if}

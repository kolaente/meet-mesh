<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader, BookingRow, VoteRow, SlotManager } from '$lib/components/dashboard';
	import { Button, Card, Badge, Spinner } from '$lib/components/ui';

	type Link = components['schemas']['Link'];
	type Booking = components['schemas']['Booking'];
	type Vote = components['schemas']['Vote'];
	type Slot = components['schemas']['Slot'];

	let link = $state<Link | null>(null);
	let bookings = $state<Booking[]>([]);
	let votes = $state<Vote[]>([]);
	let slots = $state<Slot[]>([]);
	let loading = $state(true);
	let deleting = $state(false);
	let copySuccess = $state(false);

	const linkId = $derived(Number($page.params.id));

	const typeLabels: Record<1 | 2, string> = {
		1: 'Booking',
		2: 'Poll'
	};

	const statusVariant = $derived<'active' | 'cancelled'>(link?.status === 1 ? 'active' : 'cancelled');

	const publicUrl = $derived(link ? `${window.location.origin}/p/${link.slug}` : '');

	onMount(async () => {
		try {
			const { data: linkData } = await api.GET('/links/{id}', {
				params: { path: { id: linkId } }
			});

			if (!linkData) {
				goto('/links');
				return;
			}

			link = linkData;

			// Fetch slots
			const { data: slotsData } = await api.GET('/links/{id}/slots', {
				params: { path: { id: linkId } }
			});
			if (slotsData) {
				slots = slotsData;
			}

			// Fetch bookings or votes based on type
			if (link.type === 1) {
				const { data: bookingsData } = await api.GET('/links/{id}/bookings', {
					params: { path: { id: linkId } }
				});
				if (bookingsData) {
					bookings = bookingsData;
				}
			} else {
				const { data: votesData } = await api.GET('/links/{id}/votes', {
					params: { path: { id: linkId } }
				});
				if (votesData) {
					votes = votesData;
				}
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
		if (!confirm('Are you sure you want to delete this link? This action cannot be undone.')) {
			return;
		}

		deleting = true;
		try {
			await api.DELETE('/links/{id}', {
				params: { path: { id: linkId } }
			});
			goto('/links');
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
			<Button variant="secondary" onclick={() => goto(`/links/${linkId}/edit`)}>
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
						<span class="inline-flex items-center px-2.5 py-1 text-sm font-medium rounded-full {link!.type === 1 ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-indigo-50 text-indigo-700 border border-indigo-200'}">
							{typeLabels[link!.type]}
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

					{#if link!.description}
						<div class="col-span-2 md:col-span-4">
							<p class="text-sm font-medium text-gray-500">Description</p>
							<p class="text-gray-700">{link!.description}</p>
						</div>
					{/if}
				</div>
			{/snippet}
		</Card>

		<!-- Slots (for polls only) -->
		{#if link.type === 2}
			<SlotManager {linkId} bind:slots />
		{/if}

		<!-- Bookings or Votes -->
		{#if link.type === 1}
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
		{:else}
			<Card>
				{#snippet header()}
					<h2 class="text-lg font-medium text-gray-900">
						Votes ({votes.length})
					</h2>
				{/snippet}

				{#snippet children()}
					{#if votes.length === 0}
						<p class="text-gray-500 text-center py-4">No votes yet</p>
					{:else}
						<div class="divide-y divide-gray-100">
							{#each votes as vote (vote.id)}
								<VoteRow {vote} {slots} />
							{/each}
						</div>
					{/if}
				{/snippet}
			</Card>
		{/if}
	</div>
{/if}

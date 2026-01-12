<script lang="ts">
	import type { components } from '$lib/api/types';
	import { Badge, Button } from '$lib/components/ui';
	import { api } from '$lib/api/client';

	type Booking = components['schemas']['Booking'];
	type BookingStatus = components['schemas']['BookingStatus'];

	interface Props {
		booking: Booking;
		onUpdate?: (booking: Booking) => void;
	}

	let { booking, onUpdate }: Props = $props();

	let approving = $state(false);
	let declining = $state(false);

	const statusVariant: Record<BookingStatus, 'pending' | 'confirmed' | 'cancelled'> = {
		1: 'pending',
		2: 'confirmed',
		3: 'cancelled'
	};

	const formatDateTime = (dateStr: string) => {
		const date = new Date(dateStr);
		return date.toLocaleString(undefined, {
			dateStyle: 'medium',
			timeStyle: 'short'
		});
	};

	async function handleApprove() {
		approving = true;
		try {
			const { data } = await api.POST('/bookings/{id}/approve', {
				params: { path: { id: booking.id } }
			});
			if (data && onUpdate) {
				onUpdate(data);
			}
		} finally {
			approving = false;
		}
	}

	async function handleDecline() {
		declining = true;
		try {
			const { data } = await api.POST('/bookings/{id}/decline', {
				params: { path: { id: booking.id } }
			});
			if (data && onUpdate) {
				onUpdate(data);
			}
		} finally {
			declining = false;
		}
	}
</script>

<div class="flex items-center justify-between py-3 border-b border-gray-100 last:border-b-0">
	<div class="flex items-center gap-4 min-w-0 flex-1">
		<div class="min-w-0 flex-1">
			<p class="font-medium text-gray-900 truncate">
				{booking.guest_name || booking.guest_email}
			</p>
			{#if booking.guest_name}
				<p class="text-sm text-gray-500 truncate">{booking.guest_email}</p>
			{/if}
		</div>

		<div class="text-sm text-gray-600 whitespace-nowrap">
			{formatDateTime(booking.slot.start_time)}
		</div>

		<Badge variant={statusVariant[booking.status]} size="sm" />
	</div>

	{#if booking.status === 1}
		<div class="flex items-center gap-2 ml-4">
			<Button
				variant="secondary"
				size="sm"
				onclick={handleApprove}
				loading={approving}
				disabled={declining}
			>
				{#snippet children()}Approve{/snippet}
			</Button>
			<Button
				variant="danger"
				size="sm"
				onclick={handleDecline}
				loading={declining}
				disabled={approving}
			>
				{#snippet children()}Decline{/snippet}
			</Button>
		</div>
	{/if}
</div>

<script lang="ts">
	import type { components } from '$lib/api/types';
	import { IconButton } from '$lib/components/ui';
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

	const isPending = $derived(booking.status === 1);
	const isConfirmed = $derived(booking.status === 2);

	const getDurationMinutes = (startTime: string, endTime: string) => {
		const start = new Date(startTime);
		const end = new Date(endTime);
		return Math.round((end.getTime() - start.getTime()) / (1000 * 60));
	};

	const durationMinutes = $derived(getDurationMinutes(booking.slot.start_time, booking.slot.end_time));

	const formatDateTime = (dateStr: string) => {
		const date = new Date(dateStr);
		const now = new Date();
		const tomorrow = new Date(now);
		tomorrow.setDate(tomorrow.getDate() + 1);

		const isToday = date.toDateString() === now.toDateString();
		const isTomorrow = date.toDateString() === tomorrow.toDateString();

		const timeStr = date.toLocaleTimeString(undefined, {
			hour: 'numeric',
			minute: '2-digit'
		});

		if (isToday) {
			return `Today, ${timeStr}`;
		} else if (isTomorrow) {
			return `Tomorrow, ${timeStr}`;
		} else {
			const dateStr = date.toLocaleDateString(undefined, {
				month: 'short',
				day: 'numeric'
			});
			return `${dateStr}, ${timeStr}`;
		}
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

<div class="booking-item">
	<div class="booking-status" class:pending={isPending} class:confirmed={isConfirmed}></div>

	<div class="booking-info">
		<div class="booking-guest">{booking.guest_name || booking.guest_email}</div>
		<div class="booking-meta">{durationMinutes} min</div>
	</div>

	<div class="booking-time">{formatDateTime(booking.slot.start_time)}</div>

	<div class="booking-actions">
		{#if isPending}
			<IconButton
				variant="approve"
				title="Approve"
				onclick={handleApprove}
				disabled={approving || declining}
			>
				{#snippet children()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
					</svg>
				{/snippet}
			</IconButton>
			<IconButton
				variant="decline"
				title="Decline"
				onclick={handleDecline}
				disabled={approving || declining}
			>
				{#snippet children()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
					</svg>
				{/snippet}
			</IconButton>
		{:else if isConfirmed}
			<IconButton title="View">
				{#snippet children()}
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
						<path
							stroke-linecap="round"
							stroke-width="2"
							d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
						/>
					</svg>
				{/snippet}
			</IconButton>
		{/if}
	</div>
</div>

<style>
	.booking-item {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 1rem 1.25rem;
		border-bottom: var(--border-light);
		transition: background var(--transition);
	}

	.booking-item:last-child {
		border-bottom: none;
	}

	.booking-item:hover {
		background: var(--bg-tertiary);
	}

	.booking-status {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		border: 2px solid var(--border-color);
		flex-shrink: 0;
	}

	.booking-status.pending {
		background: var(--amber);
		animation: pulse 2s ease infinite;
	}

	.booking-status.confirmed {
		background: var(--emerald);
	}

	@keyframes pulse {
		0%,
		100% {
			transform: scale(1);
			opacity: 1;
		}
		50% {
			transform: scale(1.1);
			opacity: 0.7;
		}
	}

	.booking-info {
		flex: 1;
		min-width: 0;
	}

	.booking-guest {
		font-weight: 700;
		font-size: 0.9rem;
	}

	.booking-meta {
		font-size: 0.8rem;
		color: var(--text-muted);
		margin-top: 0.1rem;
	}

	.booking-time {
		font-size: 0.75rem;
		font-weight: 700;
		padding: 0.4rem 0.6rem;
		border-radius: var(--radius);
		background: var(--bg-tertiary);
		border: 1px solid var(--border-color);
		white-space: nowrap;
	}

	.booking-actions {
		display: flex;
		gap: 0.35rem;
	}
</style>

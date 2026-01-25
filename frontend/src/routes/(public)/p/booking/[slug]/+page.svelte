<script lang="ts">
	import { page } from '$app/state';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { Spinner, Card } from '$lib/components/ui';
	import { BookingPage } from '$lib/components/booking';

	type Slot = components['schemas']['Slot'];
	type CustomField = components['schemas']['CustomField'];

	interface PublicBookingLink {
		name: string;
		description?: string;
		custom_fields?: CustomField[];
		require_email?: boolean;
	}

	let bookingLink = $state<PublicBookingLink | null>(null);
	let slots = $state<Slot[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const slug = $derived(page.params.slug ?? '');

	$effect(() => {
		if (slug) {
			loadBookingLink();
		}
	});

	async function loadBookingLink() {
		if (!slug) return;

		loading = true;
		error = null;

		// Fetch booking link info
		const { data: linkData, error: linkError } = await api.GET('/p/booking/{slug}', {
			params: { path: { slug } }
		});

		if (linkError) {
			error = 'This booking page is not available.';
			loading = false;
			return;
		}

		bookingLink = linkData ?? null;

		// Fetch availability (slots)
		const now = new Date();
		const start = now.toISOString();
		const end = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000).toISOString(); // 30 days ahead

		const { data: availabilityData } = await api.GET('/p/booking/{slug}/availability', {
			params: { path: { slug }, query: { start, end } }
		});

		if (availabilityData?.slots) {
			slots = availabilityData.slots;
		}

		loading = false;
	}
</script>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if error}
	<Card class="text-center py-8">
		<div class="text-red-500 mb-4">
			<svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
			</svg>
		</div>
		<h1 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Page Not Found</h1>
		<p class="text-gray-600 dark:text-gray-400">{error}</p>
	</Card>
{:else if bookingLink}
	<BookingPage
		link={{
			type: 1,
			name: bookingLink.name,
			description: bookingLink.description,
			custom_fields: bookingLink.custom_fields,
			slots,
			require_email: bookingLink.require_email
		}}
		{slug}
	/>
{/if}

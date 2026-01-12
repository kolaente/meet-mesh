<script lang="ts">
	import type { components } from '$lib/api/types';
	import { Card, Badge } from '$lib/components/ui';

	type Link = components['schemas']['Link'];

	interface Props {
		link: Link;
	}

	let { link }: Props = $props();

	const typeLabels: Record<1 | 2, string> = {
		1: 'Booking',
		2: 'Poll'
	};

	const statusVariant = $derived<'active' | 'cancelled'>(link.status === 1 ? 'active' : 'cancelled');
	const typeVariant = $derived<'confirmed' | 'active'>(link.type === 1 ? 'confirmed' : 'active');
</script>

<a href="/links/{link.id}" class="block">
	<Card class="hover:shadow-md transition-shadow cursor-pointer">
		{#snippet children()}
			<div class="space-y-3">
				<div class="flex items-start justify-between gap-2">
					<h3 class="font-medium text-gray-900 truncate">{link.name}</h3>
					<span class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full {link.type === 1 ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-indigo-50 text-indigo-700 border border-indigo-200'}">
						{typeLabels[link.type]}
					</span>
				</div>

				<div class="flex items-center gap-2">
					<Badge variant={statusVariant} size="sm" />
				</div>

				<div class="text-sm text-gray-500 truncate">
					<span class="font-mono">/p/{link.slug}</span>
				</div>
			</div>
		{/snippet}
	</Card>
</a>

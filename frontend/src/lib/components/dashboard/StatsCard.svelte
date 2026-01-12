<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Card } from '$lib/components/ui';

	interface Props {
		label: string;
		value: number | string;
		icon?: Snippet;
	}

	let { label, value, icon }: Props = $props();
</script>

<Card class="flex-1">
	{#snippet children()}
		<div class="flex items-center gap-3">
			{#if icon}
				<div class="p-2 rounded-lg bg-indigo-50 text-indigo-600">
					{@render icon()}
				</div>
			{/if}
			<div>
				<p class="text-sm font-medium text-gray-500">{label}</p>
				<p class="text-2xl font-semibold text-gray-900">{value}</p>
			</div>
		</div>
	{/snippet}
</Card>

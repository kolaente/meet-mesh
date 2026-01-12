<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description: string;
		icon?: Snippet;
		action?: Snippet;
	}

	let { title, description, icon, action }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center py-12 px-6 text-center">
	{#if icon}
		<div class="text-slate-400 mb-4">
			{@render icon()}
		</div>
	{/if}

	<h3 class="text-lg font-medium text-slate-900 mb-2">{title}</h3>
	<p class="text-sm text-slate-500 max-w-sm mb-6">{description}</p>

	{#if action}
		<div class="mt-2">
			{@render action()}
		</div>
	{/if}
</div>

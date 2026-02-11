<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		class?: string;
		children: Snippet;
		header?: Snippet;
		footer?: Snippet;
	}

	let { class: className = '', children, header, footer }: Props = $props();
</script>

<div
	class="bg-bg-secondary rounded-brutalist-lg border-2 border-border shadow-brutalist overflow-hidden {className}"
>
	{#if header}
		<div class="flex items-center justify-between px-5 py-4 border-b-2 border-border bg-bg-tertiary">
			{@render header()}
		</div>
	{/if}

	<div class="px-5 py-4">
		{@render children()}
	</div>

	{#if footer}
		<div class="px-5 py-4 border-t-2 border-border bg-bg-tertiary rounded-b-brutalist-lg">
			{@render footer()}
		</div>
	{/if}
</div>

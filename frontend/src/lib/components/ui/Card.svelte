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
	class="bg-[var(--bg-secondary)] rounded-[var(--radius-lg)] border-[length:var(--border-width)] border-[var(--border-color)] shadow-[var(--shadow)] overflow-hidden {className}"
>
	{#if header}
		<div class="flex items-center justify-between px-5 py-4 border-b-[length:var(--border-width)] border-[var(--border-color)] bg-[var(--bg-tertiary)]">
			{@render header()}
		</div>
	{/if}

	<div class="px-5 py-4">
		{@render children()}
	</div>

	{#if footer}
		<div class="px-5 py-4 border-t-[length:var(--border-width)] border-[var(--border-color)] bg-[var(--bg-tertiary)] rounded-b-[var(--radius-lg)]">
			{@render footer()}
		</div>
	{/if}
</div>

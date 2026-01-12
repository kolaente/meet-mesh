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
	class="bg-white rounded-[var(--radius-lg)] border border-slate-200 shadow-[var(--shadow-sm)] {className}"
>
	{#if header}
		<div class="px-5 py-4 border-b border-slate-200">
			{@render header()}
		</div>
	{/if}

	<div class="px-5 py-4">
		{@render children()}
	</div>

	{#if footer}
		<div class="px-5 py-4 border-t border-slate-200 bg-slate-50 rounded-b-[var(--radius-lg)]">
			{@render footer()}
		</div>
	{/if}
</div>

<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Dialog } from 'bits-ui';

	interface Props {
		open?: boolean;
		title: string;
		children: Snippet;
		footer?: Snippet;
	}

	let { open = $bindable(false), title, children, footer }: Props = $props();
</script>

<Dialog.Root bind:open>
	<Dialog.Portal>
		<Dialog.Overlay
			class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0"
		/>
		<Dialog.Content
			class="fixed left-1/2 top-1/2 z-50 w-full max-w-lg -translate-x-1/2 -translate-y-1/2 rounded-[var(--radius-lg)] border border-gray-200 dark:border-neutral-700 bg-white dark:bg-neutral-900 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95"
		>
			<div class="flex items-center justify-between border-b border-gray-200 dark:border-neutral-700 px-6 py-4">
				<Dialog.Title class="text-lg font-semibold text-gray-900 dark:text-gray-100">
					{title}
				</Dialog.Title>
				<Dialog.Close
					class="rounded-[var(--radius-sm)] p-1 text-gray-400 dark:text-gray-500 hover:bg-gray-100 dark:hover:bg-neutral-800 hover:text-gray-500 dark:hover:text-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-neutral-900"
				>
					<svg
						class="h-5 w-5"
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
					>
						<path
							d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z"
						/>
					</svg>
					<span class="sr-only">Close</span>
				</Dialog.Close>
			</div>

			<div class="px-6 py-4">
				{@render children()}
			</div>

			{#if footer}
				<div class="flex justify-end gap-3 border-t border-gray-200 dark:border-neutral-700 bg-gray-50 dark:bg-neutral-800 px-6 py-4 rounded-b-[var(--radius-lg)]">
					{@render footer()}
				</div>
			{/if}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

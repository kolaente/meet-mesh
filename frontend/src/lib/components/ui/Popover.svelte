<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Popover } from 'bits-ui';

	interface Props {
		open?: boolean;
		trigger: Snippet;
		content: Snippet;
	}

	let { open = $bindable(false), trigger, content }: Props = $props();
</script>

<Popover.Root bind:open>
	<Popover.Trigger>
		{#snippet child({ props })}
			<span {...props}>
				{@render trigger()}
			</span>
		{/snippet}
	</Popover.Trigger>

	<Popover.Portal>
		<Popover.Content
			class="z-50 w-72 rounded-brutalist-md border border-gray-200 dark:border-neutral-700 bg-white dark:bg-neutral-900 p-4 shadow-md outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"
			sideOffset={4}
		>
			{@render content()}
			<Popover.Arrow class="fill-white dark:fill-neutral-900" />
		</Popover.Content>
	</Popover.Portal>
</Popover.Root>

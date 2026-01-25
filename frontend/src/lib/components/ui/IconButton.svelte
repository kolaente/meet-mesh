<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	interface Props extends HTMLButtonAttributes {
		variant?: 'default' | 'approve' | 'decline';
		children: Snippet;
	}

	let { variant = 'default', children, ...rest }: Props = $props();
</script>

<button class="icon-button {variant}" {...rest}>
	{@render children()}
</button>

<style>
	.icon-button {
		width: 32px;
		height: 32px;
		border-radius: var(--radius);
		border: var(--border);
		background: var(--bg-secondary);
		color: var(--text-secondary);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all var(--transition);
	}

	.icon-button :global(svg) {
		width: 14px;
		height: 14px;
	}

	.icon-button:hover {
		transform: translate(-1px, -1px);
		box-shadow: var(--shadow-sm);
	}

	.icon-button:active {
		transform: translate(1px, 1px);
		box-shadow: none;
	}

	.icon-button.approve:hover {
		background: var(--emerald);
		color: white;
		border-color: var(--border-color);
	}

	.icon-button.decline:hover {
		background: var(--rose);
		color: white;
		border-color: var(--border-color);
	}
</style>

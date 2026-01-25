<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import Spinner from './Spinner.svelte';

	type Variant = 'primary' | 'secondary' | 'ghost' | 'danger';
	type Size = 'sm' | 'md' | 'lg';

	interface Props extends HTMLButtonAttributes {
		variant?: Variant;
		size?: Size;
		loading?: boolean;
		children: Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		disabled = false,
		loading = false,
		type = 'button',
		onclick,
		children,
		class: className = '',
		...restProps
	}: Props = $props();

	const baseClasses =
		'inline-flex items-center justify-center font-bold transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50';

	const variantClasses: Record<Variant, string> = {
		primary: 'btn-primary',
		secondary: 'btn-secondary',
		ghost: 'btn-ghost',
		danger: 'btn-danger'
	};

	const sizeClasses: Record<Size, string> = {
		sm: 'h-8 px-3 text-sm rounded-[var(--radius)] gap-1.5',
		md: 'h-10 px-4 text-[0.85rem] rounded-[var(--radius)] gap-2',
		lg: 'h-12 px-6 text-base rounded-[var(--radius)] gap-2'
	};

	const spinnerSizes: Record<Size, 'sm' | 'md' | 'lg'> = {
		sm: 'sm',
		md: 'sm',
		lg: 'md'
	};
</script>

<button
	{type}
	{disabled}
	{onclick}
	class="{baseClasses} {variantClasses[variant]} {sizeClasses[size]} {className}"
	aria-disabled={disabled || loading}
	{...restProps}
>
	{#if loading}
		<Spinner size={spinnerSizes[size]} />
	{/if}
	{@render children()}
</button>

<style>
	/* Neo-brutalist button styles */
	.btn-primary {
		background-color: var(--sky);
		color: white;
		border: var(--border);
		box-shadow: var(--shadow);
	}

	.btn-primary:hover:not(:disabled) {
		transform: translate(-2px, -2px);
		box-shadow: var(--shadow-hover);
	}

	.btn-primary:active:not(:disabled) {
		transform: translate(2px, 2px);
		box-shadow: none;
	}

	.btn-secondary {
		background-color: var(--bg-secondary);
		color: var(--text-primary);
		border: var(--border);
		box-shadow: var(--shadow-sm);
	}

	.btn-secondary:hover:not(:disabled) {
		background-color: var(--violet);
		color: white;
		transform: translate(-1px, -1px);
		box-shadow: var(--shadow);
	}

	.btn-secondary:active:not(:disabled) {
		transform: translate(1px, 1px);
		box-shadow: none;
	}

	.btn-ghost {
		background-color: transparent;
		color: var(--text-secondary);
		border: 2px solid transparent;
	}

	.btn-ghost:hover:not(:disabled) {
		background-color: var(--bg-tertiary);
		color: var(--text-primary);
		border-color: var(--border-color);
	}

	.btn-ghost:active:not(:disabled) {
		background-color: var(--bg-tertiary);
	}

	.btn-danger {
		background-color: var(--rose);
		color: white;
		border: var(--border);
		box-shadow: var(--shadow-sm);
	}

	.btn-danger:hover:not(:disabled) {
		transform: translate(-1px, -1px);
		box-shadow: var(--shadow);
	}

	.btn-danger:active:not(:disabled) {
		transform: translate(1px, 1px);
		box-shadow: none;
	}
</style>

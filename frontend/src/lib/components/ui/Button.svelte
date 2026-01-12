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
		'inline-flex items-center justify-center font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50';

	const variantClasses: Record<Variant, string> = {
		primary:
			'bg-indigo-600 text-white hover:bg-indigo-700 shadow-sm active:bg-indigo-800',
		secondary:
			'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 active:bg-gray-100',
		ghost: 'text-gray-700 hover:bg-gray-100 active:bg-gray-200',
		danger:
			'bg-red-50 text-red-700 hover:bg-red-100 active:bg-red-200 border border-red-200'
	};

	const sizeClasses: Record<Size, string> = {
		sm: 'h-8 px-3 text-sm rounded-[var(--radius-sm)] gap-1.5',
		md: 'h-10 px-4 text-sm rounded-[var(--radius-md)] gap-2',
		lg: 'h-12 px-6 text-base rounded-[var(--radius-md)] gap-2'
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

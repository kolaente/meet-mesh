<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends Omit<HTMLInputAttributes, 'value'> {
		label?: string;
		name: string;
		type?: HTMLInputAttributes['type'];
		value?: string;
		placeholder?: string;
		error?: string;
		description?: string;
		required?: boolean;
		disabled?: boolean;
		inputmode?: 'none' | 'text' | 'tel' | 'url' | 'email' | 'numeric' | 'decimal' | 'search';
	}

	let {
		label,
		name,
		type = 'text',
		value = $bindable(''),
		placeholder,
		error,
		description,
		required = false,
		disabled = false,
		inputmode,
		class: className = '',
		...restProps
	}: Props = $props();

	let inputId = $derived(`input-${name}`);
	let descriptionId = $derived(description ? `${inputId}-description` : undefined);
	let errorId = $derived(error ? `${inputId}-error` : undefined);

	// Auto-detect inputmode from type if not explicitly set
	let effectiveInputmode = $derived(
		inputmode ?? (type === 'email' ? 'email' : type === 'tel' ? 'tel' : type === 'url' ? 'url' : undefined)
	);
</script>

<div class="space-y-1.5 {className}">
	{#if label}
		<label for={inputId} class="block text-sm font-medium text-text-secondary">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	{#if description}
		<p id={descriptionId} class="text-sm text-text-muted">
			{description}
		</p>
	{/if}

	<input
		id={inputId}
		{name}
		{type}
		bind:value
		{placeholder}
		{required}
		{disabled}
		inputmode={effectiveInputmode}
		aria-invalid={!!error}
		aria-describedby={[descriptionId, errorId].filter(Boolean).join(' ') || undefined}
		class="block w-full min-h-[44px] rounded-brutalist border-2 border-border bg-bg-secondary px-3 py-2.5 sm:py-2 text-base sm:text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-accent-sky disabled:cursor-not-allowed disabled:opacity-50
		{error ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
		{...restProps}
	/>

	{#if error}
		<p id={errorId} class="text-sm text-red-600">
			{error}
		</p>
	{/if}
</div>

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
		class: className = '',
		...restProps
	}: Props = $props();

	let inputId = $derived(`input-${name}`);
	let descriptionId = $derived(description ? `${inputId}-description` : undefined);
	let errorId = $derived(error ? `${inputId}-error` : undefined);
</script>

<div class="space-y-1.5 {className}">
	{#if label}
		<label for={inputId} class="block text-sm font-medium text-gray-700">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	{#if description}
		<p id={descriptionId} class="text-sm text-gray-500">
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
		aria-invalid={!!error}
		aria-describedby={[descriptionId, errorId].filter(Boolean).join(' ') || undefined}
		class="block w-full rounded-[var(--radius-md)] border px-3 py-2 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-0 disabled:cursor-not-allowed disabled:bg-gray-50 disabled:text-gray-500
		{error
			? 'border-red-300 focus:border-red-500 focus:ring-red-500'
			: 'border-gray-300 focus:border-indigo-500 focus:ring-indigo-500'}"
		{...restProps}
	/>

	{#if error}
		<p id={errorId} class="text-sm text-red-600">
			{error}
		</p>
	{/if}
</div>

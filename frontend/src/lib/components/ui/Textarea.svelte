<script lang="ts">
	import type { HTMLTextareaAttributes } from 'svelte/elements';

	interface Props extends Omit<HTMLTextareaAttributes, 'value'> {
		label?: string;
		name: string;
		value?: string;
		placeholder?: string;
		error?: string;
		description?: string;
		required?: boolean;
		disabled?: boolean;
		rows?: number;
	}

	let {
		label,
		name,
		value = $bindable(''),
		placeholder,
		error,
		description,
		required = false,
		disabled = false,
		rows = 4,
		class: className = '',
		...restProps
	}: Props = $props();

	let textareaId = $derived(`textarea-${name}`);
	let descriptionId = $derived(description ? `${textareaId}-description` : undefined);
	let errorId = $derived(error ? `${textareaId}-error` : undefined);
</script>

<div class="space-y-1.5 {className}">
	{#if label}
		<label for={textareaId} class="block text-sm font-medium text-text-secondary">
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

	<textarea
		id={textareaId}
		{name}
		bind:value
		{placeholder}
		{required}
		{disabled}
		{rows}
		aria-invalid={!!error}
		aria-describedby={[descriptionId, errorId].filter(Boolean).join(' ') || undefined}
		class="block w-full min-h-[100px] rounded-brutalist border-2 border-border bg-bg-secondary px-3 py-2.5 sm:py-2 text-base sm:text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:ring-2 focus:ring-offset-0 focus:ring-accent-sky disabled:cursor-not-allowed disabled:opacity-50
		{error ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
		{...restProps}
	></textarea>

	{#if error}
		<p id={errorId} class="text-sm text-red-600">
			{error}
		</p>
	{/if}
</div>

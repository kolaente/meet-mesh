<script lang="ts">
	import Input from '../ui/Input.svelte';
	import Select from '../ui/Select.svelte';
	import Button from '../ui/Button.svelte';
	import type { components } from '$lib/api/types';

	type CustomField = components['schemas']['CustomField'];

	interface FormData {
		email: string;
		name?: string;
		customFields: Record<string, string>;
	}

	interface Props {
		customFields?: CustomField[];
		onSubmit: (data: FormData) => void;
		loading?: boolean;
		class?: string;
	}

	let {
		customFields = [],
		onSubmit,
		loading = false,
		class: className = ''
	}: Props = $props();

	let email = $state('');
	let name = $state('');
	let customFieldValues = $state<Record<string, string>>({});

	// Initialize custom field values
	$effect(() => {
		const values: Record<string, string> = {};
		for (const field of customFields) {
			if (!(field.name in customFieldValues)) {
				values[field.name] = '';
			}
		}
		if (Object.keys(values).length > 0) {
			customFieldValues = { ...customFieldValues, ...values };
		}
	});

	function handleSubmit(event: Event) {
		event.preventDefault();
		onSubmit({
			email,
			name: name || undefined,
			customFields: customFieldValues
		});
	}

	function getFieldType(field: CustomField): 'text' | 'email' | 'tel' {
		// CustomFieldType: 1=text, 2=email, 3=phone, 4=select, 5=textarea
		switch (field.type) {
			case 2:
				return 'email';
			case 3:
				return 'tel';
			default:
				return 'text';
		}
	}

	function isTextarea(field: CustomField): boolean {
		return field.type === 5;
	}

	function isSelect(field: CustomField): boolean {
		return field.type === 4;
	}
</script>

<form onsubmit={handleSubmit} class="space-y-4 sm:space-y-5 {className}">
	<!-- Email field (always present) -->
	<Input
		name="email"
		label="Email"
		type="email"
		inputmode="email"
		autocomplete="email"
		bind:value={email}
		placeholder="you@example.com"
		required
	/>

	<!-- Name field -->
	<Input
		name="name"
		label="Name"
		type="text"
		autocomplete="name"
		bind:value={name}
		placeholder="Your name"
	/>

	<!-- Custom fields -->
	{#each customFields as field (field.name)}
		{#if isSelect(field)}
			<Select
				name={field.name}
				label={field.label}
				options={(field.options ?? []).map((opt) => ({ value: opt, label: opt }))}
				bind:value={customFieldValues[field.name]}
				placeholder="Select an option"
			/>
		{:else if isTextarea(field)}
			<div class="space-y-1.5">
				<label for={`field-${field.name}`} class="block text-sm font-medium text-gray-700">
					{field.label}
					{#if field.required}
						<span class="text-red-500">*</span>
					{/if}
				</label>
				<textarea
					id={`field-${field.name}`}
					name={field.name}
					bind:value={customFieldValues[field.name]}
					required={field.required}
					rows="3"
					class="block w-full min-h-[100px] rounded-[var(--radius-md)] border border-gray-300 px-3 py-3 text-base sm:text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-0 focus:border-indigo-500"
				></textarea>
			</div>
		{:else}
			<Input
				name={field.name}
				label={field.label}
				type={getFieldType(field)}
				inputmode={getFieldType(field) === 'tel' ? 'tel' : getFieldType(field) === 'email' ? 'email' : 'text'}
				bind:value={customFieldValues[field.name]}
				required={field.required}
			/>
		{/if}
	{/each}

	<div class="pt-3 sm:pt-2">
		<Button type="submit" {loading} class="w-full min-h-[48px]">
			{loading ? 'Booking...' : 'Confirm Booking'}
		</Button>
	</div>
</form>

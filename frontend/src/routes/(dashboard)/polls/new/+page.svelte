<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader } from '$lib/components/dashboard';
	import { Button, Card, Checkbox, Input, Textarea } from '$lib/components/ui';

	type CustomField = components['schemas']['CustomField'];
	type CustomFieldType = components['schemas']['CustomFieldType'];

	let name = $state('');
	let description = $state('');
	let showResults = $state(false);
	let requireEmail = $state(false);
	let customFields = $state<CustomField[]>([]);
	let saving = $state(false);
	let error = $state('');

	const fieldTypeOptions = [
		{ value: '1', label: 'Text' },
		{ value: '2', label: 'Email' },
		{ value: '3', label: 'Phone' },
		{ value: '4', label: 'Select' },
		{ value: '5', label: 'Textarea' }
	];

	function addCustomField() {
		customFields = [
			...customFields,
			{
				name: '',
				label: '',
				type: 1 as CustomFieldType,
				required: false,
				options: []
			}
		];
	}

	function removeCustomField(index: number) {
		customFields = customFields.filter((_, i) => i !== index);
	}

	function updateField(index: number, key: keyof CustomField, value: unknown) {
		customFields = customFields.map((field, i) => {
			if (i !== index) return field;
			return { ...field, [key]: value };
		});
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		saving = true;

		try {
			const { data, error: apiError } = await api.POST('/polls', {
				body: {
					name,
					description: description || undefined,
					show_results: showResults,
					require_email: requireEmail,
					custom_fields: customFields.length > 0 ? customFields : undefined
				}
			});

			if (apiError) {
				error = 'Failed to create poll';
				return;
			}

			if (data) {
				goto(`/polls/${data.id}`);
			}
		} catch {
			error = 'An unexpected error occurred';
		} finally {
			saving = false;
		}
	}
</script>

<DashboardHeader title="Create Poll" />

<Card class="max-w-2xl">
	{#snippet children()}
		<form onsubmit={handleSubmit} class="space-y-6">
			{#if error}
				<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
					{error}
				</div>
			{/if}

			<Input label="Name" name="name" bind:value={name} required placeholder="Team Meeting Poll" />

			<Textarea
				label="Description"
				name="description"
				bind:value={description}
				rows={3}
				placeholder="A brief description of this poll..."
			/>

			<!-- Poll-specific options -->
			<div class="space-y-4">
				<Checkbox
					bind:checked={showResults}
					label="Show results to voters"
				/>

				<Checkbox
					bind:checked={requireEmail}
					label="Require email from voters"
				/>
			</div>

			<!-- Custom Fields -->
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<span class="block text-sm font-medium text-gray-700">Custom Fields</span>
					<Button variant="secondary" size="sm" type="button" onclick={addCustomField}>
						{#snippet children()}Add Field{/snippet}
					</Button>
				</div>

				{#if customFields.length > 0}
					<div class="space-y-4">
						{#each customFields as field, index (index)}
							<div class="p-4 bg-gray-50 rounded-lg space-y-3">
								<div class="flex items-start justify-between gap-4">
									<div class="flex-1 grid grid-cols-2 gap-3">
										<Input
											label="Field Name"
											name={`field-name-${index}`}
											value={field.name}
											oninput={(e: Event) =>
												updateField(index, 'name', (e.target as HTMLInputElement).value)}
											placeholder="company"
											required
										/>
										<Input
											label="Label"
											name={`field-label-${index}`}
											value={field.label}
											oninput={(e: Event) =>
												updateField(index, 'label', (e.target as HTMLInputElement).value)}
											placeholder="Company Name"
											required
										/>
									</div>
									<button
										type="button"
										onclick={() => removeCustomField(index)}
										class="mt-6 p-1 text-gray-400 hover:text-red-500"
										aria-label="Remove field"
									>
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M6 18L18 6M6 6l12 12"
											/>
										</svg>
									</button>
								</div>

								<div class="grid grid-cols-2 gap-3">
									<div class="space-y-1.5">
										<label for={`field-type-${index}`} class="block text-sm font-medium text-gray-700">Field Type</label>
										<select
											id={`field-type-${index}`}
											name={`field-type-${index}`}
											value={String(field.type)}
											onchange={(e: Event) =>
												updateField(index, 'type', Number((e.target as HTMLSelectElement).value) as CustomFieldType)}
											class="block w-full h-10 rounded-[var(--radius-md)] border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:border-indigo-500 focus:ring-indigo-500"
										>
											{#each fieldTypeOptions as opt (opt.value)}
												<option value={opt.value}>{opt.label}</option>
											{/each}
										</select>
									</div>
									<div class="flex items-end pb-1">
										<Checkbox
											checked={field.required}
											onchange={(checked) => updateField(index, 'required', checked)}
											label="Required"
										/>
									</div>
								</div>

								{#if field.type === 4}
									<Input
										label="Options (comma-separated)"
										name={`field-options-${index}`}
										value={field.options?.join(', ') || ''}
										oninput={(e: Event) =>
											updateField(
												index,
												'options',
												(e.target as HTMLInputElement).value.split(',').map((s) => s.trim())
											)}
										placeholder="Option 1, Option 2, Option 3"
									/>
								{/if}
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-gray-500">No custom fields. Add fields to collect additional information from voters.</p>
				{/if}
			</div>

			<div class="flex justify-end gap-3 pt-4">
				<Button variant="secondary" type="button" onclick={() => goto('/polls')}>
					{#snippet children()}Cancel{/snippet}
				</Button>
				<Button variant="primary" type="submit" loading={saving}>
					{#snippet children()}Create Poll{/snippet}
				</Button>
			</div>
		</form>
	{/snippet}
</Card>

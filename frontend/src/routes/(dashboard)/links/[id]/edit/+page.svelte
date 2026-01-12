<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { DashboardHeader } from '$lib/components/dashboard';
	import { Button, Card, Input, Select, Spinner } from '$lib/components/ui';

	type Link = components['schemas']['Link'];
	type CustomField = components['schemas']['CustomField'];
	type CustomFieldType = components['schemas']['CustomFieldType'];
	type LinkStatus = components['schemas']['LinkStatus'];

	let link = $state<Link | null>(null);
	let name = $state('');
	let status = $state<'1' | '2'>('1');
	let description = $state('');
	let customFields = $state<CustomField[]>([]);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state('');

	const linkId = $derived(Number($page.params.id));

	const statusOptions = [
		{ value: '1', label: 'Active' },
		{ value: '2', label: 'Closed' }
	];

	const fieldTypeOptions = [
		{ value: '1', label: 'Text' },
		{ value: '2', label: 'Email' },
		{ value: '3', label: 'Phone' },
		{ value: '4', label: 'Select' },
		{ value: '5', label: 'Textarea' }
	];

	onMount(async () => {
		try {
			const { data } = await api.GET('/links/{id}', {
				params: { path: { id: linkId } }
			});

			if (!data) {
				goto('/links');
				return;
			}

			link = data;
			name = data.name;
			status = String(data.status) as '1' | '2';
			description = data.description || '';
			customFields = data.custom_fields || [];
		} finally {
			loading = false;
		}
	});

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
			const { data, error: apiError } = await api.PUT('/links/{id}', {
				params: { path: { id: linkId } },
				body: {
					name,
					status: Number(status) as LinkStatus,
					description: description || undefined,
					custom_fields: customFields.length > 0 ? customFields : undefined
				}
			});

			if (apiError) {
				error = 'Failed to update link';
				return;
			}

			if (data) {
				goto(`/links/${linkId}`);
			}
		} catch {
			error = 'An unexpected error occurred';
		} finally {
			saving = false;
		}
	}
</script>

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else if link}
	<DashboardHeader title="Edit Link" />

	<Card class="max-w-2xl">
		{#snippet children()}
			<form onsubmit={handleSubmit} class="space-y-6">
				{#if error}
					<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
						{error}
					</div>
				{/if}

				<Input label="Name" name="name" bind:value={name} required placeholder="30 Minute Meeting" />

				<Select label="Status" name="status" options={statusOptions} bind:value={status} />

				<div class="space-y-1.5">
					<label for="edit-description" class="block text-sm font-medium text-gray-700">Description</label>
					<textarea
						id="edit-description"
						name="description"
						bind:value={description}
						rows="3"
						class="block w-full rounded-[var(--radius-md)] border border-gray-300 px-3 py-2 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:border-indigo-500 focus:ring-indigo-500"
						placeholder="A brief description of this scheduling link..."
					></textarea>
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
											<label for={`edit-field-type-${index}`} class="block text-sm font-medium text-gray-700">Field Type</label>
											<select
												id={`edit-field-type-${index}`}
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
											<label class="flex items-center gap-2 text-sm text-gray-700">
												<input
													type="checkbox"
													checked={field.required}
													onchange={(e: Event) =>
														updateField(index, 'required', (e.target as HTMLInputElement).checked)}
													class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
												/>
												Required
											</label>
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
						<p class="text-sm text-gray-500">
							No custom fields. Add fields to collect additional information from guests.
						</p>
					{/if}
				</div>

				<div class="flex justify-end gap-3 pt-4">
					<Button variant="secondary" type="button" onclick={() => goto(`/links/${linkId}`)}>
						{#snippet children()}Cancel{/snippet}
					</Button>
					<Button variant="primary" type="submit" loading={saving}>
						{#snippet children()}Save Changes{/snippet}
					</Button>
				</div>
			</form>
		{/snippet}
	</Card>
{/if}

<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib/api/client';
	import type { components } from '$lib/api/types';
	import { AvailabilityEditor } from '$lib/components/booking';
	import { DashboardHeader } from '$lib/components/dashboard';
	import { Button, Card, Checkbox, Input, Select, Textarea } from '$lib/components/ui';

	type AvailabilityRule = components['schemas']['AvailabilityRule'];
	type CustomField = components['schemas']['CustomField'];
	type CustomFieldType = components['schemas']['CustomFieldType'];

	let name = $state('');
	let description = $state('');
	let meetingLink = $state('');
	let slotDurationsMinutes = $state<number[]>([30]);
	let bufferMinutes = $state('0');
	let autoConfirm = $state(true);
	let customFields = $state<CustomField[]>([]);
	let availabilityRules = $state<AvailabilityRule[]>([
		{ days_of_week: [1, 2, 3, 4, 5], start_time: '09:00', end_time: '17:00' }
	]);
	let saving = $state(false);
	let error = $state('');

	const fieldTypeOptions = [
		{ value: '1', label: 'Text' },
		{ value: '2', label: 'Email' },
		{ value: '3', label: 'Phone' },
		{ value: '4', label: 'Select' },
		{ value: '5', label: 'Textarea' }
	];

	const slotDurationOptions = [
		{ value: 15, label: '15 min' },
		{ value: 30, label: '30 min' },
		{ value: 45, label: '45 min' },
		{ value: 60, label: '60 min' },
		{ value: 90, label: '90 min' },
		{ value: 120, label: '120 min' }
	];

	function toggleDuration(duration: number) {
		if (slotDurationsMinutes.includes(duration)) {
			// Only allow removal if more than one duration is selected
			if (slotDurationsMinutes.length > 1) {
				slotDurationsMinutes = slotDurationsMinutes.filter((d) => d !== duration);
			}
		} else {
			slotDurationsMinutes = [...slotDurationsMinutes, duration].sort((a, b) => a - b);
		}
	}

	const bufferTimeOptions = [
		{ value: '0', label: 'No buffer' },
		{ value: '5', label: '5 minutes' },
		{ value: '10', label: '10 minutes' },
		{ value: '15', label: '15 minutes' },
		{ value: '30', label: '30 minutes' }
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
			const { data, error: apiError } = await api.POST('/booking-links', {
				body: {
					name,
					description: description || undefined,
					meeting_link: meetingLink || undefined,
					slot_durations_minutes: slotDurationsMinutes,
					buffer_minutes: Number(bufferMinutes),
					auto_confirm: autoConfirm,
					availability_rules: availabilityRules.length > 0 ? availabilityRules : undefined,
					custom_fields: customFields.length > 0 ? customFields : undefined
				}
			});

			if (apiError) {
				error = 'Failed to create booking link';
				return;
			}

			if (data) {
				goto(`/booking-links/${data.id}`);
			}
		} catch {
			error = 'An unexpected error occurred';
		} finally {
			saving = false;
		}
	}
</script>

<DashboardHeader title="Create Booking Link" />

<Card class="max-w-2xl">
	{#snippet children()}
		<form onsubmit={handleSubmit} class="space-y-6">
			{#if error}
				<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
					{error}
				</div>
			{/if}

			<Input label="Name" name="name" bind:value={name} required placeholder="30 Minute Meeting" />

			<Textarea
				label="Description"
				name="description"
				bind:value={description}
				rows={3}
				placeholder="A brief description of this booking link..."
			/>

			<Input
				label="Meeting Link"
				name="meetingLink"
				type="url"
				bind:value={meetingLink}
				placeholder="https://zoom.us/j/... or https://meet.google.com/..."
				description="Video conference link (Zoom, Google Meet, etc.) to include in calendar events"
			/>

			<!-- Slot Durations -->
			<div class="space-y-2">
				<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
					Slot Durations
				</label>
				<p class="text-sm text-gray-500 dark:text-gray-400">
					Select which duration options guests can choose from
				</p>
				<div class="flex flex-wrap gap-3 pt-1">
					{#each slotDurationOptions as option (option.value)}
						<label
							class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border cursor-pointer transition-colors
								{slotDurationsMinutes.includes(option.value)
									? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-700 dark:text-indigo-300'
									: 'bg-white dark:bg-neutral-800 border-gray-300 dark:border-neutral-600 text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-neutral-500'}"
						>
							<input
								type="checkbox"
								checked={slotDurationsMinutes.includes(option.value)}
								onchange={() => toggleDuration(option.value)}
								class="sr-only"
							/>
							<span class="text-sm font-medium">{option.label}</span>
							{#if slotDurationsMinutes.includes(option.value)}
								<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
								</svg>
							{/if}
						</label>
					{/each}
				</div>
			</div>

			<!-- Buffer Time -->
			<Select
				label="Buffer Time"
				name="bufferMinutes"
				options={bufferTimeOptions}
				bind:value={bufferMinutes}
			/>

			<!-- Confirmation Setting -->
			<div class="pt-2">
				<Checkbox
					bind:checked={autoConfirm}
					label="Auto-confirm bookings"
					description="When disabled, bookings require your approval before guests receive confirmation"
				/>
			</div>

			<!-- Availability -->
			<AvailabilityEditor
				rules={availabilityRules}
				onchange={(rules) => availabilityRules = rules}
			/>

			<!-- Custom Fields -->
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<span class="block text-sm font-medium text-gray-700 dark:text-gray-300">Custom Fields</span>
					<Button variant="secondary" size="sm" type="button" onclick={addCustomField}>
						{#snippet children()}Add Field{/snippet}
					</Button>
				</div>

				{#if customFields.length > 0}
					<div class="space-y-4">
						{#each customFields as field, index (index)}
							<div class="p-4 bg-gray-50 dark:bg-neutral-800 rounded-lg space-y-3">
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
										class="mt-6 p-1 text-gray-400 dark:text-gray-500 hover:text-red-500 dark:hover:text-red-400"
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
										<label for={`field-type-${index}`} class="block text-sm font-medium text-gray-700 dark:text-gray-300">Field Type</label>
										<select
											id={`field-type-${index}`}
											name={`field-type-${index}`}
											value={String(field.type)}
											onchange={(e: Event) =>
												updateField(index, 'type', Number((e.target as HTMLSelectElement).value) as CustomFieldType)}
											class="block w-full h-10 rounded-brutalist-md border border-gray-300 dark:border-neutral-600 bg-white dark:bg-neutral-700 px-3 py-2 text-sm text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-0 focus:border-indigo-500 focus:ring-indigo-500"
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
					<p class="text-sm text-gray-500 dark:text-gray-400">No custom fields. Add fields to collect additional information from guests.</p>
				{/if}
			</div>

			<div class="flex justify-end gap-3 pt-4">
				<Button variant="secondary" type="button" onclick={() => goto('/booking-links')}>
					{#snippet children()}Cancel{/snippet}
				</Button>
				<Button variant="primary" type="submit" loading={saving}>
					{#snippet children()}Create Booking Link{/snippet}
				</Button>
			</div>
		</form>
	{/snippet}
</Card>

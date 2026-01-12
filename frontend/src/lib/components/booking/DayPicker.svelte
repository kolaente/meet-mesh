<script lang="ts">
	import { Select } from 'bits-ui';

	interface Props {
		availableDates: string[];
		selectedDate?: string;
		class?: string;
	}

	let {
		availableDates,
		selectedDate = $bindable(),
		class: className = ''
	}: Props = $props();

	// Convert available dates to options
	let dateOptions = $derived(
		availableDates.map((dateStr) => {
			const date = new Date(dateStr);
			return {
				value: dateStr,
				label: date.toLocaleDateString('en-US', {
					weekday: 'long',
					year: 'numeric',
					month: 'long',
					day: 'numeric'
				})
			};
		})
	);

	let selectedOption = $derived(dateOptions.find((opt) => opt.value === selectedDate));
</script>

<div class="space-y-1.5 {className}">
	<span id="day-picker-label" class="block text-sm font-medium text-gray-700">Select Date</span>

	<Select.Root type="single" bind:value={selectedDate}>
		<Select.Trigger
			aria-labelledby="day-picker-label"
			class="flex h-12 w-full items-center justify-between rounded-[var(--radius-md)] border border-gray-300 bg-white px-4 py-2 text-base focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-0 focus:border-indigo-500"
		>
			<span class={selectedOption ? 'text-gray-900' : 'text-gray-400'}>
				{selectedOption?.label ?? 'Choose a date'}
			</span>
			<svg
				class="h-5 w-5 text-gray-400"
				xmlns="http://www.w3.org/2000/svg"
				viewBox="0 0 20 20"
				fill="currentColor"
				aria-hidden="true"
			>
				<path
					fill-rule="evenodd"
					d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z"
					clip-rule="evenodd"
				/>
			</svg>
		</Select.Trigger>

		<Select.Portal>
			<Select.Content
				class="z-50 max-h-60 min-w-[8rem] overflow-auto rounded-[var(--radius-md)] border border-gray-200 bg-white shadow-md animate-in fade-in-0 zoom-in-95"
				sideOffset={4}
			>
				<Select.Viewport class="p-1">
					{#each dateOptions as option (option.value)}
						<Select.Item
							value={option.value}
							label={option.label}
							class="relative flex cursor-pointer select-none items-center rounded-[var(--radius-sm)] py-3 pl-10 pr-4 text-sm text-gray-900 outline-none transition-colors hover:bg-gray-100 focus:bg-gray-100 data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
						>
							{#snippet children({ selected })}
								{#if selected}
									<span class="absolute left-3 flex h-4 w-4 items-center justify-center">
										<svg
											class="h-4 w-4 text-indigo-600"
											xmlns="http://www.w3.org/2000/svg"
											viewBox="0 0 20 20"
											fill="currentColor"
										>
											<path
												fill-rule="evenodd"
												d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z"
												clip-rule="evenodd"
											/>
										</svg>
									</span>
								{/if}
								{option.label}
							{/snippet}
						</Select.Item>
					{/each}
				</Select.Viewport>
			</Select.Content>
		</Select.Portal>
	</Select.Root>

	{#if selectedDate}
		<p class="text-sm text-gray-600 mt-2">
			Selected: <span class="font-medium">{selectedOption?.label}</span>
		</p>
	{/if}
</div>

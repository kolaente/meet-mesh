<script lang="ts">
	import { Select } from 'bits-ui';

	interface Option {
		value: string;
		label: string;
	}

	interface Props {
		label?: string;
		name: string;
		options: Option[];
		value?: string;
		placeholder?: string;
		error?: string;
		disabled?: boolean;
	}

	let {
		label,
		name,
		options,
		value = $bindable(''),
		placeholder = 'Select an option',
		error,
		disabled = false
	}: Props = $props();

	let selectId = $derived(`select-${name}`);
	let errorId = $derived(error ? `${selectId}-error` : undefined);
	let selectedOption = $derived(options.find((opt) => opt.value === value));
</script>

<div class="space-y-1.5">
	{#if label}
		<label for={selectId} class="block text-sm font-medium text-gray-700">
			{label}
		</label>
	{/if}

	<Select.Root
		type="single"
		{name}
		{disabled}
		bind:value
	>
		<Select.Trigger
			id={selectId}
			aria-invalid={!!error}
			aria-describedby={errorId}
			class="flex min-h-[44px] h-auto sm:h-10 w-full items-center justify-between rounded-[var(--radius-md)] border bg-white px-3 py-2.5 sm:py-2 text-base sm:text-sm focus:outline-none focus:ring-2 focus:ring-offset-0 disabled:cursor-not-allowed disabled:bg-gray-50 disabled:text-gray-500
			{error
				? 'border-red-300 focus:border-red-500 focus:ring-red-500'
				: 'border-gray-300 focus:border-indigo-500 focus:ring-indigo-500'}"
		>
			<span class={selectedOption ? 'text-gray-900' : 'text-gray-400'}>
				{selectedOption?.label ?? placeholder}
			</span>
			<svg
				class="h-4 w-4 text-gray-400 flex-shrink-0"
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
				class="z-50 min-w-[8rem] max-h-[50vh] overflow-hidden rounded-[var(--radius-md)] border border-gray-200 bg-white shadow-md animate-in fade-in-0 zoom-in-95"
				sideOffset={4}
			>
				<Select.Viewport class="p-1 overflow-y-auto">
					{#each options as option (option.value)}
						<Select.Item
							value={option.value}
							label={option.label}
							class="relative flex cursor-pointer select-none items-center rounded-[var(--radius-sm)] py-3 sm:py-2 pl-8 pr-3 text-base sm:text-sm text-gray-900 outline-none transition-colors hover:bg-gray-100 focus:bg-gray-100 active:bg-gray-200 data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
						>
							{#snippet children({ selected })}
								{#if selected}
									<span class="absolute left-2 flex h-4 w-4 items-center justify-center">
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

	{#if error}
		<p id={errorId} class="text-sm text-red-600">
			{error}
		</p>
	{/if}
</div>

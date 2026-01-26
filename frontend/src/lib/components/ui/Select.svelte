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
		onchange?: (value: string) => void;
	}

	let {
		label,
		name,
		options,
		value = $bindable(''),
		placeholder = 'Select an option',
		error,
		disabled = false,
		onchange
	}: Props = $props();

	let selectId = $derived(`select-${name}`);
	let errorId = $derived(error ? `${selectId}-error` : undefined);
	let selectedOption = $derived(options.find((opt) => opt.value === value));
</script>

<div class="space-y-1.5">
	{#if label}
		<label for={selectId} class="block text-sm font-medium text-[var(--text-secondary)]">
			{label}
		</label>
	{/if}

	<Select.Root
		type="single"
		{name}
		{disabled}
		bind:value
		onValueChange={(v) => onchange?.(v)}
	>
		<Select.Trigger
			id={selectId}
			aria-invalid={!!error}
			aria-describedby={errorId}
			class="flex min-h-[44px] h-auto sm:h-10 w-full items-center justify-between rounded-[var(--radius)] border-2 border-[var(--border-color)] bg-[var(--bg-secondary)] px-3 py-2.5 sm:py-2 text-base sm:text-sm text-[var(--text-primary)] focus:outline-none focus:ring-2 focus:ring-[var(--sky)] focus:ring-offset-0 disabled:cursor-not-allowed disabled:opacity-50
			{error ? 'border-red-500 focus:ring-red-500' : ''}"
		>
			<span class={selectedOption ? 'text-[var(--text-primary)]' : 'text-[var(--text-primary)] opacity-50'}>
				{selectedOption?.label ?? placeholder}
			</span>
			<svg
				class="h-4 w-4 text-[var(--text-primary)] opacity-50 flex-shrink-0"
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
				class="z-50 min-w-[8rem] max-h-[50vh] overflow-hidden rounded-[var(--radius-md)] border border-[var(--border-color)] bg-[var(--bg-secondary)] shadow-[var(--shadow)] animate-in fade-in-0 zoom-in-95"
				sideOffset={4}
			>
				<Select.Viewport class="p-1 overflow-y-auto">
					{#each options as option (option.value)}
						<Select.Item
							value={option.value}
							label={option.label}
							class="relative flex cursor-pointer select-none items-center rounded-[var(--radius-sm)] py-3 sm:py-2 pl-8 pr-3 text-base sm:text-sm text-[var(--text-primary)] outline-none transition-colors hover:bg-[var(--bg-tertiary)] focus:bg-[var(--bg-tertiary)] active:bg-[var(--bg-tertiary)] data-[highlighted]:bg-[var(--sky)] data-[highlighted]:text-white data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
						>
							{#snippet children({ selected })}
								{#if selected}
									<span class="absolute left-2 flex h-4 w-4 items-center justify-center">
										<svg
											class="h-4 w-4 text-[var(--sky)]"
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

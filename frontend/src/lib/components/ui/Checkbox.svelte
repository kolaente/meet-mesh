<script lang="ts">
	interface Props {
		checked?: boolean;
		label?: string;
		description?: string;
		disabled?: boolean;
		onchange?: (checked: boolean) => void;
	}

	let {
		checked = $bindable(false),
		label,
		description,
		disabled = false,
		onchange
	}: Props = $props();

	function handleChange() {
		onchange?.(checked);
	}
</script>

<label class="checkbox-wrapper" class:disabled>
	<input
		type="checkbox"
		bind:checked
		{disabled}
		onchange={handleChange}
		class="checkbox"
	/>
	{#if label || description}
		<div class="label-container">
			{#if label}
				<span class="label">{label}</span>
			{/if}
			{#if description}
				<span class="description">{description}</span>
			{/if}
		</div>
	{/if}
</label>

<style>
	.checkbox-wrapper {
		display: flex;
		align-items: flex-start;
		gap: 0.75rem;
		cursor: pointer;
	}

	.checkbox-wrapper.disabled {
		cursor: not-allowed;
		opacity: 0.5;
	}

	.checkbox {
		width: 20px;
		height: 20px;
		min-width: 20px;
		border: 2px solid var(--border-color);
		border-radius: var(--radius-sm);
		background: var(--bg-secondary);
		cursor: pointer;
		appearance: none;
		display: grid;
		place-content: center;
		margin-top: 2px;
		transition: background var(--transition), border-color var(--transition), box-shadow var(--transition);
	}

	.checkbox:disabled {
		cursor: not-allowed;
	}

	.checkbox:not(:disabled):hover {
		box-shadow: var(--shadow-sm);
	}

	.checkbox:checked {
		background: var(--sky);
		border-color: var(--sky);
	}

	.checkbox:checked::before {
		content: '';
		width: 10px;
		height: 10px;
		background: white;
		clip-path: polygon(14% 44%, 0 65%, 50% 100%, 100% 16%, 80% 0%, 43% 62%);
	}

	.checkbox:focus-visible {
		outline: 2px solid var(--sky);
		outline-offset: 2px;
	}

	.label-container {
		display: flex;
		flex-direction: column;
		gap: 0.125rem;
	}

	.label {
		color: var(--text-primary);
		font-size: 0.9rem;
		font-weight: 500;
		line-height: 1.5;
	}

	.description {
		color: var(--text-secondary);
		font-size: 0.8rem;
		line-height: 1.4;
	}
</style>

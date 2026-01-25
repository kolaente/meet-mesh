<script lang="ts">
	import type { Snippet } from 'svelte';

	interface QuickAction {
		href: string;
		icon: Snippet;
		title: string;
		description: string;
	}

	interface Props {
		actions?: QuickAction[];
	}

	let { actions = [] }: Props = $props();

	// Color styles for each action position
	const colorStyles = [
		{
			// First: cyan/sky
			bg: 'var(--sky)',
			borderGlow: 'rgba(14, 165, 233, 0.3)',
			shadowGlow: 'rgba(14, 165, 233, 0.25)'
		},
		{
			// Second: purple/violet
			bg: 'var(--violet)',
			borderGlow: 'rgba(139, 92, 246, 0.3)',
			shadowGlow: 'rgba(139, 92, 246, 0.25)'
		},
		{
			// Third: blue/sky
			bg: 'var(--sky)',
			borderGlow: 'rgba(14, 165, 233, 0.3)',
			shadowGlow: 'rgba(14, 165, 233, 0.25)'
		}
	];

	function getColorStyle(index: number) {
		return colorStyles[index % colorStyles.length];
	}
</script>

<div class="quick-actions">
	{#each actions as action, index}
		{@const style = getColorStyle(index)}
		<a
			href={action.href}
			class="quick-action"
			style="--icon-bg: {style.bg}; --icon-border-glow: {style.borderGlow}; --icon-shadow-glow: {style.shadowGlow};"
		>
			<div class="quick-action-icon">
				{@render action.icon()}
			</div>
			<div class="quick-action-text">
				<div class="quick-action-title">{action.title}</div>
				<div class="quick-action-desc">{action.description}</div>
			</div>
			<svg
				class="quick-action-arrow"
				width="16"
				height="16"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
				aria-hidden="true"
			>
				<path stroke-linecap="round" stroke-width="2.5" d="M9 5l7 7-7 7" />
			</svg>
		</a>
	{/each}
</div>

<style>
	.quick-actions {
		padding: 0.75rem;
	}

	.quick-action {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem;
		border-radius: var(--radius);
		text-decoration: none;
		color: var(--text-primary);
		transition: all var(--transition);
		border: 2px solid transparent;
	}

	.quick-action:hover {
		background: var(--bg-tertiary);
		border-color: var(--border-color);
		transform: translateX(4px);
	}

	.quick-action-icon {
		width: 36px;
		height: 36px;
		border-radius: var(--radius);
		border: var(--border);
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--icon-bg);
		color: white;
		flex-shrink: 0;
		transition:
			box-shadow var(--transition),
			border-color var(--transition);
	}

	.quick-action-icon :global(svg) {
		width: 16px;
		height: 16px;
	}

	/* Dark mode glow effect */
	:global([data-theme='dark']) .quick-action-icon {
		border-color: var(--icon-border-glow);
		box-shadow: 0 0 16px var(--icon-shadow-glow);
	}

	.quick-action-text {
		flex: 1;
		min-width: 0;
	}

	.quick-action-title {
		font-weight: 700;
		font-size: 0.85rem;
	}

	.quick-action-desc {
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.quick-action-arrow {
		color: var(--text-muted);
		flex-shrink: 0;
		transition:
			transform var(--transition),
			color var(--transition);
	}

	.quick-action:hover .quick-action-arrow {
		transform: translateX(4px);
		color: var(--sky);
	}
</style>

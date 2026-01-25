<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		label: string;
		value: number | string;
		icon?: Snippet;
		trend?: { value: string; type: 'up' | 'alert' };
		color: 'sky' | 'violet' | 'amber' | 'emerald';
	}

	let { label, value, icon, trend, color }: Props = $props();

	const colorStyles = {
		sky: {
			bg: 'var(--sky)',
			borderGlow: 'rgba(14, 165, 233, 0.3)',
			shadowGlow: 'rgba(14, 165, 233, 0.25)'
		},
		violet: {
			bg: 'var(--violet)',
			borderGlow: 'rgba(139, 92, 246, 0.3)',
			shadowGlow: 'rgba(139, 92, 246, 0.25)'
		},
		amber: {
			bg: 'var(--amber)',
			borderGlow: 'rgba(245, 158, 11, 0.3)',
			shadowGlow: 'rgba(245, 158, 11, 0.25)'
		},
		emerald: {
			bg: 'var(--emerald)',
			borderGlow: 'rgba(16, 185, 129, 0.3)',
			shadowGlow: 'rgba(16, 185, 129, 0.25)'
		}
	};

	const currentColor = $derived(colorStyles[color]);
</script>

<div
	class="stat-card"
	style="--icon-bg: {currentColor.bg}; --icon-border-glow: {currentColor.borderGlow}; --icon-shadow-glow: {currentColor.shadowGlow};"
>
	<div class="stat-header">
		{#if icon}
			<div class="stat-icon">
				{@render icon()}
			</div>
		{/if}
		{#if trend}
			<span class="stat-trend" class:up={trend.type === 'up'} class:alert={trend.type === 'alert'}>
				{trend.value}
			</span>
		{/if}
	</div>
	<div class="stat-value">{value}</div>
	<div class="stat-label">{label}</div>
</div>

<style>
	.stat-card {
		background: var(--bg-secondary);
		border: var(--border);
		border-radius: var(--radius-lg);
		padding: 1.25rem;
		box-shadow: var(--shadow);
		transition: all var(--transition);
		position: relative;
		overflow: hidden;
		animation: slideUp 0.3s ease backwards;
	}

	.stat-card:hover {
		transform: translate(-2px, -2px);
		box-shadow: var(--shadow-hover);
	}

	.stat-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.75rem;
	}

	.stat-icon {
		width: 42px;
		height: 42px;
		border-radius: 10px;
		border: var(--border);
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--icon-bg);
		color: white;
		transition: box-shadow var(--transition), border-color var(--transition);
	}

	.stat-icon :global(svg) {
		width: 20px;
		height: 20px;
	}

	/* Dark mode glow effect */
	:global([data-theme='dark']) .stat-icon {
		border-color: var(--icon-border-glow);
		box-shadow: 0 0 20px var(--icon-shadow-glow);
	}

	.stat-trend {
		font-size: 0.7rem;
		font-weight: 700;
		padding: 0.2rem 0.5rem;
		border-radius: var(--radius);
		border: 1px solid var(--border-color);
	}

	.stat-trend.up {
		background: var(--emerald);
		color: white;
	}

	.stat-trend.alert {
		background: var(--amber);
		color: white;
	}

	.stat-value {
		font-size: 2rem;
		font-weight: 800;
		letter-spacing: -0.03em;
		line-height: 1;
		margin-bottom: 0.25rem;
		color: var(--text-primary);
	}

	.stat-label {
		font-size: 0.8rem;
		font-weight: 500;
		color: var(--text-muted);
	}

	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(15px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>

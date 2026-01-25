<script lang="ts">
	type Variant = 'pending' | 'confirmed' | 'active' | 'cancelled' | 'draft';
	type Size = 'sm' | 'md';

	interface Props {
		variant?: Variant;
		size?: Size;
	}

	let { variant = 'pending', size = 'md' }: Props = $props();

	const labels: Record<Variant, string> = {
		pending: 'Pending',
		confirmed: 'Confirmed',
		active: 'Active',
		cancelled: 'Cancelled',
		draft: 'Draft'
	};
</script>

<span
	class="badge badge-{variant} badge-{size}"
>
	{labels[variant]}
</span>

<style>
	.badge {
		display: inline-flex;
		align-items: center;
		font-weight: 700;
		border-radius: var(--radius, 8px);
		border: 1px solid var(--border-color, #1a1a1a);
		color: white;
	}

	/* Sizes */
	.badge-sm {
		font-size: 0.7rem;
		padding: 0.2rem 0.5rem;
	}

	.badge-md {
		font-size: 0.75rem;
		padding: 0.25rem 0.6rem;
	}

	/* Variants - Neo-brutalist with solid accent backgrounds */
	.badge-pending {
		background: var(--amber, #f59e0b);
	}

	.badge-confirmed {
		background: var(--emerald, #10b981);
	}

	.badge-active {
		background: var(--sky, #0ea5e9);
	}

	.badge-cancelled {
		background: var(--text-muted, #888888);
	}

	.badge-draft {
		background: var(--text-secondary, #555555);
	}
</style>

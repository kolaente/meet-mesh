<script lang="ts">
	type Variant = 'pending' | 'confirmed' | 'active' | 'cancelled' | 'draft';
	type Size = 'sm' | 'md';

	interface Props {
		variant?: Variant;
		size?: Size;
	}

	let { variant = 'pending', size = 'md' }: Props = $props();

	const baseClasses = 'inline-flex items-center font-medium rounded-full';

	const variantClasses: Record<Variant, string> = {
		pending: 'bg-amber-50 text-amber-700 border border-amber-200',
		confirmed: 'bg-green-50 text-green-700 border border-green-200',
		active: 'bg-indigo-50 text-indigo-700 border border-indigo-200',
		cancelled: 'bg-gray-100 text-gray-600 border border-gray-200',
		draft: 'bg-slate-100 text-slate-600 border border-slate-200'
	};

	const sizeClasses: Record<Size, string> = {
		sm: 'px-2 py-0.5 text-xs',
		md: 'px-2.5 py-1 text-sm'
	};

	const labels: Record<Variant, string> = {
		pending: 'Pending',
		confirmed: 'Confirmed',
		active: 'Active',
		cancelled: 'Cancelled',
		draft: 'Draft'
	};
</script>

<span class="{baseClasses} {variantClasses[variant]} {sizeClasses[size]}">
	{labels[variant]}
</span>

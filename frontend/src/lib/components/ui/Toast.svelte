<script lang="ts">
	import { fly } from 'svelte/transition';

	type Variant = 'success' | 'error' | 'info';

	interface Props {
		variant?: Variant;
		message: string;
		onClose: () => void;
	}

	let { variant = 'info', message, onClose }: Props = $props();

	const variantClasses: Record<Variant, string> = {
		success: 'bg-green-50 text-green-800 border-green-200',
		error: 'bg-red-50 text-red-800 border-red-200',
		info: 'bg-blue-50 text-blue-800 border-blue-200'
	};

	const iconPaths: Record<Variant, string> = {
		success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
		error: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
		info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
	};
</script>

<div
	role="alert"
	class="flex items-center gap-3 px-4 py-3 rounded-[var(--radius-md)] border shadow-[var(--shadow-md)] {variantClasses[variant]}"
	transition:fly={{ x: 300, duration: 300 }}
>
	<svg
		class="h-5 w-5 flex-shrink-0"
		xmlns="http://www.w3.org/2000/svg"
		fill="none"
		viewBox="0 0 24 24"
		stroke="currentColor"
		stroke-width="2"
		aria-hidden="true"
	>
		<path stroke-linecap="round" stroke-linejoin="round" d={iconPaths[variant]} />
	</svg>

	<p class="flex-1 text-sm font-medium">{message}</p>

	<button
		type="button"
		onclick={onClose}
		class="flex-shrink-0 p-1 rounded hover:bg-black/5 transition-colors"
		aria-label="Dismiss notification"
	>
		<svg
			class="h-4 w-4"
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke="currentColor"
			stroke-width="2"
			aria-hidden="true"
		>
			<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
		</svg>
	</button>
</div>

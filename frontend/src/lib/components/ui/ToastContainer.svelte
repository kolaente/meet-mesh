<script lang="ts">
	import { onMount } from 'svelte';
	import { getToasts } from '$lib/stores/toast.svelte';
	import Toast from './Toast.svelte';

	const toasts = getToasts();

	const AUTO_DISMISS_MS = 4000;

	onMount(() => {
		// Set up auto-dismiss interval
		const interval = setInterval(() => {
			const allToasts = toasts.all;
			if (allToasts.length > 0) {
				toasts.remove(allToasts[0].id);
			}
		}, AUTO_DISMISS_MS);

		return () => clearInterval(interval);
	});
</script>

<div
	class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none"
	aria-live="polite"
	aria-label="Notifications"
>
	{#each toasts.all as toast (toast.id)}
		<div class="pointer-events-auto">
			<Toast variant={toast.variant} message={toast.message} onClose={() => toasts.remove(toast.id)} />
		</div>
	{/each}
</div>

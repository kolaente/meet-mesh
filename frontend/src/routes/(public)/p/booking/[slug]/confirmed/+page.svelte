<script lang="ts">
	import { page } from '$app/state';
	import { Card } from '$lib/components/ui';
</script>

<Card class="text-center py-8">
	<div class="text-green-500 mb-4">
		<svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
		</svg>
	</div>
	<h1 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Booking Confirmed!</h1>
	<p class="text-gray-600 dark:text-gray-400">You will receive a confirmation email shortly.</p>
</Card>

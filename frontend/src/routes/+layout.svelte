<script lang="ts">
	import '@fontsource-variable/outfit';
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ToastContainer } from '$lib/components/ui';
	import { onMount } from 'svelte';
	import { getDateFormat } from '$lib/stores/dateFormat.svelte';

	let { children } = $props();

	const dateFormat = getDateFormat();

	onMount(() => {
		dateFormat.init();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
<ToastContainer />

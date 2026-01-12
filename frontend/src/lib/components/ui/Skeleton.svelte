<script lang="ts">
	interface Props {
		class?: string;
		width?: string;
		height?: string;
	}

	let { class: className = '', width, height }: Props = $props();

	const style = $derived(
		[
			width ? `width: ${width}` : '',
			height ? `height: ${height}` : ''
		].filter(Boolean).join('; ')
	);
</script>

<div
	class="animate-pulse bg-slate-200 rounded-[var(--radius-sm)] {className}"
	style={style || undefined}
	aria-hidden="true"
></div>

<style>
	@keyframes shimmer {
		0% {
			background-position: -200% 0;
		}
		100% {
			background-position: 200% 0;
		}
	}

	div {
		background: linear-gradient(
			90deg,
			rgb(226 232 240) 0%,
			rgb(241 245 249) 50%,
			rgb(226 232 240) 100%
		);
		background-size: 200% 100%;
		animation: shimmer 1.5s ease-in-out infinite;
	}
</style>

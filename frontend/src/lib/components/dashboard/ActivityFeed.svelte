<script lang="ts">
	interface Activity {
		type: 'booking' | 'poll' | 'vote';
		text: string;
		time: string;
		boldParts?: string[];
	}

	interface Props {
		activities: Activity[];
	}

	let { activities }: Props = $props();

	function formatText(text: string, boldParts?: string[]): string {
		if (!boldParts || boldParts.length === 0) {
			return text;
		}

		let result = text;
		for (const part of boldParts) {
			result = result.replace(part, `<strong>${part}</strong>`);
		}
		return result;
	}
</script>

<div class="activity-list">
	{#each activities as activity}
		<div class="activity-item">
			<div class="activity-dot {activity.type}"></div>
			<div>
				<div class="activity-text">
					{@html formatText(activity.text, activity.boldParts)}
				</div>
				<div class="activity-time">{activity.time}</div>
			</div>
		</div>
	{/each}
</div>

<style>
	.activity-list {
		padding: 0.5rem 0;
	}

	.activity-item {
		display: flex;
		gap: 0.75rem;
		padding: 0.65rem 1.25rem;
		transition: background var(--transition);
	}

	.activity-item:hover {
		background: var(--bg-tertiary);
	}

	.activity-dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		margin-top: 6px;
		flex-shrink: 0;
		border: 1px solid var(--border-color);
	}

	.activity-dot.booking {
		background: var(--rose);
	}

	.activity-dot.poll {
		background: var(--violet);
	}

	.activity-dot.vote {
		background: var(--emerald);
	}

	.activity-text {
		font-size: 0.8rem;
		line-height: 1.4;
		color: var(--text-primary);
	}

	.activity-text :global(strong) {
		font-weight: 700;
	}

	.activity-time {
		font-size: 0.7rem;
		color: var(--text-muted);
		margin-top: 0.15rem;
	}
</style>

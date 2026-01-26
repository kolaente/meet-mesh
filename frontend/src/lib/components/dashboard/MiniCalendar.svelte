<script lang="ts">
	import { getDateFormat } from '$lib/stores/dateFormat.svelte';

	interface Props {
		events?: Date[];
		onDateClick?: (date: Date) => void;
	}

	let { events = [], onDateClick }: Props = $props();

	// Current displayed month (starts with current month)
	let displayedDate = $state(new Date());

	const dateFormat = getDateFormat();
	const weekDays = $derived(dateFormat.getWeekDays('narrow'));
	const MONTHS = [
		'January',
		'February',
		'March',
		'April',
		'May',
		'June',
		'July',
		'August',
		'September',
		'October',
		'November',
		'December'
	];

	function isSameDay(date1: Date, date2: Date): boolean {
		return (
			date1.getFullYear() === date2.getFullYear() &&
			date1.getMonth() === date2.getMonth() &&
			date1.getDate() === date2.getDate()
		);
	}

	function hasEvent(date: Date): boolean {
		return events.some((eventDate) => isSameDay(eventDate, date));
	}

	function isToday(date: Date): boolean {
		const today = new Date();
		return isSameDay(date, today);
	}

	function prevMonth() {
		displayedDate = new Date(displayedDate.getFullYear(), displayedDate.getMonth() - 1, 1);
	}

	function nextMonth() {
		displayedDate = new Date(displayedDate.getFullYear(), displayedDate.getMonth() + 1, 1);
	}

	function handleDateClick(date: Date) {
		onDateClick?.(date);
	}

	// Calculate calendar days for the displayed month
	const calendarDays = $derived.by(() => {
		const year = displayedDate.getFullYear();
		const month = displayedDate.getMonth();

		// First day of the month
		const firstDayOfMonth = new Date(year, month, 1);
		const startDayOfWeek = dateFormat.getDayIndex(firstDayOfMonth);

		// Last day of the month
		const lastDayOfMonth = new Date(year, month + 1, 0);
		const daysInMonth = lastDayOfMonth.getDate();

		// Previous month days to fill
		const prevMonthLastDay = new Date(year, month, 0).getDate();

		const days: { date: Date; isOtherMonth: boolean }[] = [];

		// Add days from previous month
		for (let i = startDayOfWeek - 1; i >= 0; i--) {
			days.push({
				date: new Date(year, month - 1, prevMonthLastDay - i),
				isOtherMonth: true
			});
		}

		// Add days of current month
		for (let day = 1; day <= daysInMonth; day++) {
			days.push({
				date: new Date(year, month, day),
				isOtherMonth: false
			});
		}

		// Add days from next month to complete the grid (6 rows * 7 days = 42)
		const remainingDays = 42 - days.length;
		for (let day = 1; day <= remainingDays; day++) {
			days.push({
				date: new Date(year, month + 1, day),
				isOtherMonth: true
			});
		}

		return days;
	});

	const monthYearLabel = $derived(
		`${MONTHS[displayedDate.getMonth()]} ${displayedDate.getFullYear()}`
	);
</script>

<div class="calendar-body">
	<div class="calendar-header">
		<span class="calendar-month">{monthYearLabel}</span>
		<div class="calendar-nav">
			<button type="button" onclick={prevMonth} aria-label="Previous month">
				<svg
					width="12"
					height="12"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					aria-hidden="true"
				>
					<path stroke-linecap="round" stroke-width="3" d="M15 19l-7-7 7-7" />
				</svg>
			</button>
			<button type="button" onclick={nextMonth} aria-label="Next month">
				<svg
					width="12"
					height="12"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
					aria-hidden="true"
				>
					<path stroke-linecap="round" stroke-width="3" d="M9 5l7 7-7 7" />
				</svg>
			</button>
		</div>
	</div>
	<div class="calendar-grid">
		{#each weekDays as day}
			<div class="calendar-day-label">{day}</div>
		{/each}
		{#each calendarDays as { date, isOtherMonth }}
			<button
				type="button"
				class="calendar-day"
				class:other-month={isOtherMonth}
				class:today={isToday(date)}
				class:has-event={hasEvent(date)}
				onclick={() => handleDateClick(date)}
			>
				{date.getDate()}
			</button>
		{/each}
	</div>
</div>

<style>
	.calendar-body {
		padding: 1rem;
	}

	.calendar-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 0.75rem;
	}

	.calendar-month {
		font-weight: 700;
		font-size: 0.85rem;
		color: var(--text-primary);
	}

	.calendar-nav {
		display: flex;
		gap: 0.25rem;
	}

	.calendar-nav button {
		width: 26px;
		height: 26px;
		border-radius: var(--radius);
		border: var(--border);
		background: var(--bg-secondary);
		color: var(--text-primary);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all var(--transition);
	}

	.calendar-nav button:hover {
		background: var(--cyan);
		color: white;
	}

	.calendar-grid {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: 2px;
	}

	.calendar-day-label {
		text-align: center;
		font-size: 0.65rem;
		font-weight: 700;
		color: var(--text-muted);
		padding: 0.4rem;
		text-transform: uppercase;
	}

	.calendar-day {
		text-align: center;
		padding: 0.4rem;
		font-size: 0.75rem;
		font-weight: 600;
		border-radius: var(--radius);
		cursor: pointer;
		transition: all var(--transition);
		position: relative;
		border: 1px solid transparent;
		background: transparent;
		color: var(--text-primary);
	}

	.calendar-day:hover {
		background: var(--bg-tertiary);
		border-color: var(--border-color);
	}

	.calendar-day.today {
		background: var(--sky);
		color: white;
		border-color: var(--border-color);
	}

	.calendar-day.has-event::after {
		content: '';
		position: absolute;
		bottom: 2px;
		left: 50%;
		transform: translateX(-50%);
		width: 4px;
		height: 4px;
		border-radius: 50%;
		background: var(--rose);
	}

	.calendar-day.today.has-event::after {
		background: white;
	}

	.calendar-day.other-month {
		color: var(--text-muted);
		opacity: 0.4;
	}
</style>

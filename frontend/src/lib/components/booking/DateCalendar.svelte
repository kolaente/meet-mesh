<script lang="ts">
	interface Props {
		availableDates: string[];
		selectedDate?: string;
		class?: string;
	}

	let {
		availableDates,
		selectedDate = $bindable(),
		class: className = ''
	}: Props = $props();

	// Current view state
	let viewDate = $state(new Date());

	// Compute available dates as a Set for fast lookup
	let availableDateSet = $derived(new Set(availableDates.map((d) => d.split('T')[0])));

	// Calendar computation
	let year = $derived(viewDate.getFullYear());
	let month = $derived(viewDate.getMonth());
	let monthName = $derived(viewDate.toLocaleString('default', { month: 'long', year: 'numeric' }));

	let firstDayOfMonth = $derived(new Date(year, month, 1).getDay());
	let daysInMonth = $derived(new Date(year, month + 1, 0).getDate());

	// Generate calendar grid
	let calendarDays = $derived.by(() => {
		const days: (number | null)[] = [];

		// Add empty cells for days before the first of the month
		for (let i = 0; i < firstDayOfMonth; i++) {
			days.push(null);
		}

		// Add days of the month
		for (let day = 1; day <= daysInMonth; day++) {
			days.push(day);
		}

		return days;
	});

	const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

	function prevMonth() {
		viewDate = new Date(year, month - 1, 1);
	}

	function nextMonth() {
		viewDate = new Date(year, month + 1, 1);
	}

	function selectDate(day: number) {
		const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
		selectedDate = dateStr;
	}

	function isAvailable(day: number): boolean {
		const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
		return availableDateSet.has(dateStr);
	}

	function isSelected(day: number): boolean {
		if (!selectedDate) return false;
		const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
		return selectedDate.startsWith(dateStr);
	}

	function isToday(day: number): boolean {
		const today = new Date();
		return (
			today.getFullYear() === year &&
			today.getMonth() === month &&
			today.getDate() === day
		);
	}
</script>

<div class="select-none {className}">
	<!-- Header with month navigation -->
	<div class="flex items-center justify-between mb-4">
		<button
			type="button"
			onclick={prevMonth}
			class="p-2 min-w-[44px] min-h-[44px] flex items-center justify-center text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-[var(--radius-md)] transition-colors"
			aria-label="Previous month"
		>
			<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
		</button>

		<h2 class="text-base sm:text-lg font-semibold text-gray-900">{monthName}</h2>

		<button
			type="button"
			onclick={nextMonth}
			class="p-2 min-w-[44px] min-h-[44px] flex items-center justify-center text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-[var(--radius-md)] transition-colors"
			aria-label="Next month"
		>
			<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
			</svg>
		</button>
	</div>

	<!-- Weekday headers - abbreviated on mobile -->
	<div class="grid grid-cols-7 mb-2">
		{#each weekDays as day}
			<div class="text-center text-xs sm:text-sm font-medium text-gray-500 py-2">
				<span class="sm:hidden">{day.charAt(0)}</span>
				<span class="hidden sm:inline">{day}</span>
			</div>
		{/each}
	</div>

	<!-- Calendar grid with touch-friendly sizing -->
	<div class="grid grid-cols-7 gap-0.5 sm:gap-1">
		{#each calendarDays as day}
			{#if day === null}
				<div class="aspect-square min-h-[40px] sm:min-h-0"></div>
			{:else}
				{@const available = isAvailable(day)}
				{@const selected = isSelected(day)}
				{@const today = isToday(day)}
				<button
					type="button"
					onclick={() => available && selectDate(day)}
					disabled={!available}
					class="aspect-square min-h-[40px] sm:min-h-0 flex flex-col items-center justify-center relative rounded-[var(--radius-md)] text-xs sm:text-sm transition-colors
						{selected
							? 'bg-indigo-600 text-white'
							: available
								? 'hover:bg-gray-100 text-gray-900 cursor-pointer active:bg-gray-200'
								: 'text-gray-300 cursor-not-allowed'}
						{today && !selected ? 'ring-2 ring-indigo-600 ring-inset' : ''}"
				>
					<span>{day}</span>
					{#if available && !selected}
						<span class="absolute bottom-1 sm:bottom-1.5 w-1 h-1 rounded-full bg-indigo-500"></span>
					{/if}
				</button>
			{/if}
		{/each}
	</div>
</div>

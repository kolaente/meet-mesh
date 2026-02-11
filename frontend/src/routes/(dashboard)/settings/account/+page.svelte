<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { getDateFormat } from '$lib/stores/dateFormat.svelte';
	import { getToasts } from '$lib/stores/toast.svelte';
	import { DashboardHeader } from '$lib/components/dashboard';
	import { Card, Button, Input, Select, Spinner } from '$lib/components/ui';

	const auth = getAuth();
	const dateFormat = getDateFormat();
	const toast = getToasts();

	let displayName = $state('');
	let saving = $state(false);
	let loading = $state(true);

	// Reactive bindings for date format selects
	let timeFormatValue = $derived(dateFormat.timeFormat);
	let weekStartDayValue = $derived(dateFormat.weekStartDay);

	const timeFormatOptions = [
		{ value: '12h', label: '12-hour (2:00 PM)' },
		{ value: '24h', label: '24-hour (14:00)' }
	];

	const weekStartOptions = [
		{ value: 'sunday', label: 'Sunday' },
		{ value: 'monday', label: 'Monday' }
	];

	onMount(async () => {
		// Load current user data
		const { data } = await api.GET('/auth/me');
		if (data) {
			displayName = data.name ?? '';
		}
		loading = false;
	});

	async function saveProfile() {
		saving = true;
		try {
			const { data, error } = await api.PUT('/auth/me', {
				body: { name: displayName }
			});

			if (error) {
				toast.error('Failed to save profile');
				return;
			}

			if (data) {
				toast.success('Profile updated');
			}
		} finally {
			saving = false;
		}
	}
</script>

<DashboardHeader title="Account Settings" />

{#if loading}
	<div class="flex items-center justify-center py-12">
		<Spinner size="lg" />
	</div>
{:else}
	<div class="space-y-6">
		<!-- Profile Section -->
		<section>
			<div class="mb-4">
				<h2 class="text-lg font-medium text-text-primary">Profile</h2>
				<p class="text-sm text-text-secondary">Manage your display name and profile information</p>
			</div>

			<Card>
				{#snippet children()}
					<div class="space-y-6">
						<div class="flex flex-col sm:flex-row sm:items-start gap-4">
							<div class="flex-1">
								<Input
									name="displayName"
									label="Display Name"
									description="This name is shown to guests on your booking pages and polls."
									bind:value={displayName}
									placeholder="Enter your name"
								/>
							</div>
						</div>

						<div>
							<Input
								name="email"
								label="Email"
								description="Your email address is managed by your identity provider and cannot be changed here."
								value={auth.user?.email ?? ''}
								disabled={true}
							/>
						</div>

						<div class="flex justify-end pt-2 border-t border-border">
							<Button variant="primary" onclick={saveProfile} disabled={saving}>
								{#snippet children()}
									{#if saving}
										<Spinner size="sm" />
										Saving...
									{:else}
										Save Profile
									{/if}
								{/snippet}
							</Button>
						</div>
					</div>
				{/snippet}
			</Card>
		</section>

		<!-- Date & Time Format Section -->
		<section>
			<div class="mb-4">
				<h2 class="text-lg font-medium text-text-primary">Date & Time Format</h2>
				<p class="text-sm text-text-secondary">Customize how dates and times are displayed</p>
			</div>

			<Card>
				{#snippet children()}
					<div class="space-y-6">
						<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
							<div>
								<p class="font-medium text-text-primary">Time Format</p>
								<p class="text-sm text-text-secondary">Choose 12-hour or 24-hour time display</p>
							</div>
							<div class="w-full sm:w-48">
								<Select
									name="timeFormat"
									options={timeFormatOptions}
									value={timeFormatValue}
									onchange={(value) => dateFormat.setTimeFormat(value as '12h' | '24h')}
								/>
							</div>
						</div>

						<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
							<div>
								<p class="font-medium text-text-primary">Week Starts On</p>
								<p class="text-sm text-text-secondary">Choose which day starts your week</p>
							</div>
							<div class="w-full sm:w-48">
								<Select
									name="weekStartDay"
									options={weekStartOptions}
									value={weekStartDayValue}
									onchange={(value) => dateFormat.setWeekStartDay(value as 'sunday' | 'monday')}
								/>
							</div>
						</div>

						<div class="pt-2 border-t border-border">
							<Button variant="secondary" onclick={() => dateFormat.reset()}>
								{#snippet children()}Reset to Browser Defaults{/snippet}
							</Button>
						</div>
					</div>
				{/snippet}
			</Card>
		</section>
	</div>
{/if}


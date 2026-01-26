<script lang="ts">
    import { page } from '$app/state';
    import { Card } from '$lib/components/ui';

    // BookingStatus: 1=pending, 2=confirmed, 3=declined
    const status = $derived(Number(page.url.searchParams.get('status')) || 2);
    const isPending = $derived(status === 1);
</script>

<svelte:head>
    <title>{isPending ? 'Booking Pending' : 'Booking Confirmed'} | Meet Mesh</title>
</svelte:head>

<Card class="text-center py-8">
    {#if isPending}
        <div class="text-amber-500 mb-4">
            <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        </div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Booking Pending Approval</h1>
        <p class="text-gray-600 dark:text-gray-400">Your booking request has been submitted.</p>
        <p class="text-gray-600 dark:text-gray-400 mt-2">You will receive a confirmation email once the organizer approves your request.</p>
    {:else}
        <div class="text-green-500 mb-4">
            <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
        </div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-2">Booking Confirmed!</h1>
        <p class="text-gray-600 dark:text-gray-400">You will receive a confirmation email shortly.</p>
    {/if}
</Card>

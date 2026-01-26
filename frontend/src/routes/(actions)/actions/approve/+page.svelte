<script lang="ts">
  import { page } from '$app/stores'
  import { onMount } from 'svelte'
  import { Card, Spinner } from '$lib/components/ui'

  type PageState = 'loading' | 'success' | 'error'

  let pageState: PageState = $state('loading')
  let message = $state('')

  onMount(async () => {
    const token = $page.url.searchParams.get('token')

    if (!token) {
      pageState = 'error'
      message = 'Missing token'
      return
    }

    try {
      const response = await fetch(`/api/actions/approve?token=${encodeURIComponent(token)}`)
      const data = await response.json()

      if (!response.ok) {
        pageState = 'error'
        message = data.message || 'Invalid or expired link'
      } else {
        pageState = 'success'
        message = data.message || 'The booking has been confirmed.'
      }
    } catch {
      pageState = 'error'
      message = 'An unexpected error occurred'
    }
  })
</script>

<svelte:head>
  <title>Approve Booking | Meet Mesh</title>
</svelte:head>

<Card>
  <div class="text-center py-6">
    {#if pageState === 'loading'}
      <div class="flex flex-col items-center gap-4">
        <div class="text-indigo-600">
          <Spinner size="lg" />
        </div>
        <p class="text-slate-600">Approving booking...</p>
      </div>
    {:else if pageState === 'success'}
      <div class="flex flex-col items-center gap-4">
        <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center">
          <svg
            class="w-8 h-8 text-green-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-semibold text-slate-900">Booking Approved!</h1>
          <p class="text-slate-600 mt-1">{message}</p>
        </div>
      </div>
    {:else}
      <div class="flex flex-col items-center gap-4">
        <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
          <svg
            class="w-8 h-8 text-red-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </div>
        <div>
          <h1 class="text-xl font-semibold text-slate-900">Unable to Approve</h1>
          <p class="text-slate-600 mt-1">{message}</p>
        </div>
      </div>
    {/if}
  </div>
</Card>

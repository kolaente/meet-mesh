<script lang="ts">
  import { page } from '$app/state'
  import { api } from '$lib/api/client'
  import type { components } from '$lib/api/types'
  import { Spinner, Card } from '$lib/components/ui'
  import { BookingPage } from '$lib/components/booking'
  import { PollPage } from '$lib/components/poll'

  type PublicLink = components['schemas']['Slot'] extends infer _ ? {
    type: components['schemas']['LinkType'];
    name: string;
    description?: string;
    custom_fields?: components['schemas']['CustomField'][];
    slots: components['schemas']['Slot'][];
    show_results?: boolean;
    require_email?: boolean;
  } : never;

  let link = $state<PublicLink | null>(null)
  let loading = $state(true)
  let error = $state<string | null>(null)

  const slug = $derived(page.params.slug ?? '')

  $effect(() => {
    if (slug) {
      loadLink()
    }
  })

  async function loadLink() {
    if (!slug) return

    loading = true
    error = null

    const { data, error: apiError } = await api.GET('/p/{slug}', {
      params: { path: { slug } }
    })

    if (apiError) {
      error = 'This scheduling page is not available.'
      loading = false
      return
    }

    link = data ?? null
    loading = false
  }
</script>

{#if loading}
  <div class="flex items-center justify-center py-12">
    <Spinner size="lg" />
  </div>
{:else if error}
  <Card class="text-center py-8">
    <div class="text-red-500 mb-4">
      <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
    </div>
    <h1 class="text-xl font-semibold text-gray-900 mb-2">Page Not Found</h1>
    <p class="text-gray-600">{error}</p>
  </Card>
{:else if link}
  {#if link.type === 1}
    <!-- Booking type -->
    <BookingPage {link} {slug} />
  {:else if link.type === 2}
    <!-- Poll type -->
    <PollPage {link} {slug} />
  {:else}
    <Card class="text-center py-8">
      <p class="text-gray-600">Unknown link type</p>
    </Card>
  {/if}
{/if}

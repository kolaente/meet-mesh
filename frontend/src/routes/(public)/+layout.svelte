<script lang="ts">
  import { page } from '$app/state'
  import { fade } from 'svelte/transition'
  import { pageIn, pageOut } from '$lib/utils/transitions'
  import ThemeToggle from '$lib/components/ui/ThemeToggle.svelte'
  import { DateFormatToggle } from '$lib/components/ui'

  let { children } = $props()
</script>

<div class="min-h-screen bg-[var(--bg-primary)]">
  <!-- Theme toggle in top-right corner -->
  <div class="fixed top-4 right-4 z-50">
    <ThemeToggle />
  </div>

  <div class="max-w-2xl mx-auto px-4 py-8">
    {#key page.url.pathname}
      <div in:fade={pageIn} out:fade={pageOut}>
        {@render children()}
      </div>
    {/key}
  </div>

  <!-- Date Format Toggle -->
  <div class="fixed bottom-4 right-4 z-50">
    <DateFormatToggle />
  </div>
</div>

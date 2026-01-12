<script lang="ts">
  import { page } from '$app/state'
  import { getAuth } from '$lib/stores/auth.svelte'
  import { goto } from '$app/navigation'
  import { fade } from 'svelte/transition'
  import { Sidebar, MobileMenu } from '$lib/components/dashboard'
  import { LoadingScreen } from '$lib/components/shared'
  import { pageIn, pageOut } from '$lib/utils/transitions'

  let { children } = $props()
  const auth = getAuth()

  let mobileMenuOpen = $state(false)

  $effect(() => { auth.check() })
  $effect(() => {
    if (!auth.loading && !auth.isAuthenticated) {
      goto('/auth/login')
    }
  })
</script>

{#if auth.loading}
  <LoadingScreen />
{:else if auth.isAuthenticated}
  <div class="flex min-h-screen">
    <!-- Mobile menu (hamburger + drawer) -->
    <MobileMenu bind:open={mobileMenuOpen} />

    <!-- Sidebar (hidden on mobile, icons on tablet, full on desktop) -->
    <Sidebar />

    <!-- Main content with responsive padding -->
    <main class="flex-1 pt-16 px-4 pb-6 sm:pt-6 sm:px-6 bg-slate-50">
      {#key page.url.pathname}
        <div in:fade={pageIn} out:fade={pageOut}>
          {@render children()}
        </div>
      {/key}
    </main>
  </div>
{/if}

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

    <!-- Main content with margin for fixed sidebar -->
    <main class="main-content">
      {#key page.url.pathname}
        <div in:fade={pageIn} out:fade={pageOut}>
          {@render children()}
        </div>
      {/key}
    </main>
  </div>
{/if}

<style>
  .main-content {
    flex: 1;
    margin-left: var(--spacing-sidebar);
    padding: 1.5rem 2rem;
    background: var(--color-bg-primary);
    min-height: 100vh;
  }

  @media (max-width: 900px) {
    .main-content {
      margin-left: 0;
      padding-top: 4.5rem;
    }
  }
</style>

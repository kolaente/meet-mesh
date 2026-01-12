<script lang="ts">
  import { page } from '$app/state'
  import { fade } from 'svelte/transition'
  import { pageIn, pageOut } from '$lib/utils/transitions'

  let { children } = $props()
</script>

<div class="min-h-screen bg-slate-50">
  <div class="max-w-2xl mx-auto px-4 py-8">
    {#key page.url.pathname}
      <div in:fade={pageIn} out:fade={pageOut}>
        {@render children()}
      </div>
    {/key}
  </div>
</div>

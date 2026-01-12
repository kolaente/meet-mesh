<script lang="ts">
  import type { Snippet } from 'svelte'

  let { title, actions }: { title: string; actions?: Snippet } = $props()
</script>

<header class="flex items-center justify-between mb-6">
  <h1 class="text-2xl font-semibold text-gray-900">{title}</h1>
  {#if actions}
    <div class="flex items-center gap-3">
      {@render actions()}
    </div>
  {/if}
</header>

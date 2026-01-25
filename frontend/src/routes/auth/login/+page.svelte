<script lang="ts">
  import { onMount } from 'svelte'
  import { Spinner } from '$lib/components/ui'

  onMount(() => {
    window.location.href = '/api/auth/login'
  })
</script>

<div class="min-h-screen flex items-center justify-center">
  <div class="text-center">
    <Spinner size="lg" />
    <p class="mt-4 text-gray-600 dark:text-gray-400">Redirecting to login...</p>
  </div>
</div>

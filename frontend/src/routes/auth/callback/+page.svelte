<script lang="ts">
  import { goto } from '$app/navigation'
  import { Spinner } from '$lib/components/ui'
  import { getAuth } from '$lib/stores/auth.svelte'

  const auth = getAuth()

  $effect(() => {
    auth.check()
  })

  $effect(() => {
    if (!auth.loading && auth.isAuthenticated) {
      goto('/')
    }
  })
</script>

<div class="min-h-screen flex items-center justify-center">
  <div class="text-center">
    <Spinner size="lg" />
    <p class="mt-4 text-gray-600">Completing login...</p>
  </div>
</div>

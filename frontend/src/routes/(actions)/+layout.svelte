<script lang="ts">
  let { children } = $props()
</script>

<div class="min-h-screen bg-slate-50">
  <div class="max-w-lg mx-auto px-4 py-16">
    {@render children()}
  </div>
</div>

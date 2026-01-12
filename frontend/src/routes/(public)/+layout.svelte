<script lang="ts">
  let { children } = $props()
</script>

<div class="min-h-screen bg-slate-50">
  <div class="max-w-2xl mx-auto px-4 py-8">
    {@render children()}
  </div>
</div>

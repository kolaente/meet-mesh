<script lang="ts">
  import type { Snippet } from 'svelte'
  import ThemeToggle from '$lib/components/ui/ThemeToggle.svelte'

  let { title, subtitle, actions }: { title: string; subtitle?: string; actions?: Snippet } = $props()
</script>

<header class="header">
  <div class="header-content">
    <h1 class="header-title">{title}</h1>
    {#if subtitle}
      <p class="header-subtitle">{subtitle}</p>
    {/if}
  </div>
  <div class="header-actions">
    <ThemeToggle />
    {#if actions}
      {@render actions()}
    {/if}
  </div>
</header>

<style>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1.75rem;
    padding-bottom: 1.5rem;
    border-bottom: var(--border-light);
  }

  .header-content {
    display: flex;
    flex-direction: column;
  }

  .header-title {
    font-size: 1.75rem;
    font-weight: 800;
    letter-spacing: -0.03em;
    color: var(--text-primary);
  }

  .header-subtitle {
    color: var(--text-secondary);
    font-size: 0.9rem;
    margin-top: 0.25rem;
  }

  .header-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  /* Responsive */
  @media (max-width: 600px) {
    .header {
      flex-direction: column;
      gap: 1rem;
    }

    .header-actions {
      width: 100%;
    }
  }
</style>

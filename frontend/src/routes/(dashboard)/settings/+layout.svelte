<script lang="ts">
  import { page } from '$app/state'

  let { children } = $props()

  const tabs = [
    { href: '/settings/account', label: 'Account' },
    { href: '/settings', label: 'Calendar' }
  ]

  function isActive(href: string): boolean {
    if (href === '/settings') {
      return page.url.pathname === '/settings'
    }
    return page.url.pathname.startsWith(href)
  }
</script>

<div class="settings-tabs">
  {#each tabs as tab}
    <a href={tab.href} class="settings-tab" class:active={isActive(tab.href)}>
      {tab.label}
    </a>
  {/each}
</div>

{@render children()}

<style>
  .settings-tabs {
    display: flex;
    gap: 0.25rem;
    margin-bottom: 1.5rem;
    border-bottom: 2px solid var(--color-border);
    padding-bottom: 0;
  }

  .settings-tab {
    padding: 0.6rem 1.25rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text-secondary);
    text-decoration: none;
    border-bottom: 2px solid transparent;
    margin-bottom: -2px;
    transition: all 0.15s ease;
  }

  .settings-tab:hover {
    color: var(--color-text-primary);
  }

  .settings-tab.active {
    color: var(--color-text-primary);
    border-bottom-color: var(--color-accent-sky);
  }
</style>


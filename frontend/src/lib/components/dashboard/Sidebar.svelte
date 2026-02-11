<script lang="ts">
  import { page } from '$app/state'
  import { getAuth } from '$lib/stores/auth.svelte'
  import { DropdownMenu } from 'bits-ui'

  const auth = getAuth()

  // Navigation items grouped by section
  type NavItem = {
    href: string
    label: string
    icon: string
    badge?: string | number | null
  }

  const menuItems: NavItem[] = [
    { href: '/', label: 'Dashboard', icon: 'dashboard' },
    { href: '/booking-links', label: 'Booking Links', icon: 'link' },
    { href: '/polls', label: 'Polls', icon: 'poll' }
  ]

  const settingsItems: NavItem[] = [
    { href: '/settings/account', label: 'Account', icon: 'user' },
    { href: '/settings', label: 'Calendar', icon: 'calendar' }
  ]

  function isActive(href: string): boolean {
    if (href === '/') {
      return page.url.pathname === '/'
    }
    if (href === '/settings') {
      return page.url.pathname === '/settings'
    }
    return page.url.pathname.startsWith(href)
  }

  function handleLogout() {
    auth.logout()
  }

  // Get user initials from email
  function getUserInitials(): string {
    const email = auth.user?.email ?? ''
    if (!email) return '?'
    const name = email.split('@')[0]
    const parts = name.split(/[._-]/)
    if (parts.length >= 2) {
      return (parts[0][0] + parts[1][0]).toUpperCase()
    }
    return name.substring(0, 2).toUpperCase()
  }

  // Get display name from email
  function getUserDisplayName(): string {
    const email = auth.user?.email ?? ''
    if (!email) return 'User'
    const name = email.split('@')[0]
    return name.split(/[._-]/).map(p => p.charAt(0).toUpperCase() + p.slice(1)).join(' ')
  }
</script>

<aside class="sidebar">
  <!-- Logo -->
  <a href="/" class="logo">
    <div class="logo-icon">
      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
      </svg>
    </div>
    <span class="logo-text">Meet Mesh</span>
  </a>

  <!-- Navigation -->
  <nav class="nav">
    {#each menuItems as item}
      <a href={item.href} class="nav-item" class:active={isActive(item.href)}>
        {#if item.icon === 'dashboard'}
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/>
          </svg>
        {:else if item.icon === 'link'}
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
          </svg>
        {:else if item.icon === 'poll'}
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
          </svg>
        {/if}
        {item.label}
        {#if item.badge}
          <span class="nav-badge">{item.badge}</span>
        {/if}
      </a>
    {/each}

    <div class="nav-section-label">Settings</div>
    {#each settingsItems as item}
      <a href={item.href} class="nav-item" class:active={isActive(item.href)}>
        {#if item.icon === 'user'}
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
          </svg>
        {:else if item.icon === 'calendar'}
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        {/if}
        {item.label}
      </a>
    {/each}
  </nav>

  <!-- User section -->
  <div class="sidebar-footer">
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
          <button {...props} class="user-card">
            <div class="user-avatar">{getUserInitials()}</div>
            <div class="user-info">
              <div class="user-name">{getUserDisplayName()}</div>
              <div class="user-email">{auth.user?.email ?? ''}</div>
            </div>
            <svg class="user-chevron" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
            </svg>
          </button>
        {/snippet}
      </DropdownMenu.Trigger>

      <DropdownMenu.Portal>
        <DropdownMenu.Content class="user-menu-content" side="top" sideOffset={8} align="start">
          <DropdownMenu.Item class="user-menu-item" onSelect={handleLogout}>
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Log out
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  </div>
</aside>

<style>
  .sidebar {
    width: var(--spacing-sidebar);
    background: var(--color-bg-secondary);
    border-right: 2px solid var(--color-border);
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
    transition: transform 0.15s ease, background 0.15s ease;
  }

  /* Logo */
  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding-bottom: 1.5rem;
    margin-bottom: 1.5rem;
    border-bottom: 1px solid var(--color-border);
    text-decoration: none;
    cursor: pointer;
  }

  .logo-icon {
    width: 40px;
    height: 40px;
    background: var(--color-accent-sky);
    border: 2px solid var(--color-border);
    border-radius: var(--radius-brutalist);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-brutalist-sm);
    transition: box-shadow 0.15s ease, border-color 0.15s ease;
  }

  :global([data-theme="dark"]) .logo-icon {
    border-color: rgba(14, 165, 233, 0.3);
    box-shadow: 0 0 20px rgba(14, 165, 233, 0.3);
  }

  .logo-icon svg {
    width: 22px;
    height: 22px;
    color: white;
  }

  .logo-text {
    font-size: 1.25rem;
    font-weight: 800;
    color: var(--color-text-primary);
    letter-spacing: -0.03em;
  }

  /* Navigation */
  .nav {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.7rem 0.75rem;
    border-radius: var(--radius-brutalist);
    color: var(--color-text-secondary);
    text-decoration: none;
    font-size: 0.9rem;
    font-weight: 600;
    transition: all 0.15s ease;
    border: 2px solid transparent;
    cursor: pointer;
  }

  .nav-item:hover {
    background: var(--color-bg-tertiary);
    color: var(--color-text-primary);
  }

  .nav-item.active {
    background: var(--color-accent-sky);
    color: white;
    border-color: var(--color-border);
    box-shadow: var(--shadow-brutalist-sm);
  }

  :global([data-theme="dark"]) .nav-item.active {
    border-color: transparent;
    box-shadow: 0 0 16px rgba(14, 165, 233, 0.3);
  }

  .nav-item svg {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
  }

  .nav-badge {
    margin-left: auto;
    background: var(--color-accent-rose);
    color: white;
    font-size: 0.7rem;
    font-weight: 700;
    padding: 0.15rem 0.5rem;
    border-radius: 99px;
    border: 1px solid var(--color-border);
  }

  .nav-section-label {
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
    padding: 1rem 0.75rem 0.35rem;
  }

  /* Sidebar footer */
  .sidebar-footer {
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .user-card {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem;
    border-radius: var(--radius-brutalist);
    transition: background 0.15s ease;
    cursor: pointer;
    width: 100%;
    border: none;
    background: transparent;
    text-align: left;
  }

  .user-card:hover {
    background: var(--color-bg-tertiary);
  }

  .user-avatar {
    width: 36px;
    height: 36px;
    border-radius: var(--radius-brutalist);
    background: var(--color-accent-amber);
    border: 2px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: 700;
    font-size: 0.8rem;
    flex-shrink: 0;
  }

  .user-info {
    min-width: 0;
    flex: 1;
  }

  .user-name {
    font-weight: 600;
    font-size: 0.85rem;
    color: var(--color-text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-email {
    font-size: 0.7rem;
    color: var(--color-text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-chevron {
    width: 16px;
    height: 16px;
    color: var(--color-text-muted);
    flex-shrink: 0;
    transition: color 0.15s ease;
  }

  .user-card:hover .user-chevron {
    color: var(--color-text-primary);
  }

  :global(.user-menu-content) {
    min-width: 200px;
    background: var(--color-bg-secondary);
    border: 2px solid var(--color-border);
    border-radius: var(--radius-brutalist);
    padding: 0.35rem;
    box-shadow: var(--shadow-brutalist);
    z-index: 200;
  }

  :global(.user-menu-content[data-state="open"]) {
    animation: menuIn 0.15s ease;
  }

  :global(.user-menu-content[data-state="closed"]) {
    animation: menuOut 0.1s ease;
  }

  @keyframes menuIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @keyframes menuOut {
    from { opacity: 1; transform: translateY(0); }
    to { opacity: 0; transform: translateY(4px); }
  }

  :global(.user-menu-item) {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.6rem 0.75rem;
    border-radius: calc(var(--radius-brutalist) - 2px);
    color: var(--color-text-secondary);
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s ease;
    border: none;
    background: transparent;
    width: 100%;
  }

  :global(.user-menu-item:hover),
  :global(.user-menu-item[data-highlighted]) {
    background: var(--color-bg-tertiary);
    color: var(--color-text-primary);
  }

  :global(.user-menu-item svg) {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
  }

  /* Responsive - hide on mobile (below 640px), MobileMenu handles that */
  @media (max-width: 640px) {
    .sidebar {
      display: none;
    }
  }
</style>

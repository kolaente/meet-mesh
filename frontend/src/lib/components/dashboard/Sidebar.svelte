<script lang="ts">
  import { page } from '$app/state'
  import { getAuth } from '$lib/stores/auth.svelte'

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
    { href: '/settings', label: 'Calendar', icon: 'calendar' }
  ]

  function isActive(href: string): boolean {
    if (href === '/') {
      return page.url.pathname === '/'
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

    {#each settingsItems as item}
      <a href={item.href} class="nav-item" class:active={isActive(item.href)}>
        {#if item.icon === 'calendar'}
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
    <button class="user-card" onclick={handleLogout} title="Click to logout">
      <div class="user-avatar">{getUserInitials()}</div>
      <div class="user-info">
        <div class="user-name">{getUserDisplayName()}</div>
        <div class="user-email">{auth.user?.email ?? ''}</div>
      </div>
    </button>
  </div>
</aside>

<style>
  .sidebar {
    width: var(--sidebar-width);
    background: var(--bg-secondary);
    border-right: var(--border);
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
    transition: transform var(--transition), background var(--transition);
  }

  /* Logo */
  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding-bottom: 1.5rem;
    margin-bottom: 1.5rem;
    border-bottom: var(--border-light);
    text-decoration: none;
  }

  .logo-icon {
    width: 40px;
    height: 40px;
    background: var(--cyan);
    border: var(--border);
    border-radius: var(--radius);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-sm);
    transition: box-shadow var(--transition), border-color var(--transition);
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
    color: var(--text-primary);
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
    border-radius: var(--radius);
    color: var(--text-secondary);
    text-decoration: none;
    font-size: 0.9rem;
    font-weight: 600;
    transition: all var(--transition);
    border: 2px solid transparent;
  }

  .nav-item:hover {
    background: var(--bg-tertiary);
    color: var(--text-primary);
  }

  .nav-item.active {
    background: var(--cyan);
    color: white;
    border-color: var(--border-color);
    box-shadow: var(--shadow-sm);
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
    background: var(--pink);
    color: white;
    font-size: 0.7rem;
    font-weight: 700;
    padding: 0.15rem 0.5rem;
    border-radius: 99px;
    border: 1px solid var(--border-color);
  }

  /* Sidebar footer */
  .sidebar-footer {
    padding-top: 1rem;
    border-top: var(--border-light);
  }

  .user-card {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem;
    border-radius: var(--radius);
    transition: background var(--transition);
    cursor: pointer;
    width: 100%;
    border: none;
    background: transparent;
    text-align: left;
  }

  .user-card:hover {
    background: var(--bg-tertiary);
  }

  .user-avatar {
    width: 36px;
    height: 36px;
    border-radius: var(--radius);
    background: var(--orange);
    border: var(--border);
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
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-email {
    font-size: 0.7rem;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  /* Responsive - hide on mobile (below 640px), MobileMenu handles that */
  @media (max-width: 640px) {
    .sidebar {
      display: none;
    }
  }
</style>

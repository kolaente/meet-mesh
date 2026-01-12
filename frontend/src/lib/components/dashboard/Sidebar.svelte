<script lang="ts">
  import { page } from '$app/state'
  import { getAuth } from '$lib/stores/auth.svelte'

  const auth = getAuth()

  export const navItems = [
    { href: '/', label: 'Dashboard', icon: 'dashboard' },
    { href: '/booking-links', label: 'Booking Links', icon: 'calendar' },
    { href: '/polls', label: 'Polls', icon: 'poll' },
    { href: '/settings', label: 'Settings', icon: 'settings' }
  ] as const

  function isActive(href: string): boolean {
    if (href === '/') {
      return page.url.pathname === '/'
    }
    return page.url.pathname.startsWith(href)
  }

  function handleLogout() {
    auth.logout()
  }

  // Tooltip state for collapsed view
  let hoveredItem = $state<string | null>(null)
</script>

<!-- Desktop: Full sidebar (lg:), Tablet: Icons only (sm: to lg:), Mobile: Hidden (below sm:) -->
<aside class="hidden sm:flex sm:w-16 lg:w-60 flex-col bg-white border-r border-gray-200 min-h-screen transition-all duration-200">
  <!-- Logo -->
  <div class="h-16 flex items-center px-3 lg:px-6 border-b border-gray-200">
    <span class="hidden lg:block text-xl font-semibold text-gray-900">Meet Mesh</span>
    <span class="lg:hidden text-xl font-semibold text-gray-900">M</span>
  </div>

  <!-- Navigation -->
  <nav class="flex-1 px-2 lg:px-3 py-4">
    <ul class="space-y-1">
      {#each navItems as item}
        <li class="relative">
          <a
            href={item.href}
            onmouseenter={() => (hoveredItem = item.label)}
            onmouseleave={() => (hoveredItem = null)}
            class="flex items-center justify-center lg:justify-start gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors {isActive(item.href)
              ? 'bg-indigo-50 text-indigo-700'
              : 'text-gray-700 hover:bg-gray-100'}"
          >
            {#if item.icon === 'dashboard'}
              <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
              </svg>
            {:else if item.icon === 'calendar'}
              <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            {:else if item.icon === 'poll'}
              <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            {:else if item.icon === 'settings'}
              <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            {/if}
            <span class="hidden lg:inline">{item.label}</span>
          </a>
          <!-- Tooltip for collapsed sidebar (tablet view) -->
          {#if hoveredItem === item.label}
            <div class="hidden sm:block lg:hidden absolute left-full ml-2 top-1/2 -translate-y-1/2 z-50 px-2 py-1 text-xs font-medium text-white bg-gray-900 rounded shadow-lg whitespace-nowrap">
              {item.label}
            </div>
          {/if}
        </li>
      {/each}
    </ul>
  </nav>

  <!-- User section -->
  <div class="border-t border-gray-200 p-2 lg:p-4">
    <p class="hidden lg:block text-sm text-gray-600 truncate mb-2">{auth.user?.email ?? ''}</p>
    <button
      onclick={handleLogout}
      class="flex items-center justify-center lg:justify-start w-full gap-2 text-sm text-gray-500 hover:text-gray-700 transition-colors p-2 lg:p-0"
      title="Logout"
    >
      <svg class="w-5 h-5 lg:hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
      </svg>
      <span class="hidden lg:inline">Logout</span>
    </button>
  </div>
</aside>

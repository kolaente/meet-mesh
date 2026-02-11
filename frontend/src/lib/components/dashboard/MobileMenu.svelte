<script lang="ts">
	import { page } from '$app/state';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { fly, fade } from 'svelte/transition';

	const auth = getAuth();

	interface Props {
		open?: boolean;
		onClose?: () => void;
	}

	let { open = $bindable(false), onClose }: Props = $props();

	const navItems = [
		{ href: '/', label: 'Dashboard', icon: 'dashboard' },
		{ href: '/booking-links', label: 'Booking Links', icon: 'link' },
		{ href: '/polls', label: 'Polls', icon: 'poll' },
		{ href: '/settings/account', label: 'Account', icon: 'user' },
		{ href: '/settings', label: 'Calendar', icon: 'calendar' }
	] as const;

	function isActive(href: string): boolean {
		if (href === '/') {
			return page.url.pathname === '/';
		}
		if (href === '/settings') {
			return page.url.pathname === '/settings';
		}
		return page.url.pathname.startsWith(href);
	}

	function handleLogout() {
		auth.logout();
		close();
	}

	function close() {
		open = false;
		onClose?.();
	}

	function handleNavClick() {
		close();
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			close();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- Mobile menu - only shown on mobile (below sm:) -->
<div class="sm:hidden">
	<!-- Hamburger button - neo-brutalist style -->
	<button
		type="button"
		onclick={() => (open = true)}
		class="mobile-menu-button"
		aria-label="Open menu"
	>
		<svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-width="2.5" d="M4 6h16M4 12h16M4 18h16" />
		</svg>
	</button>

	<!-- Overlay -->
	{#if open}
		<div
			class="overlay"
			transition:fade={{ duration: 200 }}
			onclick={close}
			onkeydown={(e) => e.key === 'Enter' && close()}
			role="button"
			tabindex="0"
			aria-label="Close menu"
		></div>
	{/if}

	<!-- Drawer -->
	{#if open}
		<aside class="drawer" transition:fly={{ x: -280, duration: 200 }}>
			<!-- Logo -->
			<a href="/" onclick={handleNavClick} class="logo">
				<div class="logo-icon">
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
						/>
					</svg>
				</div>
				<span class="logo-text">Meet Mesh</span>
			</a>

			<!-- Navigation -->
			<nav class="nav">
				{#each navItems as item}
					<a
						href={item.href}
						onclick={handleNavClick}
						class="nav-item"
						class:active={isActive(item.href)}
					>
						{#if item.icon === 'dashboard'}
							<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-width="2"
									d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
								/>
							</svg>
						{:else if item.icon === 'link'}
							<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
								/>
							</svg>
						{:else if item.icon === 'poll'}
							<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-width="2"
									d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
								/>
							</svg>
						{:else if item.icon === 'user'}
							<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
								/>
							</svg>
						{:else if item.icon === 'calendar'}
							<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
								/>
							</svg>
						{/if}
						{item.label}
					</a>
				{/each}
			</nav>

			<!-- User section -->
			<div class="sidebar-footer">
				<div class="user-card">
					<div class="user-avatar">
						{auth.user?.email?.charAt(0).toUpperCase() ?? 'U'}
					</div>
					<div class="user-info">
						<div class="user-email">{auth.user?.email ?? ''}</div>
					</div>
				</div>
				<button onclick={handleLogout} class="logout-button">
					<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-width="2"
							d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
						/>
					</svg>
					Logout
				</button>
			</div>
		</aside>
	{/if}
</div>

<style>
	/* CSS Variables matching the prototype */
	:root {
		--bg-primary: #f5f3f0;
		--bg-secondary: #ffffff;
		--bg-tertiary: #ebe8e4;
		--text-primary: #1a1a1a;
		--text-secondary: #555555;
		--text-muted: #888888;
		--border-color: #1a1a1a;
		--cyan: #0ea5e9;
		--orange: #f59e0b;
		--shadow: 3px 3px 0 var(--border-color);
		--shadow-sm: 2px 2px 0 var(--border-color);
		--border: 2px solid var(--border-color);
		--border-light: 1px solid var(--border-color);
		--radius: 8px;
		--transition: 0.15s ease;
	}

	/* Mobile hamburger button */
	.mobile-menu-button {
		position: fixed;
		top: 1rem;
		left: 1rem;
		z-index: 200;
		width: 44px;
		height: 44px;
		border-radius: var(--radius);
		border: var(--border);
		background: var(--bg-secondary);
		color: var(--text-primary);
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: var(--shadow);
		transition: all var(--transition);
	}

	.mobile-menu-button:hover {
		transform: translate(-1px, -1px);
		box-shadow: 4px 4px 0 var(--border-color);
	}

	.mobile-menu-button:active {
		transform: translate(1px, 1px);
		box-shadow: 1px 1px 0 var(--border-color);
	}

	/* Overlay */
	.overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 50;
	}

	/* Drawer */
	.drawer {
		position: fixed;
		inset-block: 0;
		left: 0;
		z-index: 100;
		width: 260px;
		background: var(--bg-secondary);
		border-right: var(--border);
		padding: 1.25rem;
		display: flex;
		flex-direction: column;
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
		cursor: pointer;
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
	}

	.logo-icon svg {
		width: 22px;
		height: 22px;
		color: white;
		stroke-width: 2.5;
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
		cursor: pointer;
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

	.nav-item svg {
		width: 18px;
		height: 18px;
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
		margin-bottom: 0.75rem;
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
	}

	.user-info {
		flex: 1;
		min-width: 0;
	}

	.user-email {
		font-size: 0.75rem;
		color: var(--text-muted);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.logout-button {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		width: 100%;
		padding: 0.7rem 0.75rem;
		border-radius: var(--radius);
		color: var(--text-secondary);
		background: transparent;
		border: 2px solid transparent;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		transition: all var(--transition);
	}

	.logout-button:hover {
		background: var(--bg-tertiary);
		color: var(--text-primary);
	}

	.logout-button svg {
		width: 18px;
		height: 18px;
	}

	/* Dark mode support */
	:global([data-theme='dark']) .mobile-menu-button {
		--bg-secondary: #0a0a0a;
		--text-primary: #f5f5f5;
		--border-color: #000000;
	}

	:global([data-theme='dark']) .drawer {
		--bg-secondary: #0a0a0a;
		--bg-tertiary: #111111;
		--text-primary: #f5f5f5;
		--text-secondary: #a0a0a0;
		--text-muted: #606060;
		--border-color: #000000;
	}

	:global([data-theme='dark']) .logo-icon {
		border-color: rgba(14, 165, 233, 0.3);
		box-shadow: 0 0 20px rgba(14, 165, 233, 0.3);
	}

	:global([data-theme='dark']) .nav-item.active {
		border-color: transparent;
		box-shadow: 0 0 16px rgba(14, 165, 233, 0.3);
	}
</style>

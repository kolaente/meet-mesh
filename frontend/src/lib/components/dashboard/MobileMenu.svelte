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
		{ href: '/booking-links', label: 'Booking Links', icon: 'calendar' },
		{ href: '/polls', label: 'Polls', icon: 'poll' },
		{ href: '/settings', label: 'Settings', icon: 'settings' }
	] as const;

	function isActive(href: string): boolean {
		if (href === '/') {
			return page.url.pathname === '/';
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
	<!-- Hamburger button -->
	<button
		type="button"
		onclick={() => (open = true)}
		class="fixed top-4 left-4 z-40 p-2 rounded-lg bg-white border border-gray-200 shadow-sm hover:bg-gray-50 transition-colors"
		aria-label="Open menu"
	>
		<svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M4 6h16M4 12h16M4 18h16"
			/>
		</svg>
	</button>

	<!-- Overlay -->
	{#if open}
		<div
			class="fixed inset-0 z-40 bg-black/50"
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
		<aside
			class="fixed inset-y-0 left-0 z-50 w-72 bg-white shadow-xl flex flex-col"
			transition:fly={{ x: -288, duration: 200 }}
		>
			<!-- Header -->
			<div class="h-16 flex items-center justify-between px-6 border-b border-gray-200">
				<span class="text-xl font-semibold text-gray-900">Meet Mesh</span>
				<button
					type="button"
					onclick={close}
					class="p-2 -mr-2 rounded-lg text-gray-500 hover:bg-gray-100 transition-colors"
					aria-label="Close menu"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			</div>

			<!-- Navigation -->
			<nav class="flex-1 px-3 py-4 overflow-y-auto">
				<ul class="space-y-1">
					{#each navItems as item}
						<li>
							<a
								href={item.href}
								onclick={handleNavClick}
								class="flex items-center gap-3 px-3 py-3 rounded-lg text-base font-medium transition-colors {isActive(
									item.href
								)
									? 'bg-indigo-50 text-indigo-700'
									: 'text-gray-700 hover:bg-gray-100'}"
							>
								{#if item.icon === 'dashboard'}
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
										/>
									</svg>
								{:else if item.icon === 'calendar'}
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
										/>
									</svg>
								{:else if item.icon === 'poll'}
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
										/>
									</svg>
								{:else if item.icon === 'settings'}
									<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
										/>
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
										/>
									</svg>
								{/if}
								{item.label}
							</a>
						</li>
					{/each}
				</ul>
			</nav>

			<!-- User section -->
			<div class="border-t border-gray-200 p-4">
				<p class="text-sm text-gray-600 truncate mb-3">{auth.user?.email ?? ''}</p>
				<button
					onclick={handleLogout}
					class="flex items-center gap-2 w-full px-3 py-3 text-sm font-medium text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
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

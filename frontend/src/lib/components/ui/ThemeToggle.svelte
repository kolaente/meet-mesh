<script lang="ts">
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';

  let theme = $state<'light' | 'dark'>('light');

  onMount(() => {
    const saved = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    theme = (saved || (prefersDark ? 'dark' : 'light')) as 'light' | 'dark';
    document.documentElement.setAttribute('data-theme', theme);
  });

  function toggleTheme() {
    theme = theme === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }
</script>

<button
  onclick={toggleTheme}
  class="theme-toggle"
  title="Toggle theme"
  aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
>
  {#if theme === 'light'}
    <!-- Sun icon -->
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <circle cx="12" cy="12" r="5" stroke-width="2"/>
      <path stroke-linecap="round" stroke-width="2" d="M12 1v2m0 18v2M4.22 4.22l1.42 1.42m12.72 12.72l1.42 1.42M1 12h2m18 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
    </svg>
  {:else}
    <!-- Moon icon -->
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"/>
    </svg>
  {/if}
</button>

<style>
  .theme-toggle {
    width: 40px;
    height: 40px;
    border-radius: var(--radius);
    border: var(--border);
    background: var(--bg-secondary);
    color: var(--text-primary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition);
    box-shadow: var(--shadow-sm);
  }

  .theme-toggle:hover {
    transform: translate(-1px, -1px);
    box-shadow: var(--shadow);
  }

  .theme-toggle:active {
    transform: translate(1px, 1px);
    box-shadow: none;
  }

  .theme-toggle svg {
    width: 18px;
    height: 18px;
  }
</style>

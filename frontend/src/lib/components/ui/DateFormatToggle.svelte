<script lang="ts">
  import { getDateFormat } from '$lib/stores/dateFormat.svelte';

  let expanded = $state(false);
  const dateFormat = getDateFormat();

  function toggleTimeFormat() {
    dateFormat.setTimeFormat(dateFormat.timeFormat === '12h' ? '24h' : '12h');
  }

  function toggleWeekStart() {
    dateFormat.setWeekStartDay(dateFormat.weekStartDay === 'sunday' ? 'monday' : 'sunday');
  }
</script>

<div class="date-format-toggle">
  <button
    type="button"
    class="toggle-button"
    onclick={() => expanded = !expanded}
    aria-expanded={expanded}
    aria-label="Date format settings"
  >
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  </button>

  {#if expanded}
    <div class="settings-panel">
      <div class="setting-row">
        <span class="setting-label">Time</span>
        <button
          type="button"
          class="setting-value"
          onclick={toggleTimeFormat}
        >
          {dateFormat.timeFormat === '12h' ? '12h' : '24h'}
        </button>
      </div>
      <div class="setting-row">
        <span class="setting-label">Week starts</span>
        <button
          type="button"
          class="setting-value"
          onclick={toggleWeekStart}
        >
          {dateFormat.weekStartDay === 'sunday' ? 'Sun' : 'Mon'}
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .date-format-toggle {
    position: relative;
  }

  .toggle-button {
    width: 36px;
    height: 36px;
    border-radius: var(--radius);
    border: var(--border);
    background: var(--bg-secondary);
    color: var(--text-secondary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition);
    box-shadow: var(--shadow-sm);
  }

  .toggle-button:hover {
    color: var(--text-primary);
    background: var(--bg-tertiary);
  }

  .settings-panel {
    position: absolute;
    bottom: 100%;
    right: 0;
    margin-bottom: 0.5rem;
    padding: 0.5rem;
    background: var(--bg-secondary);
    border: var(--border);
    border-radius: var(--radius);
    box-shadow: var(--shadow);
    min-width: 140px;
  }

  .setting-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.25rem 0;
  }

  .setting-label {
    font-size: 0.75rem;
    color: var(--text-secondary);
  }

  .setting-value {
    font-size: 0.75rem;
    font-weight: 600;
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    background: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    cursor: pointer;
    transition: all var(--transition);
  }

  .setting-value:hover {
    background: var(--sky);
    color: white;
    border-color: var(--sky);
  }
</style>

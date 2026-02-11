<script lang="ts">
  import { Spinner } from '$lib/components/ui';

  interface Props {
    /** Current avatar URL, empty string or undefined if no avatar */
    avatarUrl?: string;
    /** Initials to show as fallback when no avatar is set */
    initials: string;
    /** Called when avatar is uploaded or deleted. Receives the new avatar_url (empty string if deleted). */
    onchange?: (avatarUrl: string) => void;
  }

  let { avatarUrl = '', initials, onchange }: Props = $props();

  let uploading = $state(false);
  let deleting = $state(false);
  let dragOver = $state(false);
  let error = $state('');
  let fileInput = $state<HTMLInputElement | null>(null);

  const MAX_SIZE = 2 * 1024 * 1024; // 2MB
  const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/webp'];

  function handleFileSelect(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files?.length) {
      uploadFile(input.files[0]);
    }
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    dragOver = false;
    if (e.dataTransfer?.files.length) {
      uploadFile(e.dataTransfer.files[0]);
    }
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    dragOver = true;
  }

  function handleDragLeave() {
    dragOver = false;
  }

  async function uploadFile(file: File) {
    error = '';

    // Client-side validation
    if (!ALLOWED_TYPES.includes(file.type)) {
      error = 'Invalid file type. Use JPEG, PNG, or WebP.';
      return;
    }

    if (file.size > MAX_SIZE) {
      error = 'File too large. Maximum size is 2MB.';
      return;
    }

    uploading = true;
    try {
      const formData = new FormData();
      formData.append('avatar', file);

      const response = await fetch('/api/avatars/', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const text = await response.text();
        error = text || 'Upload failed';
        return;
      }

      const data = await response.json();
      onchange?.(data.avatar_url);
    } catch (err) {
      error = 'Upload failed. Please try again.';
    } finally {
      uploading = false;
      // Reset file input
      if (fileInput) fileInput.value = '';
    }
  }

  async function deleteAvatar() {
    error = '';
    deleting = true;
    try {
      const response = await fetch('/api/avatars/', {
        method: 'DELETE',
      });

      if (!response.ok) {
        error = 'Failed to remove avatar';
        return;
      }

      onchange?.('');
    } catch (err) {
      error = 'Failed to remove avatar. Please try again.';
    } finally {
      deleting = false;
    }
  }
</script>

<div class="avatar-upload">
  <!-- Current Avatar / Drop Zone -->
  <button
    type="button"
    class="avatar-dropzone"
    class:drag-over={dragOver}
    class:has-avatar={!!avatarUrl}
    ondrop={handleDrop}
    ondragover={handleDragOver}
    ondragleave={handleDragLeave}
    onclick={() => fileInput?.click()}
    disabled={uploading}
  >
    {#if uploading}
      <div class="avatar-preview">
        <Spinner size="lg" />
      </div>
    {:else if avatarUrl}
      <img src={avatarUrl} alt="Avatar" class="avatar-preview" />
      <div class="avatar-overlay">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="overlay-icon">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        <span>Change</span>
      </div>
    {:else}
      <div class="avatar-placeholder">
        <span class="avatar-initials">{initials}</span>
      </div>
      <div class="avatar-overlay">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="overlay-icon">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        <span>Upload</span>
      </div>
    {/if}
  </button>

  <input
    bind:this={fileInput}
    type="file"
    accept="image/jpeg,image/png,image/webp"
    onchange={handleFileSelect}
    class="hidden-input"
  />

  <div class="avatar-actions">
    <p class="avatar-hint">Click or drag & drop. JPEG, PNG, or WebP. Max 2MB.</p>
    {#if avatarUrl}
      <button
        type="button"
        class="remove-btn"
        onclick={deleteAvatar}
        disabled={deleting}
      >
        {#if deleting}
          Removing...
        {:else}
          Remove avatar
        {/if}
      </button>
    {/if}
  </div>

  {#if error}
    <p class="avatar-error">{error}</p>
  {/if}
</div>

<style>
  .avatar-upload {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .avatar-dropzone {
    position: relative;
    width: 96px;
    height: 96px;
    border-radius: 50%;
    overflow: hidden;
    cursor: pointer;
    border: 2px dashed var(--color-border);
    background: var(--color-bg-tertiary);
    padding: 0;
    transition: border-color 0.2s;
  }

  .avatar-dropzone:hover,
  .avatar-dropzone.drag-over {
    border-color: var(--color-cyan);
  }

  .avatar-dropzone.has-avatar {
    border-style: solid;
  }

  .avatar-preview {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--color-orange);
  }

  .avatar-initials {
    color: white;
    font-weight: 700;
    font-size: 1.5rem;
  }

  .avatar-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.5);
    color: white;
    opacity: 0;
    transition: opacity 0.2s;
    font-size: 0.75rem;
    font-weight: 600;
  }

  .avatar-dropzone:hover .avatar-overlay {
    opacity: 1;
  }

  .overlay-icon {
    width: 20px;
    height: 20px;
    margin-bottom: 2px;
  }

  .hidden-input {
    display: none;
  }

  .avatar-actions {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .avatar-hint {
    font-size: 0.75rem;
    color: var(--color-text-muted);
  }

  .remove-btn {
    background: none;
    border: none;
    color: var(--color-pink);
    font-size: 0.8rem;
    font-weight: 600;
    cursor: pointer;
    padding: 0;
    text-align: left;
  }

  .remove-btn:hover {
    text-decoration: underline;
  }

  .remove-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .avatar-error {
    color: var(--color-pink);
    font-size: 0.8rem;
  }
</style>



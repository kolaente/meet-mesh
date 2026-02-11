---
# meet-mesh-4rqq
title: Rewrite app.css with @theme block and dark mode overrides
status: completed
type: task
priority: normal
created_at: 2026-02-11T18:37:32Z
updated_at: 2026-02-11T18:55:40Z
parent: meet-mesh-gpjw
---

# Rewrite app.css with @theme block and dark mode overrides

This is the foundational change. All other tasks depend on this one landing first.

## File to Modify

- `frontend/src/app.css`

## Overview

Tailwind v4's `@theme` directive registers design tokens as CSS custom properties at `:root` level and auto-generates utility classes based on namespace. For colors that change between light/dark mode, define the light values in `@theme` and override in a dark mode selector.

## Replacement Content

Replace the entire `app.css` with:

```css
@import 'tailwindcss';

@custom-variant dark (&:where([data-theme="dark"], [data-theme="dark"] *));

@theme {
  /* Background colors */
  --color-bg-primary: #f5f3f0;
  --color-bg-secondary: #ffffff;
  --color-bg-tertiary: #ebe8e4;

  /* Text colors */
  --color-text-primary: #1a1a1a;
  --color-text-secondary: #555555;
  --color-text-muted: #888888;
  --color-text-tertiary: #666666;

  /* Border color */
  --color-border: #1a1a1a;

  /* Accent colors */
  --color-accent-amber: #f59e0b;
  --color-accent-emerald: #10b981;
  --color-accent-sky: #0ea5e9;
  --color-accent-sky-hover: #0284c7;
  --color-accent-violet: #8b5cf6;
  --color-accent-rose: #f43f5e;

  /* Border radius */
  --radius-brutalist-sm: 0.375rem;
  --radius-brutalist-md: 0.5rem;
  --radius-brutalist: 8px;
  --radius-brutalist-lg: 12px;

  /* Shadows */
  --shadow-brutalist: 3px 3px 0 #1a1a1a;
  --shadow-brutalist-hover: 5px 5px 0 #1a1a1a;
  --shadow-brutalist-sm: 2px 2px 0 #1a1a1a;
  --shadow-DEFAULT: 0 4px 6px -1px rgb(0 0 0 / 0.1);

  /* Spacing */
  --spacing-sidebar: 250px;
}

/* Dark mode overrides - these override the @theme-generated CSS vars */
@layer base {
  [data-theme="dark"] {
    --color-bg-primary: #1a1a1a;
    --color-bg-secondary: #0a0a0a;
    --color-bg-tertiary: #111111;
    --color-text-primary: #f5f5f5;
    --color-text-secondary: #a0a0a0;
    --color-text-muted: #606060;
    --color-text-tertiary: #505050;
    --color-border: #333333;
    --shadow-brutalist: 3px 3px 0 #000000;
    --shadow-brutalist-hover: 5px 5px 0 #000000;
    --shadow-brutalist-sm: 2px 2px 0 #000000;
  }
}

body {
  font-family: 'Outfit Variable', -apple-system, sans-serif;
  background: var(--color-bg-primary);
  color: var(--color-text-primary);
  min-height: 100vh;
  transition: background 0.15s ease, color 0.15s ease;
}

button {
  cursor: pointer;
}
```

## Generated Utility Classes

This `@theme` block auto-generates these Tailwind utility classes:

- `bg-bg-primary`, `bg-bg-secondary`, `bg-bg-tertiary`
- `text-text-primary`, `text-text-secondary`, `text-text-muted`, `text-text-tertiary`
- `border-border` (from `--color-border`)
- `bg-accent-sky`, `text-accent-sky`, `border-accent-sky`, etc.
- `rounded-brutalist-sm`, `rounded-brutalist-md`, `rounded-brutalist`, `rounded-brutalist-lg`
- `shadow-brutalist`, `shadow-brutalist-hover`, `shadow-brutalist-sm`

## Key Decisions

- **`--text-tertiary`** was used in components but never defined. Defined here as `#666666` (light) / `#505050` (dark).
- **`--sky-hover`** was used but never defined. Defined as `--color-accent-sky-hover: #0284c7`.
- **Radii** use `brutalist-*` prefix to avoid collision with Tailwind built-in `rounded-sm`, `rounded-md`, `rounded-lg`.

## Steps

1. Rewrite `app.css` as shown above.
2. Run: `cd frontend && pnpm build`
3. Commit: `git commit -m "feat: add @theme block with Tailwind v4 design tokens"`

## Summary of Changes

- Added @theme block with Tailwind v4 design tokens
- Created proper namespaced CSS variables (--color-*, --radius-*, --shadow-*)
- Added dark mode overrides in @layer base
- Added legacy CSS variable aliases for backwards compatibility with existing component code
- Build passes successfully

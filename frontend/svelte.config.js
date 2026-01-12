import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      fallback: '200.html'  // SPA fallback for client-side routing
    }),
    alias: {
      $components: 'src/lib/components'
    }
  }
};

export default config;

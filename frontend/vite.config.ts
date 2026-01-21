import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { readFileSync, existsSync } from 'fs';
import { parse } from 'yaml';

function getApiTarget(): string {
	const configPath = '../config.yaml';
	const defaultTarget = 'http://localhost:8080';

	if (!existsSync(configPath)) {
		console.warn(`Config file not found at ${configPath}, using default: ${defaultTarget}`);
		return defaultTarget;
	}

	try {
		const configContent = readFileSync(configPath, 'utf-8');
		const config = parse(configContent);
		const baseUrl = config?.server?.base_url;
		if (baseUrl) {
			return baseUrl;
		}
		console.warn(`No server.base_url in config, using default: ${defaultTarget}`);
		return defaultTarget;
	} catch (error) {
		console.warn(`Failed to parse config file: ${error}, using default: ${defaultTarget}`);
		return defaultTarget;
	}
}

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: getApiTarget(),
				changeOrigin: true
			}
		}
	}
});

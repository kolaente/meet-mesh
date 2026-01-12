<script lang="ts">
	import Input from '../ui/Input.svelte';
	import Button from '../ui/Button.svelte';

	interface FormData {
		name?: string;
		email?: string;
	}

	interface Props {
		onSubmit: (data: FormData) => void;
		loading?: boolean;
		requireEmail?: boolean;
		class?: string;
	}

	let {
		onSubmit,
		loading = false,
		requireEmail = false,
		class: className = ''
	}: Props = $props();

	let name = $state('');
	let email = $state('');

	function handleSubmit(event: Event) {
		event.preventDefault();
		onSubmit({
			name: name || undefined,
			email: email || undefined
		});
	}
</script>

<form onsubmit={handleSubmit} class="space-y-4 {className}">
	<div class="border-t border-gray-200 pt-4">
		<h3 class="text-base font-medium text-gray-900 mb-4">Your Information</h3>

		<!-- Name field -->
		<Input
			name="name"
			label="Name"
			type="text"
			bind:value={name}
			placeholder="Your name (optional)"
		/>

		<!-- Email field -->
		<div class="mt-4">
			<Input
				name="email"
				label="Email"
				type="email"
				bind:value={email}
				placeholder="you@example.com"
				required={requireEmail}
				description={requireEmail ? undefined : 'Recommended for updates about this poll'}
			/>
		</div>
	</div>

	<div class="pt-2">
		<Button type="submit" {loading} class="w-full">
			{loading ? 'Submitting...' : 'Submit Vote'}
		</Button>
	</div>
</form>

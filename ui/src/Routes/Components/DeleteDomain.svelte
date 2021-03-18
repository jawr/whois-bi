<script>
	import { navigate } from 'svelte-routing'
	import { fetchJSON } from '../../fetchJSON'

	export let name = ''

	let confirm = ''
	let error = ''
	let started = false

	const deleteDomain = async () => {
		if (confirm === name) {
			try {
				const response = await fetchJSON(`/api/user/domain/${name}`, {method: 'DELETE'})
				navigate('/')
			} catch (err) {
				error = err.message
			}
		} else {
			started = true
		}
	}
</script>

<div class="cf">
	<input
		class="f5 fr pv1 h2 tc dib bb bt-0 bl-0 br-0 bw2 b--dark-red bg-animate bg-red hover-bg-red white pointer br2 grow w-10"
		type="submit"
		value="Delete"
		on:click|preventDefault={deleteDomain}
	/>
	{#if started}
		<input
			class="mr4 f5 h2 button-reset fr pv1 dib bl bt bb bw1 b--light-red tc br--left br2 w-30"
			type="text"
			placeholder={`type ${name} to confirm`}
			bind:value={confirm}
		/>
	{/if}
</div>

{#if error.length > 0}
<div class="cf">
	<p class="mt3 pr3 fr light-red">{error}</p>
</div>
{/if}

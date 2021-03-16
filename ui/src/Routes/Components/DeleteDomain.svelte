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
		class="mt6 f5 fr pv1 h2 tc dib ba bw1 b--light-red bg-animate bg-light-red hover-bg-red white pointer w-10 br--right br2"
		type="submit"
		value="Delete"
		on:click|preventDefault={deleteDomain}
	/>
	{#if started}
		<input
			class="mt6 f5 h2 button-reset fr pv1 dib bl bt bb bw1 b--light-red tc w-40 br--left br2"
			type="text"
			placeholder={`please type ${name}`}
			bind:value={confirm}
		/>
	{/if}
</div>

{#if error.length > 0}
<div class="cf">
	<p class="mt3 pr3 fr light-red">{error}</p>
</div>
{/if}

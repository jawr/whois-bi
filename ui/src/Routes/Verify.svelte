<script>
	import { onMount } from 'svelte'
	import { postJSON } from '../fetchJSON'

	export let code = ''

	let status = ''
	let error = ''

	onMount(async () => {
		try {
			const response = await postJSON(`/api/verify/${code}`)
			status = response.status
		} catch (err) {
			error = err.message
		}
	})
</script>

{#if status.length > 0}
	<h1 class="f3 f2- f1-l fw2 mv3">Success</h1>
	<p>{status}</p>
{:else if error.length > 0}
	<h1 class="f3 f2- f1-l fw2 mv3 red">{error}</h1>
{:else}
	<h1 class="f3 f2- f1-l fw2 mv3 red">You Ok?</h1>
{/if}

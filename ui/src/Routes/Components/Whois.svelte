<script>
	import { onMount } from 'svelte'
	import { fetchJSON } from '../../fetchJSON'

	export let name = ''

	let whois = []
	let current = {}
	let error = ''

	const updateWhois = async () => {
		try {
			whois = await fetchJSON(`/api/user/domain/${name}/whois`)
			if (whois.length > 0) {
				current = whois[0]
			}
		} catch (err) {
			error = err.message
			error = error
		}
	}

	onMount(async () => {
		await updateWhois()
	})
</script>

{#if whois.length > 0}
	<div class="mw8">
		<pre class="pa2 tl ba bg-light-gray pre overflow-content">{window.atob(current.raw)}</pre>
	</div>
	<p class="mt5 f7">last updated: {current.added_at}</p>
{:else}
	<p class="mt5">No whois records found.</p>
{/if}

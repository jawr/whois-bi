<script>
	import { onMount } from 'svelte'
	import { fetchJSON } from '../../fetchJSON'
	import { parseISO, format } from 'date-fns'

	export let name = ''

	let whois = []
	let current = {}
	let error = ''
	let loaded = false

	const updateWhois = async () => {
		loaded = false
		try {
			whois = await fetchJSON(`/api/user/domain/${name}/whois`)
			if (whois.length > 0) {
				current = whois[0]
			}
		} catch (err) {
			error = err.message
			error = error
		}
		loaded = true
	}

	onMount(async () => {
		await updateWhois()
	})
</script>

<div class="pt2">
	<p class="tl f5 lh-copy">A Whois record shows information regarding the ownership and administration of a domain.</p>
</div>

{#if whois.length > 0}
	<div class="pa2 tl ba bg-mid-gray washed-green">
		<pre class="overflow-scroll">{window.atob(current.raw)}</pre>
	</div>
	<p class="mt5 f7">last updated: {format(parseISO(current.added_at, new Date()), 'yyyy/MM/dd HH:mm')}</p>
{:else if loaded}
	<p class="mt5">No whois records found.</p>
{/if}

<style>
	pre {
		white-space: pre-wrap;
	}
</style>

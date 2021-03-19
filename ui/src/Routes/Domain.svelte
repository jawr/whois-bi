<script>
	import { onMount } from 'svelte'
	import { fade } from 'svelte/transition'
	import { parseISO, format } from 'date-fns'
	import { fetchJSON } from '../fetchJSON'
	import { link } from 'svelte-routing'
	import Records from './Components/Records.svelte'
	import Whois from './Components/Whois.svelte'
	import Config from './Components/Config.svelte'

	export let name = ''
	export let tab = 'records'

	let domain = {}
	let error = ''

	onMount(async () => {
		try {
			domain = await fetchJSON(`/api/user/domain/${name}`)
		} catch (err) {
			error = err.message
		}
	})

	const tabClasses = 'dark-gray link w-50 tc pointer pb3 pt3 active '
	const activeClasses = 'bw1 bb b--washed-green'
</script>

<div in:fade class="mw8 center h-100">
	{#if error.length > 0}
		<h1 class="f3 f2- f1-l fw2 mv3 red">Error: {error}</h1>
	{:else}
		<h1 class="f3 f2-m f1-l fw2 mv3">{name}</h1>

		{#if domain.last_updated_at}
			<small>Last updated {format(parseISO(domain.last_updated_at, new Date()), 'yyyy/MM/dd HH:mm')}</small>
		{/if}

		<div class="mt5 flex bb b--light-gray">
			<a 
				use:link
				href={`/domain/${name}/records`} 
				class={`${tabClasses}` + (tab === 'records' ? activeClasses : '')}
			>Records</a>
			<a
				use:link
				href={`/domain/${name}/whois`} 
				class={`${tabClasses}` + (tab === 'whois' ? activeClasses : '')}
			>Whois</a>
			<a
				use:link
				href={`/domain/${name}/config`} 
				class={`${tabClasses}` + (tab === 'config' ? activeClasses : '')}
			>Config</a>
		</div>

		{#if tab === 'records'}
			<div in:fade>
				<Records {name} />
			</div>
		{:else if tab === 'whois'}
			<div in:fade>
				<Whois {name} />
			</div>
		{:else if tab === 'config'}
			<div in:fade>
				<Config {name} />
			</div>
		{/if}
	{/if}
</div>

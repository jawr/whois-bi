<script>
	import { onMount } from 'svelte'
	import { fade } from 'svelte/transition'
	import { fetchJSON } from '../fetchJSON'
	import { Link } from 'svelte-routing'
	import Records from './Components/Records.svelte'
	import Whois from './Components/Whois.svelte'
	import Jobs from './Components/Jobs.svelte'

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

	let parts = []
	let current, subtab = ''

	$: {
		parts = tab.split('/')
		current = parts[0]
		subtab = (parts.length > 1) ? parts[1] : 'current'
	}

	const tabClasses = "link w-50 tc bg-animate pointer green hover-bg-washed-green pb3 pt3 "

</script>

<div in:fade class="mw8 center h-100">
	{#if error.length > 0}
		<h1 class="f3 f2- f1-l fw2 mv3 red">Error: {error}</h1>
	{:else}
		<h1 class="f3 f2-m f1-l fw2 mv3">Details</h1>
		<p>Look in depth at '{name}'</p>

		{#if domain.last_updated_at}
			<small>Last updated {domain.last_updated_at}</small>
		{/if}

		<div class="mt5 flex bb b--dark-green">
			<Link 
				to={`domain/${name}/records`} 
				class={`${tabClasses}` + (current === 'records' ? 'bg-washed-green' : '')}
			>Records</Link>
			<Link 
				to={`domain/${name}/whois`} 
				class={`${tabClasses}` + (current === 'whois' ? 'bg-washed-green' : '')}
			>Whois</Link>
			<Link 
				to={`domain/${name}/jobs`} 
				class={`${tabClasses}` + (current === 'jobs' ? 'bg-washed-green' : '')}
			>Jobs</Link>
		</div>

		{#if current === 'records'}
			<div in:fade>
				<Records {name} tab={subtab} />
			</div>
		{:else if current === 'whois'}
			<div in:fade>
				<Whois {name} />
			</div>
		{:else if current === 'jobs'}
			<div in:fade>
			<Jobs {name} />
			</div>
		{/if}
	{/if}
</div>

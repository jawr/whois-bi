<script>
	import { onMount } from 'svelte'
	import { fade } from 'svelte/transition'
	import { Link } from 'svelte-routing'
	import { domains, domainsQuery } from '../stores'
	import { fetchJSON } from '../fetchJSON'
	import CreateDomains from './Components/CreateDomains.svelte'
	import Search from './Components/Search.svelte'

	let error = ''

	onMount(async () => {
		try {
			$domains = await fetchJSON('/api/user/domains')
		} catch (err) {
			error = '' + err
		}
	})

	let filtered = []
	$: {
		const q = $domainsQuery.toLowerCase()
		filtered = $domains.filter(d => d.domain.toLowerCase().indexOf(q) > -1)
	}
</script>

<div in:fade>
{#if $domains.length > 0}
	<Search store={domainsQuery} text="Filter domains" />

	<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
		<thead>
			<tr class="fw3 ttu f7">
				<th class="pv2 ph3 tl w-60">Name</th>
				<th class="pv2 ph3 tr w-10">#Records</th>
				<th class="pv2 ph3 tr w-10">#Whois</th>
				<th class="pv2 ph3 tr w-30">Added</th>
			</tr>
		</thead>
		<tbody>
			{#each filtered as domain}
				<tr class="striped--near-white">
					<td class="pv3 ph3 tl"><Link class="link underline green" to={`domain/${domain.domain}/records`}>{domain.domain}</Link></td>
					<td class="pv3 ph3 tr">{domain.records}</td>
					<td class="pv3 ph3 tr">{domain.whois}</td>
					<td class="pv3 ph3 tr">{domain.added_at}</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/if}

{#if error.length > 0}
	<div class="mt3">
		<p class="red">{error}</p>
	</div>
{/if}

<CreateDomains />
</div>

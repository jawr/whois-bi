<script>
	import { onMount } from 'svelte'
	import { fade } from 'svelte/transition'
	import { Link } from 'svelte-routing'
	import { domains, domainsQuery } from '../stores'
	import { fetchJSON } from '../fetchJSON'
	import { parseISO, format } from 'date-fns'

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

	const formatDateTime = (t) => {
		if (t === '0001-01-01T00:00:00Z') {
			return ''
		}
		return format(parseISO(t, new Date()), 'yyyy/MM/dd')
	}
</script>

<div class="mw8 center">
<div in:fade class="pt6">
	{#if $domains.length > 0}
		<Search store={domainsQuery} text="Filter domains" />

		<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
			<thead>
				<tr class="fw3 ttu f7">
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Name</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">#Records</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">#Whois</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white w-20">Created</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Expires</th>
				</tr>
			</thead>
			<tbody class="tl lh-copy">
				{#each filtered as domain}
					<tr>
						<td data-label="Name"     class="pv3 pr3 bb b--black-20">
							<Link class="link underline green" to={`domain/${domain.domain}/records`}>{domain.domain}</Link>
						</td>
						<td data-label="#Records" class="pv3 pr3 bb b--black-20">{domain.records}</td>
						<td data-label="#Whois"   class="pv3 pr3 bb b--black-20">{domain.whois}</td>
						<td data-label="#Created"   class="pv3 pr3 bb b--black-20">{formatDateTime(domain.created_at)}</td>
						<td data-label="#Expires"   class="pv3 pr3 bb b--black-20">{formatDateTime(domain.expires_at)}</td>
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
</div>

	<style>
	td {
		word-break: break-all;
	}

	@media screen and (max-width: 800px) {
		table thead {
			border: none;
			clip: rect(0 0 0 0);
			height: 1px;
			margin: -1px;
			overflow: hidden;
			padding: 0;
			position: absolute;
			width: 1px;
		}

		table tr {
			border-bottom: 3px solid #ddd;
			display: block;
		}

		table td {
			border-bottom: 1px solid #ddd;
			display: block;
			text-align: right;
		}

		table td::before {
			font-weight: 200;
			content: attr(data-label);
			float: left;
		}
	}
	</style>

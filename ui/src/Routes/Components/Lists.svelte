<script>
	import { onMount, onDestroy } from 'svelte'
	import { fetchJSON, postJSON } from '../../fetchJSON'
	import { parseISO, format } from 'date-fns'
	import CreateList from './CreateList.svelte'

	let whitelists = []
	let blacklists = []

	let error = ''

	onMount(async () => {
		try {
			const response = await fetchJSON(`/api/user/lists`)
			whitelists = response.whitelists
			blacklists = response.blacklists
		} catch (err) {
			error = err.message
		}
	})

	const formatDateTime = (t) => {
		if (t === '0001-01-01T00:00:00Z') {
			return ''
		}
		return format(parseISO(t, new Date()), 'yyyy/MM/dd')
	}
</script>

<div class="mw8 center">
	<div class="fl w-100 pv2">
		<p class="tl lh-copy">By default all records trigger an email, create lists to control which records can trigger emails.</p>
	</div>

{#if whitelists.length > 0}
	<h2 class="tl">Whitelists</h2>
	<div class="mt4">
		<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
			<thead>
				<tr class="fw3 ttu f7">
					<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Domain</th>
					<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">RRType</th>
					<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Record</th>
					<th scope="col" class="w-10 fw6 bb b--black-20 tl pb3 pr3 bg-white">Added</th>
				</tr>
			</thead>
			<tbody class="tl lh-copy">
				{#each whitelists as i}
					<tr>
						<td data-label="Additions" class="pv3 pr3 bb b--black-20">{i.domain}</td>
						<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.rr_type}</td>
						<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.record}</td>
						<td data-label="Created"   class="pv3 pr3 bb b--black-20">{formatDateTime(i.added_at)}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

{#if blacklists.length > 0}
	<h2 class="tl">Blacklists</h2>
	<div class="mt4">
		<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
			<thead>
				<tr class="fw3 ttu f7">
					<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Domain</th>
					<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">RRType</th>
					<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Record</th>
					<th scope="col" class="w-10 fw6 bb b--black-20 tl pb3 pr3 bg-white">Added</th>
				</tr>
			</thead>
			<tbody class="tl lh-copy">
				{#each blacklists as i}
					<tr>
						<td data-label="Additions" class="pv3 pr3 bb b--black-20">{i.domain}</td>
						<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.rr_type}</td>
						<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.record}</td>
						<td data-label="Created"   class="pv3 pr3 bb b--black-20">{formatDateTime(i.added_at)}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<div class="mt4">
	<h2 class="tl">Create</h2>
	<CreateList />
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


<script>
	import { onMount } from 'svelte'
	import { Link } from 'svelte-routing'
	import { fetchJSON } from '../../fetchJSON'
	import { records, recordsQuery } from '../../stores'

	import Search from './Search.svelte'
	import CreateRecord from './CreateRecord.svelte'

	export let name = ''
	export let tab = ''

	let filtered = []
	let error = ''

	const updateRecords = async () => {
		try {
			$records = await fetchJSON(`/api/user/domain/${name}/records`)
		} catch (err) {
			error = err.message
			error = error
		}
	}

	onMount(async () => {
		await updateRecords()
	})

	$: {
		if (tab === 'current') {
			filtered = $records.filter(r => r.removed_at.length === 0)
		} else {
			filtered = $records.filter(r => r.removed_at.length > 0)
		}

		const q = $recordsQuery.toLowerCase()

		if (q.length > 0) {
			filtered = filtered.filter(r => (
				r.name.toLowerCase().indexOf(q) > -1 
				|| r.fields.toLowerCase().indexOf(q) > -1 
				|| r.ttl.toString() === q 
				|| r.rr_type.toLowerCase() === q
			))
		}
	}

	const tabClasses = "link w-50 tc bg-animate pointer green hover-bg-near-white pb2 pt2 "
</script>

<div class="flex bb b--dark-green">
	<Link 
		to={`domain/${name}/records/current`} 
		class={`${tabClasses}` + (tab === 'current' ? 'bg-near-white' : '')}
	>Current</Link>
	<Link 
		to={`domain/${name}/records/historical`} 
		class={`${tabClasses}` + (tab === 'historical' ? 'bg-near-white' : '')}
	>Historical</Link>
</div>

<div class="pt2">
	{#if tab === 'current'}
		<p class="tl f5 lh-copy">Overview of current, active records that we know of.</p>
	{:else if tab === 'historical'}
		<p class="tl f5 lh-copy">Overview of records that are no longer active.</p>
	{/if}
</div>

<div class="mt4">

	<Search store={recordsQuery} text="Filter records" />

	<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
		<thead>
			<tr class="fw3 ttu f7">
				<th scope="col" class="pv2 ph3 tl w-20">Record</th>
				<th scope="col" class="pv2 ph3 tr w-10">Type</th>
				<th scope="col" class="pv2 ph3 tr w-40">Fields</th>
				<th scope="col" class="pv2 ph3 tr w-10">TTL</th>
				<th scope="col" class="pv2 ph3 tr w-20">Added</th>
			</tr>
		</thead>
		<tbody>
			{#each filtered as record}
				<tr class="striped--near-white">
					<td data-label="Record" class="pv3 ph3 tl">{record.name}</td>
					<td data-label="Type" class="pv3 ph3 tr">{record.rr_type}</td>
					<td data-label="Fields" class="newline pv3 ph3 tr">{record.fields}</td>
					<td data-label="TTL" class="pv3 ph3 tr">{record.ttl}</td>
					<td data-label="Added" class="pv3 ph3 tr">{record.added_at}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<div class="mw8">
	<CreateRecord name={name} />
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

		table td.newline::before {
			content: attr(data-label) '\00000a';
			text-align: left;
			display: block;
			float: none;
			margin-bottom: 1rem;
			font-weight: 200;
		}

		table td::before {
			font-weight: 200;
			content: attr(data-label);
			float: left;
		}
	}
</style>

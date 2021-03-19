<script>
	import { onMount } from 'svelte'
	import { Link } from 'svelte-routing'
	import { fetchJSON } from '../../fetchJSON'
	import { records, recordsQuery } from '../../stores'
	import { parseISO, format } from 'date-fns'

	import Search from './Search.svelte'
	import RTypeLabel from './RTypeLabel.svelte'
	import CreateRecord from './CreateRecord.svelte'

	export let name = ''

	let current = true
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

	const parseRecordName = (rname) => {
		if (rname === name + '.') {
			return '@'
		}

		return rname.replace('.' + name + '.', '')
	}

	const formatDateTime = (t) => {
		if (t === '0001-01-01T00:00:00Z') {
			return ''
		}
		return format(parseISO(t, new Date()), 'yyyy/MM/dd')
	}

	onMount(async () => {
		await updateRecords()
	})

	$: {
		if (current) {
			filtered = $records.filter(r => r.removed_at === '0001-01-01T00:00:00Z')
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
</script>

<div class="mt4 mb4">
	<Search store={recordsQuery} text="Filter records" />

	<div class="flex items-center mb2 fr">
		<input class="mr2" type="checkbox" name="current" bind:checked={current}>
		<label for="current" class="lh-copy">Current</label>
	</div>
</div>

<div class="cf">
	<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
		<thead>
			<tr class="fw3 ttu f7">
				{#if current}
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Record</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Type</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Fields</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">TTL</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Added</th>
				{:else}
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Record</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Type</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Fields</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">TTL</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Added</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Removed</th>
				{/if}
			</tr>
		</thead>
		<tbody class="tl lh-copy">
			{#each filtered as record}
				<tr>
					<td data-label="Record" class="pv3 pr3 bb b--black-20">{parseRecordName(record.name)}</td>
					<td data-label="Type"   class="pv3 pr3 bb b--black-20"><RTypeLabel type={record.rr_type} /></td>
					<td data-label="Fields" class="truncate-xl pv3 pr3 bb b--black-20 newline fields">{record.fields}</td>
					<td data-label="TTL"    class="pv3 pr3 bb b--black-20">{record.ttl}</td>
					<td data-label="Added"  class="pv3 pr3 bb b--black-20">{formatDateTime(record.added_at)}</td>
					{#if !current}
						<td data-label="Removed" class="pv3 pr3 bb b--black-20">{formatDateTime(record.removed_at)}</td>
					{/if}
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<div class="mw8 mt4">
	<CreateRecord name={name} />
</div>

<style>
	td.fields {
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

		table td.newline {
			word-break: break-all;
		}

		table td::before {
			font-weight: 200;
			content: attr(data-label);
			float: left;
		}
	}
</style>

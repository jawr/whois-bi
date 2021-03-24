<script>
	import { onMount, onDestroy } from 'svelte'
	import { Trash2Icon } from 'svelte-feather-icons'
	import { fetchJSON, postJSON } from '../fetchJSON'
	import { whitelists, blacklists } from '../stores'
	import { parseISO, format } from 'date-fns'
	import CreateList from './Components/CreateList.svelte'

	let error = ''

	onMount(async () => {
		try {
			const response = await fetchJSON(`/api/user/lists`)
			$whitelists = response.whitelists
			$blacklists = response.blacklists
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

	const deleteList = async (id) => {
		try {
			await fetchJSON(`/api/user/lists/${id}`, {
				method: 'DELETE',
			})
			whitelists.update(arr => arr.filter(i => i.id !== id))
			blacklists.update(arr => arr.filter(i => i.id !== id))
		} catch (err) {
			error = err.message
		}
	}
</script>

<div class="mw8 center">
	<div class="fl w-100 pv2">
		<p class="tl lh-copy">By default adding or removing any records triggers an email alert. Create blacklists to prevent certain records from creating an alert, i.e. SOA records.</p>
	</div>

	{#if $whitelists.length > 0}
		<h2 class="tl">Whitelists</h2>
		<div class="mt4">
			<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
				<thead>
					<tr class="fw3 ttu f7">
						<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Domain</th>
						<th scope="col" class="w-20 fw6 bb b--black-20 tl pb3 pr3 bg-white">RRType</th>
						<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Record</th>
						<th scope="col" class="w-10 fw6 bb b--black-20 tl pb3 pr3 bg-white">Added</th>
						<th scope="col" class="w-10 fw6 bb b--black-20 tl pb3 pr3 bg-white"></th>
					</tr>
				</thead>
				<tbody class="tl lh-copy">
					{#each $whitelists as i}
						<tr>
							<td data-label="Additions" class="pv3 pr3 bb b--black-20">{i.domain}</td>
							<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.rr_type}</td>
							<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.record}</td>
							<td data-label="Created"   class="pv3 pr3 bb b--black-20">{formatDateTime(i.added_at)}</td>
							<td data-label=""          class="pv3 pr3 bb b--black-20 tc">
								<span on:click|preventDefault={() => deleteList(i.id)} class="bg-red washed-red pointer pa1 white br1"><Trash2Icon size="1x" /></span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}

	{#if $blacklists.length > 0}
		<h2 class="tl">Blacklists</h2>
		<div class="mt4">
			<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
				<thead>
					<tr class="fw3 ttu f7">
						<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Domain</th>
						<th scope="col" class="w-20 fw6 bb b--black-20 tl pb3 pr3 bg-white">RRType</th>
						<th scope="col" class="w-30 fw6 bb b--black-20 tl pb3 pr3 bg-white">Record</th>
						<th scope="col" class="w-10 fw6 bb b--black-20 tl pb3 pr3 bg-white">Added</th>
						<th scope="col" class="w-10 fw6 bb b--black-20 tl pb3 pr3 bg-white"></th>
					</tr>
				</thead>
				<tbody class="tl lh-copy">
					{#each $blacklists as i}
						<tr>
							<td data-label="Additions" class="pv3 pr3 bb b--black-20">{i.domain}</td>
							<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.rr_type}</td>
							<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{i.record}</td>
							<td data-label="Created"   class="pv3 pr3 bb b--black-20">{formatDateTime(i.added_at)}</td>
							<td data-label=""          class="pv3 pr3 bb b--black-20 tc">
								<span on:click|preventDefault={() => deleteList(i.id)} class="bg-red washed-red pointer pa1 white br1"><Trash2Icon size="1x" /></span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}

	<div class="mt4">
		<h2 class="tl">Create List</h2>
		<CreateList />
	</div>
</div>

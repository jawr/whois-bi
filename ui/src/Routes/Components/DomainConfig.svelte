<script>
	import { onMount, onDestroy } from 'svelte'
	import { fetchJSON, postJSON, putJSON } from '../../fetchJSON'
	import { parseISO, format } from 'date-fns'
	import { CheckIcon } from 'svelte-feather-icons'

	import DeleteDomain from './DeleteDomain.svelte'
	import Checkbox from './Checkbox.svelte'


	export let name = ''

	let running = true
	let success = ''
	let error = ''
	let jobs = []
	let cancel = null
	let dontBatch = false
	let loaded = false
	let loading = false

	const updateJobs = async () => {
		try {
			jobs = await fetchJSON(`/api/user/jobs/${name}`)
			running = jobs.filter(j => j.finished_at === '0001-01-01T00:00:00Z').length > 0
		} catch (err) {
			error = err.message
		}

		if (running) {
			success = 'Job queued. This page will auto update when the job is complete. Please grab a cup of tea.'
			cancel = setTimeout(async () => {
				await updateJobs()
			}, 10000)
		} else {
			success = ''
		}
	}

	const requestJob = async () => {
		await postJSON(`/api/user/jobs/${name}`)
		await updateJobs()
	}

	const setDontBatch = async () => {
		loading = true
		dontBatch = !dontBatch
		await putJSON(`/api/user/domain/${name}/batch`, {dont_batch: dontBatch})
		loading = false
	}

	const getDomain = async () => {
		loaded = false
		try {
			const domain = await fetchJSON(`/api/user/domain/${name}`)
			dontBatch = domain.dont_batch
		} catch (err) {
			error = err.message
		}
		loaded = true
		console.log('loaded', loaded)
	}

	onMount(async () => {
		await updateJobs()
		await getDomain()
	})

	onDestroy(() => {
		if (cancel) {
			clearTimeout(cancel)
		}
	})

	const formatDateTime = (t) => {
		if (t === '0001-01-01T00:00:00Z') {
			return ''
		}
		return format(parseISO(t, new Date()), 'yyyy/MM/dd HH:mm')
	}
</script>

<div class="mw8 mt4">
	<DeleteDomain {name} />

	{#if loaded}
	<div class="mw8 mt4">
		<div class="flex items-center mb2 fr">
			<Checkbox on:click={setDontBatch} disabled={loading} checked={!dontBatch}>Batch alerts with other Domains</Checkbox>
		</div>
	</div>
	{/if}

<div class="pt2 cf">
	<h2 class="f4 tl fw3">Jobs</h2>
	<p class="tl f5 lh-copy">Monitor for DNS Record changes on your domain.</p>
</div>

{#if jobs.length > 0}
	<div class="mt4">
		<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
			<thead>
				<tr class="fw3 ttu f7">
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Created</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Finished</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Additions</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Removals</th>
					<th scope="col" class="fw6 bb b--black-20 tl pb3 pr3 bg-white">Whois</th>
				</tr>
			</thead>
			<tbody class="tl lh-copy">
				{#each jobs as job}
					<tr>
						<td data-label="Created"   class="pv3 pr3 bb b--black-20">{formatDateTime(job.created_at)}</td>
						<td data-label="Finished"  class="pv3 pr3 bb b--black-20">{formatDateTime(job.finished_at)}</td>
						<td data-label="Additions" class="pv3 pr3 bb b--black-20">{job.additions}</td>
						<td data-label="Removals"  class="pv3 pr3 bb b--black-20">{job.removals}</td>
						<td data-label="Whois"     class="pv3 pr3 bb b--black-20">
							{#if job.whois_updated}
								<CheckIcon size="1x" />
							{:else}
								<br />
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<div class="cf">
	{#if !running}
		<button 
			class="f5 fr pv1 h2 tc dib bb bt-0 bl-0 br-0 bw2 b--dark-green bg-animate bg-green hover-bg-green white pointer br2 grow"
			on:click|preventDefault={requestJob}>Request Update</button>
	{/if}
</div>

<div class="cf">
	{#if success.length > 0}
		<p class="mt3 fr">{success}</p>
	{/if}
</div>

<div class="cf">
	{#if error.length > 0}
		<p class="mt3 pr3 fr light-red">{error}</p>
	{/if}
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


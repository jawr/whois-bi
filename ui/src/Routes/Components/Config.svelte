<script>
	import { onMount, onDestroy } from 'svelte'
	import { fetchJSON, postJSON } from '../../fetchJSON'
	import { parseISO, format } from 'date-fns'
	import { CheckIcon } from 'svelte-feather-icons'

	import DeleteDomain from './DeleteDomain.svelte'

	export let name = ''

	let running = true
	let success = ''
	let error = ''
	let jobs = []
	let cancel = null

	const updateJobs = async () => {
		try {
			jobs = await fetchJSON(`/api/user/jobs/${name}`)
			running = jobs.filter(j => j.finished_at === '0001-01-01T00:00:00Z').length > 0
		} catch (err) {
			error = err.message
			error = error
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

	const handleClick = async () => {
		await postJSON(`/api/user/jobs/${name}`)
		await updateJobs()
	}

	onMount(async () => {
		await updateJobs()
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
</div>

<div class="pt2">
	<h2 class="f4 tl fw3">Jobs</h2>
	<p class="tl f5 lh-copy">A Job looks for records using names you have supplied or using a common list (i.e. www.{name}). Jobs are created daily, but you can request a job be run now.</p>
</div>

{#if jobs.length > 0}
	<div class="mt4">
		<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
			<thead>
				<tr class="fw3 ttu f7">
					<th role="col" class="pv2 ph3 tl w-20">Created</th>
					<th role="col" class="pv2 ph3 tr w-20">Finished</th>
					<th role="col" class="pv2 ph3 tr">Additions</th>
					<th role="col" class="pv2 ph3 tr w-10">Removals</th>
					<th role="col" class="pv2 ph3 tc w-10">Whois</th>
				</tr>
			</thead>
			<tbody>
				{#each jobs as job}
					<tr class="striped--near-white">
						<td data-label="Created" class="pv3 ph3 tl">{formatDateTime(job.created_at)}</td>
						<td data-label="Finished" class="pv3 ph3 tr">{formatDateTime(job.finished_at)}</td>
						<td data-label="Additions" class="pv3 ph3 tr">{job.additions}</td>
						<td data-label="Removals" class="pv3 ph3 tr">{job.removals}</td>
						<td data-label="Whois" class="pv3 ph3 tc">
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
			on:click|preventDefault={handleClick}>Request Update</button>
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


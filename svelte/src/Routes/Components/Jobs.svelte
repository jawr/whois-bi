<script>
	import { onMount } from 'svelte'
	import { fetchJSON, postJSON } from '../../fetchJSON'

	export let name = ''

	let running = true
	let success = ''
	let error = ''
	let jobs = []

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
			setTimeout(async () => {
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
</script>

{#if jobs.length > 0}
	<div class="mt4">
		<table class="collapse bn br2 pv2 ph3 mt4 mb4 mw8 w-100 center">
			<thead>
				<tr class="fw3 ttu f7">
					<th class="pv2 ph3 tl w-10">Created</th>
					<th class="pv2 ph3 tr w-10">Finished</th>
					<th class="pv2 ph3 tr w-10">Additions</th>
					<th class="pv2 ph3 tr w-10">Removals</th>
					<th class="pv2 ph3 tr w-10">Whois Update</th>
				</tr>
			</thead>
			<tbody>
				{#each jobs as job}
					<tr class="striped--near-white">
						<td class="pv3 ph3 tl">{job.created_at}</td>
						<td class="pv3 ph3 tr">{job.finished_at === '0001-01-01T00:00:00Z' ? '' : job.finished_at}</td>
						<td class="pv3 ph3 tr">{job.additions}</td>
						<td class="pv3 ph3 tr">{job.removals}</td>
						<td class="pv3 ph3 tr">{job.whois_updated.toString()}</td>
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

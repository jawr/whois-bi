<script>
	import { postJSON } from '../../fetchJSON'
	import { whitelists, blacklists } from '../../stores'

	let list_type = 'whitelist'
	let domain = ''
	let rr_type = ''
	let record = ''
	let error = ''
	let disabled = false

	const create = async () => {
		error = ''
		disabled = true

		const data = {
			list_type,
			domain,
			rr_type,
			record,
		}

		try {
			const response = await postJSON(`/api/user/lists`, data)
			if (response.list_type === 'blacklist') {
				blacklists.update(arr => [...arr, response])
			} else {
				whitelists.update(arr => [...arr, response])
			}
		} catch (err) {
			error = err.message
		}

		disabled = false
	}
</script>

<form on:submit|preventDefault={create} class="mw8 center">
	<div class="cf mb2">
		<div class="fl w-100 w-third-ns pv2 pr2-ns">
			<label class="db fw6 lh-copy f6 tl" for="list type">List Type</label>
			<select class="h2 pv2 input-reset ba w-100" name="list type" bind:value={list_type}>
				<option value="whitelist">Whitelist</option>
				<option value="blacklist">Blacklist</option>
			</select>
		</div>
	</div>

	<div class="cf mb4">
		<div class="fl w-100 w-third-ns pv2 pr2-ns">
			<label class="db fw6 lh-copy f6 tl" for="domain">Domain</label>
			<input placeholder="example.com" class="h2 pv2 input-reset ba w-100" name="domain" type="text" bind:value={domain} />
		</div>

		<div class="fl w-100 w-third-ns pv2 pr2-ns">
			<label class="db fw6 lh-copy f6 tl" for="record">Record</label>
			<input placeholder="www\..*" class="h2 pv2 input-reset ba w-100" name="record" type="text" bind:value={record} />
		</div>

		<div class="fl w-100 w-third-ns pv2">
			<label class="db fw6 lh-copy f6 tl" for="rr_type">RRType</label>
			<input placeholder="(SOA|NS)" class="h2 pv2 input-reset ba w-100" name="rr_type" type="text" bind:value={rr_type} />
		</div>
	</div>

	<div class="cf mb4">
		<small class="tl lh-copy mt3">Whitelists will always override blacklists. Use <code>*</code>, regexp, or exact names.</small>
	</div>

	<div class="cf">
		<div class="fl w-100 w-third-ns pv2"></div>
		<div class="fl w-100 w-third-ns pv2"></div>
		<div class="fl w-100 w-third-ns pv2">
			<input 
				name="submit"
				class="f5 pv1 h2 tc bb bt-0 bl-0 br-0 bw2 b--dark-green bg-animate bg-green hover-bg-green white pointer br2 grow w-30"
				type="submit" value="Add" {disabled} />
		</div>
	</div>
</form>

{#if error.length > 0}
	<div>
		<p class="red">{error}</p>
	</div>
{/if}

	<style>
	.anchor {
		position: absolute;
		bottom: 0;
	}
	</style>

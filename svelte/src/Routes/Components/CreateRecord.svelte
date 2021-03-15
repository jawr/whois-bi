<script>
	import { postJSON } from '../../fetchJSON'

	export let name = ''

	let message, rawRecord = ''
	let errors = []
	let disabled = false

	const handleSubmit = async () => {
		disabled = true

		try {
			const created = await postJSON('/api/user/domain', {domain})
			if (created.errors && created.errors.length > 0) {
				errors = created.errors
			} else {
				message = `Created`
			}
		} catch (err) {
			errors = [err.message]
		}

		// domains.update(arr => [...arr, ...created.filter(i => i !== undefined)])

		errors = errors
		disabled = false

		setTimeout(() => {
			message = ''
		}, 5000)

		setTimeout(() => {
			error = ''
		}, 8000)
	}
</script>

<form class="bg-washed-green mw8 center pa4 br2-ns ba b--black-10" on:submit={handleSubmit}>
	<fieldset class="cf bn ma0 pa0">

		{#if status.length > 0}
			<p class="pa0 f5 f4-ns mb3 black-80">{status}</p>
		{:else}
			<legend class="pa0 f5 f4-ns mb3 black-80">Paste a raw record or Zone file</legend>
			<div class="cf">
				<label class="clip" for="rawRecord">Raw Record</label>
				<textarea
					class="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 br2-ns br--left-ns"
					placeholder="www.whois.bi	IN	A	123.45.67.89"
					rows="8"
					name="rawRecord" 
					bind:value={rawRecord}
				></textarea>
			</div>
			<div class="cf">
				<input
					class="bb bt-0 bl-0 br-0 bw2 b--green mt3 f6 f5-l button-reset fr pv3 tc bg-animate bg-light-green grow white pointer w-100 w-25-m w-20-l br2"
					type="submit"
					value="Add"
				/>
			</div>
		{/if}
		{#each errors as error}
			<p class="pa0 f6 word-wrap black-80">{error}</p>)
		{/each}
	</fieldset>
</form>

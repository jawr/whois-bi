<script>
	import { postJSON } from '../../fetchJSON'
	import { domains } from '../../stores'

	let messages = []
	let errors = []

	let createDomains = ''
	let disabled = false

	const handleSubmit = async () => {
		disabled = true

		const split = createDomains.split(/\r?\n/)

		const create = await Promise.all(
			split.map(async (domain) => {
				try {
					const create = await postJSON('/api/user/domain', {domain})
					messages.push(`Created '${domain}'`)
					return create
				} catch (error) {
					errors.push(`Error creating '${domain}': ${error.message}`)
				}
			})
		)

		domains.update(arr => [...arr, ...create.filter(i => i !== undefined)])

		messages = messages
		errors = errors
		disabled = false

		setTimeout(() => {
			messages = []
			messages = messages
		}, 5000)

		setTimeout(() => {
			errors = []
			errors = errors
		}, 8000)
	}

</script>

<div class="pa4-l">
	<form class="bg-washed-green mw8 center pa4 br2-ns ba b--black-10" on:submit|preventDefault={handleSubmit}>
		<fieldset class="cf bn ma0 pa0">
			<legend class="pa0 f5 f4-ns mb3 black-80">Enter a Domain per line to start monitoring</legend>
			<div class="cf">
			<label class="clip" for="domain">Domain Name</label>
			<textarea
				class="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 br2-ns br--left-ns"
				placeholder="whois.bi"
				name="domain" 
				rows="0"
				bind:value={createDomains}
				{disabled}
			/>
			</div>
			<div class="cf">
			<input
				class="bb bt-0 bl-0 br-0 bw2 b--green mt3 f6 f5-l button-reset fr pv3 tc bg-animate bg-light-green grow white pointer w-100 w-25-m w-20-l br2"
				type="submit"
				value="Start"
			/>

			{#each messages as message}
			<p class="pa0 f4 mb3 black-80">{message}</p>
			{/each}

			<div class="mt3">
			{#each errors as error}
				<p class="dark-red pa0 f6 mb0 black-80">{error}</p>
			{/each}
			</div>
		</fieldset>
	</form>
	</div>

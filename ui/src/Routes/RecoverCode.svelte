<script>
	import { navigate } from 'svelte-routing'
	import { postJSON } from '../fetchJSON'

	export let code = ''

	let password = ''
	let passwordConfirm = ''
	let status = '' 
	let error = ''
	let sending = false

	const handleSubmit = async () => {
		if (password.length === 0 || password !== passwordConfirm) {
			error = 'Password does not match'
			return
		}

		const data = {code, password, password_confirm: passwordConfirm}
		try {
			const response = await postJSON('/api/recover/code', data)
			status = response.status
		} catch (err) {
			error = '' + err
		}
	}
</script>

<h1 class="f3 f2-m f1-l fw2 mv3">
	Recover
</h1>

{#if status.length > 0}
	<p>{status}</p>
{:else}
	<form on:submit|preventDefault={handleSubmit}>
		<fieldset id="sign_up" class="ba b--transparent ph0 mh0">
			<div class="mt3">
				<label class="db fw4 lh-copy f6 mb2" for="email-address">New Password</label>
				<input 
					class="pa2 input-reset ba bg-transparent w-100 measure" 
					type="password" 
					bind:value={password}
				/>
			</div>
			<div class="mt3">
				<label class="db fw4 lh-copy f6 mb2" for="email-address">Confirm New Password</label>
				<input 
					class="pa2 input-reset ba bg-transparent w-100 measure" 
					type="password" 
					bind:value={passwordConfirm}
				/>
			</div>
		</fieldset>

		{#if error.length > 0}
			<div class="mt3">
				<p class="red">{error}</p>
			</div>
		{/if}

		<div class="mt3">
			{#if sending}
				<p>Processing..</p>
			{:else}
				<button type="submit" class="f4 bb bt-0 bl-0 br-0 bw2 b--dark-green br2 pointer ph3 pv2 mt5 fw3 mb2 dib white bg-green grow">Recover!</button>
			{/if}
		</div>
	</form>
{/if}

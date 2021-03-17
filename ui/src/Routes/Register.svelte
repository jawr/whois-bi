<script>
	import { navigate } from 'svelte-routing'
	import { postJSON } from '../fetchJSON'

	let email, password, confirmPassword, error = ''
	let sending = false

	const handleSubmit = async () => {
		if (password.length === 0) {
			error = 'Please enter a password'
			return
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match'
			return
		}

		const data = {email, password}
		try {
			await postJSON('/api/register', data)
			navigate('/registered')
		} catch (err) {
			error = err
		}
	}
</script>

<h1 class="f3 f2-m f1-l fw2 mv3">
	Register
</h1>
<p>Create an account and start monitoring immediately.</p>

<form on:submit|preventDefault={handleSubmit}>
	<fieldset id="sign_up" class="ba b--transparent ph0 mh0">
		<legend class="ph0 mh0 fw6 clip">Sign Up</legend>
		<div class="mt3">
			<label class="db fw4 lh-copy f6" for="email">Email address</label>
			<input 
				class="pa2 input-reset ba bg-transparent w-100 measure" 
				type="email"
				name="email" 
				bind:value={email}
			/>
		</div>

		<div class="mt3">
			<label class="db fw4 lh-copy f6" for="password">Password</label>
			<input 
			class="b pa2 input-reset ba bg-transparent w-100 measure" 
			type="password" 
			name="password"
			bind:value={password}
			/>
		</div>

		<div class="mt3">
			<label class="db fw4 lh-copy f6" for="confirmPassword">Confirm Password</label>
			<input 
			class="b pa2 input-reset ba bg-transparent w-100 measure" 
			type="password" 
			name="confirmPassword" 
			bind:value={confirmPassword}
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
			<p>Pigeon enroute...</p>
		{:else}
			<button type="submit" class="f4 bb bt-0 bl-0 br-0 bw2 b--dark-green br2 pointer dim ph3 pv2 mt5 fw3 mb2 dib white bg-green grow">Hit it!</button>
		{/if}
	</div>
</form>

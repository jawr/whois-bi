<script>
	import { Link, link, navigate } from 'svelte-routing'
	import { loggedIn } from '../stores'
	import { postJSON } from '../fetchJSON'

	let email, password, error = ''
	let sending = false

	const handleSubmit = async () => {
		const data = {email, password}
		try {
			await postJSON('/api/login', data)
			$loggedIn = true
			navigate('/')
		} catch (err) {
			error = '' + err
		}
	}
</script>

<h1 class="f3 f2-m f1-l fw2 mv3">
	Login
</h1>
<p>If you haven't already, <Link to="register" class="no-underline green">create an account</Link> to start monitoring.</p>

<form on:submit|preventDefault={handleSubmit}>
	<fieldset id="sign_up" class="ba b--transparent ph0 mh0">
		<legend class="ph0 mh0 fw6 clip">Sign Up</legend>
		<div class="mt3">
			<label class="db fw4 lh-copy f6 mb2" for="email-address">Email address</label>
			<input 
				class="pa2 input-reset ba bg-transparent w-100 measure" 
				type="email" 
				name="email"
				bind:value={email}
			/>
		</div>
		<div class="mt3">
			<label class="db fw4 lh-copy f6 mb2" for="password">Password</label>
			<input 
				class="b pa2 input-reset ba bg-transparent w-100 measure"
				type="password"
				name="password"  
				bind:value={password}
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
			<p>Logging in...</p>
		{:else}

			<div class="lh-copy mt3">
				<a use:link href="/recover" class="f6 link dim black db">Forgot your password?</a>
			</div>

			<button type="submit" class="f4 bb bt-0 bl-0 br-0 bw2 b--dark-green br2 pointer dim ph3 pv2 mt5 fw3 mb2 dib white bg-green grow">Let me in!</button>
		{/if}
	</div>
</form>

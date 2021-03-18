<script>
	import { Link, navigate } from 'svelte-routing'
	import { loggedIn } from './stores.js'

	const logout = async () => {
		await fetch('/api/logout')
		$loggedIn = false
		navigate('/')
	}

</script>

<div class="w-100 border-box pa3 ph5-l bg-white absolute top-0 left-0 right-0">
	<div class="mw9 center db border-box">
		<nav class="db dt-l w-100">
			<span class="db dtc-l v-mid mid-gray w-100 w-25-l tc tl-l mb2 mb0-l">
				<h3 class="f3 lh-title">
					<Link to="/" title="Home" class="bg-light-green lh-copy white pa2 tracked-tight link dim">
						Whois.
					</Link>
				</h3>
			</span>
			<div class="db dtc-l v-mid w-100 w-75-l tc tr-l">
				{#if $loggedIn }
					<Link class="link dim dark-gray f6 f5-l dib mr3 mr4-l" to="/" title="Dashboard">Dashboard</Link>
					<span on:click|preventDefault={logout} class="link pointer dim dark-gray f6 f5-l dib mr3 mr4-l" href="#" title="Logout">Logout</span>
				{:else }
					<Link class="link dim dark-gray f6 f5-l dib mr3 mr4-l" to="register" title="Register">Register</Link>
					<Link class="link dim dark-gray f6 f5-l dib mr3 mr4-l" to="login" title="Login">Login</Link>
				{/if}
			</div>
		</nav>
	</div>
</div>

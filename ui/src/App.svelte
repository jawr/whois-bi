<script>
	import { onMount } from 'svelte'
	import { Router, Link, Route } from 'svelte-routing'
	import { loggedIn } from './stores'

	import Header from './Header.svelte'
	import Footer from './Footer.svelte'
	import NotFound from './NotFound.svelte'

	import Dashboard from './Routes/Dashboard.svelte'
	import Domain from './Routes/Domain.svelte'
	import Hero from './Routes/Hero.svelte'
	import Config from './Routes/Config.svelte'
	import Register from './Routes/Register.svelte'
	import Registered from './Routes/Registered.svelte'
	import Verify from './Routes/Verify.svelte'
	import Login from './Routes/Login.svelte'
	import Recover from './Routes/Recover.svelte'
	import Recovering from './Routes/Recovering.svelte'
	import RecoverCode from './Routes/RecoverCode.svelte'

	let loading = true
	onMount(async () => {
		const response = await fetch('/api/user/status')
		$loggedIn = response.ok
		loading = false
	})
</script>

{#if !loading}
<Router>
	<Header />

	<div class="mw9 center db border-box min-vh-100 h-auto">
		<section class="tc ph4 h5 pv5 w-100 dt">
			<section class="dtc v-mid tc ph3 ph4-l pv5">
				<Route path="/">
					{#if $loggedIn}
						<Dashboard />
					{:else}
						<Hero />
					{/if}
				</Route>
				<Route path="domain/:name/*tab" let:params>
					<Domain name={params.name} tab={params.tab || 'records'} />
				</Route>
				<Route path="config" component={Config} />
				<Route path="verify/:code" component={Verify} />
				<Route path="register" component={Register} />
				<Route path="registered" component={Registered} />
				<Route path="login" component={Login} />
				<Route path="recover" component={Recover} />
				<Route path="recovering" component={Recovering} />
				<Route path="recover/:code" component={RecoverCode} />
				<Route component={NotFound} />
			</section>
		</section>
	</div>

	<Footer />
</Router>
{/if}

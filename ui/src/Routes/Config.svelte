<script>
	import Lists from './Components/Lists.svelte'
	import { fade } from 'svelte/transition'
	import { link } from 'svelte-routing'

	export let tab = 'lists'

	const tabClasses = 'dark-gray link w-50 tc pointer pb3 pt3 active '
	const activeClasses = 'bw1 bb b--washed-green'
</script>

<div in:fade class="mw8 center h-100">
	<div class="mt5 flex bb b--light-gray">
		<a 
			use:link
			href={`/config/lists`} 
			class={`${tabClasses}` + (tab === 'lists' ? activeClasses : '')}
		>Lists</a>
	</div>

	{#if tab === 'lists'}
		<div in:fade>
			<Lists />
		</div>
	{/if}
</div>

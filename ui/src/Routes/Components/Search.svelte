<script>
	export let text = ''
	export let store = {}
</script>

<input 
	class="input-reset ba b--black-20 pv1 ph2 br2 mb2 w-100" 
	placeholder={text}
	type="text" 
	bind:value={$store}
/>

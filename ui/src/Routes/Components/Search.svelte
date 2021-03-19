<script>
	export let text = ''
	export let store = {}

	let expanded = false

</script>
<div class="mw8 cf center">
	<div class="mw5 fr">
		<input 
				class="input-reset ba b--black-20 pv1 ph2 br2 mb2 db w-50 fr" 
				class:w-100={expanded}
				class:tl={expanded}
				class:tc={!expanded}
				placeholder={text}
				type="text" 
				bind:value={$store}
				on:focus|preventDefault={() => expanded = true}
				on:blur|preventDefault={() => expanded = false}
			/>
	</div>
</div>

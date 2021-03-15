<script>
	export let text = ''
	export let store = {}
</script>
<div class="mw8 cf center">
	<div class="mw5 fr">
			<input 
				class="input-reset ba b--black-20 pa1 br2 mb2 db w-100" 
				placeholder={text}
				type="text" 
				bind:value={$store}
			/>
	</div>
</div>

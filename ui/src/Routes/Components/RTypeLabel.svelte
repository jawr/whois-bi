<script>
	export let type = ''

	const colours = {
		'SOA': 'gray',
		'NS': 'pink',
		'TXT': 'blue',
		'CNAME': 'yellow',
		'A': 'red',
		'MX': 'purple',
	}

	const colour = colours[type] || 'green'

	let classes = `bg-light-${colour} dark-${colour} br2 ph1 f7 dib tc`
</script>

{#if type}
	<span class={classes}>{type}</span>
{/if}

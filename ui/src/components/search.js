import React from 'react'

const Search = ({ title, onChange }) => {
  const [query, setQuery] = React.useState('')

  React.useEffect(() => {
    onChange(query)
  }, [onChange, query])

  return (
    <div className="mw8 mb5 center">
      <form className="black-80">
        <small className="f6 black-60 db mb2">{ title }</small>
        <input 
          className="input-reset ba b--black-20 pa2 mb2 db w-100" 
          type="text" 
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </form>
    </div>
  )
}

export default Search

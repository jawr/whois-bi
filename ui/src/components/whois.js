import React from 'react'
import { useWhois } from '../context/whois'

export const Whois = ({ domain }) => {
  const { whois, getWhois } = useWhois ()

  React.useEffect(() => {
    const load = async () => {
      getWhois(domain)
    }
    load()
  }, [domain, getWhois])

	if (!whois || !whois.Raw) {
		return <p>No whois records found</p>
	}

	return (
		<>
			<div className="mw8">
				<pre className="pa2 tl ba bg-light-gray pre overflow-content">{atob(whois.Raw)}</pre>
			</div>
			<p className="mt5">Whois last updated: {whois.AddedAt}</p>
		</>
	)
}

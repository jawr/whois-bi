import React from 'react'
import { useSelector } from 'react-redux'
import { selectors } from '../store/domains'

export const Whois = ({ domain }) => {
	const whois = useSelector(selectors.whoisByID(domain.ID))

	if (!whois.Raw) {
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

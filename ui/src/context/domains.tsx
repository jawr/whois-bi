import React from 'react'
import { postJSON, getJSON } from './fetchJSON'

const Context = React.createContext()

const useDomains = () => {
	const context = React.useContext(Context)
	if (!context) {
		throw new Error(`useDomains must be used within a DomainsProvider`)
	}
	return context
}

const DomainsProvider = ({ children }) => {
	const [query, setQuery] = React.useState('')
	const [domains, setDomains] = React.useState([])
	const [domain, setDomain] = React.useState({})

	const filter = (domains) => {
		if (query.length > 0) {
			const q = query.toLowerCase()
			return domains.filter(i => (
				i.Domain.toLowerCase().indexOf(q) > -1 
			))
		}
		return [...domains]
	}

	const getDomain = React.useCallback( async (name) => {
		const domain = await getJSON(`/api/user/domain/${name}`)
		setDomain(domain)
	}, [setDomain])

	const getDomains = React.useCallback( async () => {
		const domains = await getJSON('/api/user/domains')
		setDomains(domains)
	}, [setDomains])

	const createDomain = async (toCreate) => {
		const created = await postJSON('/api/user/domain', {Domain: toCreate})
		setDomains(prevDomains => [...prevDomains, created])
		return created
	}

	const deleteDomain = async (name) => {
		await fetch(`/api/user/domain/${name}`, {
			method: 'DELETE'
		})
		setDomains(prevDomains => prevDomains.filter(i => i.Domain !== name))
	}

	const values = {
		domain,
		domains,
		getDomain,
		getDomains,
		createDomain,
		deleteDomain,
		// searching
		filter,
		setQuery,
	}

	return (
		<Context.Provider value={values}>
			{children}
		</Context.Provider>
	)
}

export default DomainsProvider
export { useDomains }

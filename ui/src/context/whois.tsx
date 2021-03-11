import React from 'react'
import { getJSON } from './fetchJSON'

const Context = React.createContext()

const useWhois = () => {
	const context = React.useContext(Context)
	if (!context) {
		throw new Error(`useWhois must be used within a WhoisProvider`)
	}
	return context
}

const WhoisProvider = ({ children }) => {
	const [whois, setWhois] = React.useState({})
	const [allWhois, setAllWhois] = React.useState([])

	const getWhois = React.useCallback( async (name) => {
		const all = await getJSON(`/api/user/domain/${name}/whois`)

		if (all.length > 0) {
			setWhois(all[0])
			setAllWhois(all)
		}
	}, [setWhois, setAllWhois])

	const values = {
		whois,
		allWhois,
		getWhois,
	}

	return (
		<Context.Provider value={values}>
			{children}
		</Context.Provider>
	)
}

export default WhoisProvider
export { useWhois }

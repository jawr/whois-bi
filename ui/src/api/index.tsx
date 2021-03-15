import React from 'react'
import { AuthAPI } from './auth'
import { DomainsAPI } from './domains'
import { DomainAPI } from './domain'

type API = {
	Auth: AuthAPI
	Domains: DomainsAPI
	Domain: DomainAPI
}

const newAPI = (): API => {
	return {
		Auth: new AuthAPI(),
		Domains: new DomainsAPI(),
		Domain: new DomainAPI(),
	}
}

const Context = React.createContext(newAPI())

const useAPI = () => {
	const context = React.useContext(Context)
	if (!context) {
		throw new Error(`useAPI must be used within a APIProvider`)
	}
	return context
}

const APIProvider = (children: React.ReactNode) => {
	const context = useAPI()
	return (
		<Context.Provider value={context}>
			{children}
		</Context.Provider>
	)
}

export default APIProvider
export { useAPI }

import React from 'react'
import { postJSON } from './fetchJSON'

const Context = React.createContext()

const useLogin = () => {
	const context = React.useContext(Context)
	if (!context) {
		throw new Error(`useLogin must be used within a LoginProvider`)
	}
	return context
}


const LoginProvider = ({ children }) => {
	const [loggedIn, setLoggedIn] = React.useState(false)

	const checkLogin = React.useCallback(async () => {
		try {
			const response = await fetch('/api/user/status')
			setLoggedIn(response.ok)
		} catch (error) {
			setLoggedIn(false)
		}
	}, [setLoggedIn])

	const postLogin = async (email, password) => {
		try {
			await postJSON(
				'/api/login', 
				{
					Email: email, 
					Password: password,
				}
			)
			setLoggedIn(true)
		} catch (error) {
			setLoggedIn(false)
			throw error
		}
	}

	const postLogout = async () => {
		await fetch('/api/logout')
		setLoggedIn(false)
	}

	const register = async (email, password) => {
		await postJSON(
			'/api/register', 
			{Email: email, Password: password},
		)
	}

	const verify = async (code) => {
		await postJSON(`/api/verify/${code}`)
	}

	const value = {
		loggedIn,
		checkLogin,
		postLogin,
		postLogout,
		register,
		verify,
	}

	return (
		<Context.Provider value={value}>
			{children}
		</Context.Provider>
	)
}

export default LoginProvider
export { useLogin }

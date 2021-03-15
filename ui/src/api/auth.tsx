import { postJSON } from './fetchJSON'

export class AuthAPI {
	loggedIn: boolean = false

	// check to see if the user is already authenticated
	check = async () => {
		try {
			await fetch(`/api/user/status`)
		} catch (error) {
			this.loggedIn = false
		}

		return this.loggedIn
	}

	// attempt to register a new user
	register = async (email: string, password: string) => {
		const payload = {email, password}
		await postJSON(`/api/register`, payload)
		return null
	}

	// verify an email by checking the provided verification code
	verify = async (code: string) => {
		await fetch(`/api/verify/${code}`, {method: 'POST'})
		return null
	}

	// attempt to login with the given credentials
	login = async (email: string, password: string) => {
		const payload = {email, password}
		await postJSON(`/api/login`, payload)
		this.loggedIn = true
		return this.loggedIn
	}

	// logout, even if not successful act as though it was
	logout = async () => {
		await fetch(`/api/logout`)
		this.loggedIn = false
	}
}

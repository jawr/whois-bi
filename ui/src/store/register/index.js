import fetchWrapper from '../fetchWrapper'

export const actions = {
	register: (email, password) => (dispatch) => (
		fetchWrapper(
			'/register',
			{
				method: 'POST', 
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					Email: email,
					Password: password,
				})
			}
		)
	),
}

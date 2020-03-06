import fetchWrapper from '../fetchWrapper'
import { push } from 'connected-react-router'

export const actions = {
	login: (email, password) => (dispatch) => (
		fetchWrapper(
			'/login',
			{
				method: 'POST', 
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({
					Email: email,
					Password: password,
				})
			}
		)
		.then((result) => {
			dispatch(push('/dashboard'))
		})
	),
}

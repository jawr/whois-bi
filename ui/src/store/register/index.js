import { post } from '../fetchWrapper'

export const actions = {
	register: (email, password) => (dispatch) => (
		post(
			'/api/register',
			{
					Email: email,
					Password: password,
			}
		)
	),

	verify: (code) => (dispatch) => (
		post('/api/verify/' + code, {})
	)
}

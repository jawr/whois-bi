import { post, get } from '../fetchWrapper'
import { push } from 'connected-react-router'
import createReducer from '../createReducer'

const SET = 'login.SET'

export const actions = {
	login: (email, password) => (dispatch) => (
		post(
			'/api/login',
			{
				Email: email,
				Password: password,
			}
		)
		.then((result) => {
			dispatch(push('/dashboard'))
			dispatch({type: SET, LoggedIn: true})
		})
		.catch(error => {
			dispatch({type: SET, LoggedIn: false})
		})
	),

	check: () => (dispatch, getState) => (
		get('/api/user/status')
		.then(r => {
			dispatch({type: SET, LoggedIn: true})
			switch (getState().router.location.pathname) {
				case '/':
				case '/login':
				case '/register':
					dispatch(push('/dashboard'))
					break
				default: break
			}
		})
		.catch(error => dispatch({type: SET, LoggedIn: false}))
	),

	logout: () => (dispatch) => {
		dispatch({type: SET, LoggedIn: false})
		return get('/api/logout')
	}
}

const initialState = {
	LoggedIn: false,
}

export const reducer = createReducer(initialState, {
	[SET]: (state, action) => ({...state, LoggedIn: action.LoggedIn}),
})

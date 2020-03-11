import { get, post } from '../fetchWrapper'
import createReducer from '../createReducer'

// overwrite any existing
const GET_ALL = 'jobs.GET_ALL'

export const selectors = {
	ByDomainID: (domainID) => (state) => state.jobs.ByDomainID[domainID] || [],
}

export const actions = {
	getAll: () => (dispatch) => (
		get('/api/user/jobs')
		.then(Jobs => dispatch({type: GET_ALL, Jobs}))
	),

	create: (domain) => (dispatch) => (
		post('/api/user/jobs/' + domain.Domain)
		.then(Job => dispatch({type: GET_ALL, Jobs: [Job], unshift: true}))
	),
}

const initialState = {
	ByDomainID: {},
}

export const reducer = createReducer(initialState, {
	[GET_ALL]: (state, action) => {
		let ByDomainID = {}
		action.Jobs.forEach(i => {
			if (!(i.DomainID in ByDomainID)) {
			if (i.DomainID in state.ByDomainID) {
				ByDomainID[i.DomainID] = [...state.ByDomainID[i.DomainID]]
			} else {
				ByDomainID[i.DomainID] = []
			}
			}
			if (!ByDomainID[i.DomainID].some(j => j.ID === i.ID)) {
				if (action.unshift) {
					ByDomainID[i.DomainID].unshift(i)
				} else {
					ByDomainID[i.DomainID].push(i)
				}
			}
		})
		return {
			...state,
			ByDomainID,
		}
	},
})


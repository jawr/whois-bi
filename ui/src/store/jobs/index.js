import { get } from '../fetchWrapper'
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
}

const initialState = {
	ByDomainID: {},
}

export const reducer = createReducer(initialState, {
	[GET_ALL]: (state, action) => {
		let ByDomainID = {...state.ByDomainID}
		action.Jobs.forEach(i => {
			if (!(i.DomainID in ByDomainID)) ByDomainID[i.DomainID] = []
			if (!ByDomainID[i.DomainID].some(j => j.ID === i.ID)) ByDomainID[i.DomainID].push(i)
		})
		return {...state, ByDomainID}
	},
})


import { get, post } from '../fetchWrapper'
import createReducer from '../createReducer'
import { actions as domainActions } from '../domains'

// overwrite any existing
const GET = 'jobs.GET'
const SET_UNFINISHED = 'jobs.SET_UNFINISHED'

export const selectors = {
	ByDomainID: (domainID) => (state) => state.jobs.ByDomainID[domainID] || [],
	UnfinishedByDomainID: (domainID) => (state) => state.jobs.UnfinishedByDomainID[domainID] || [],
}

export const actions = {
	get: (domain) => (dispatch, getState) => (
		get('/api/user/jobs/' + domain.Domain)
		.then(Jobs => {
			dispatch({type: GET, Domain: domain, Jobs})

			const unfinished = Jobs.filter(j => j.FinishedAt.length === 0)
			const stateUnfinished = selectors.UnfinishedByDomainID(domain.ID)(getState())

			console.log('unfinished', unfinished, 'stateUnfinished', stateUnfinished)

			if (unfinished.length === 0 && stateUnfinished.length === 0) return Promise.resolve()

			if (unfinished.length === 0) {
				// job is complete, lets clean up
				dispatch(domainActions.get(domain.Domain))
				dispatch({type: SET_UNFINISHED, Domain: domain, Unfinished: []})

			} else if (unfinished.length > 0) {
					setTimeout(
						() => dispatch(actions.get(domain)),
						10000,
					)

				dispatch({type: SET_UNFINISHED, Domain: domain, Unfinished: unfinished})
			}
		})
	),

	create: (domain) => (dispatch) => (
		post('/api/user/jobs/' + domain.Domain)
		.then(Job => {
			dispatch(actions.get(domain))
			dispatch({type: SET_UNFINISHED, Domain: domain, Unfinished: [Job]})
		})
	),
}

const initialState = {
	UnfinishedByDomainID: {},
	ByDomainID: {},
}

export const reducer = createReducer(initialState, {
	[GET]: (state, action) => {
		return {
			...state,
			ByDomainID: {
				...state.ByDomainID,
				[action.Domain.ID]: action.Jobs,
			},
		}
	},

	[SET_UNFINISHED]: (state, action) => {
		return {
			...state,
			UnfinishedByDomainID: {
				...state.UnfinishedByDomainID,
				[action.Domain.ID]: action.Unfinished,
			},

		}
	}
})


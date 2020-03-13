import { get, post } from '../fetchWrapper'
import createReducer from '../createReducer'

// overwrite any existing
const GET = 'domains.GET'
const GET_ALL = 'domains.GET_ALL'
const CREATE = 'domains.CREATE'
const GET_RECORDS = 'domains.GET_RECORDS'
const ADD_RECORD = 'domains.ADD_RECORD'
const SEARCH_RECORDS = 'domain.SEARCH_RECORDS'
const RESET_RECORDS = 'domain.RESET_RECORDS'

export const selectors = {
	filterRecords: (records) => (state) => {
		const query = selectors.query()(state).toLowerCase()
		if (query.length > 0) {
			return records.filter(i => (
				i.Name.toLowerCase().indexOf(query) > -1 
					|| i.Fields.toLowerCase().indexOf(query) > -1 
					|| i.TTL.toString() === query 
					|| i.RRType.toLowerCase() === query
			))
		}
		return records
	},
	domainByName: (name) => (state) => state.domains.ByName[name] || {},
	recordsByID: (id) => (state) => (selectors.filterRecords((state.domains.RecordsByID[id] || []).filter(i => i.RemovedAt.length === 0))(state)),
	historicalRecordsByID: (id) => (state) => (selectors.filterRecords(state.domains.RecordsByID[id] || [])(state)),
	whoisByID: (id) => (state) => state.domains.WhoisByID[id] || {},
	query: () => (state) => state.domains.Query,
}

export const actions = {
	search: (query) => ({type: SEARCH_RECORDS, query}),
	resetSearch: (query) => ({type: RESET_RECORDS}),
	getAll: () => (dispatch) => (
		get('/api/user/domains')
		.then(Domains => dispatch({type: GET_ALL, Domains}))
	),

	get: (name) => (dispatch) => (
		get('/api/user/domain/' + name)
		.then(data => dispatch({type: GET, data}))
		.catch(error => console.log(error))
	),

	getRecords: (domain) => (dispatch) => (
		get('/api/user/domain/' + domain.Domain + '/records')
		.then(Records => {
			if (Records.length > 0) {
				dispatch({type: GET_RECORDS, DomainID: domain.ID, Records})
			}
			return Promise.resolve()
		})
	),

	create: (domain) => (dispatch) => (
		post(
			'/api/user/domain',
			{
				Domain: domain,
			},
		)
		.then(Domain => dispatch({type: CREATE, Domain}))
	),

	addRecord: (domain, rawRecord) => (dispatch) => (
		post(
			'/api/user/domain/' + domain.Domain + '/record',
			{
				Raw: rawRecord,
			},
		)
		.then(Record => {
			dispatch({type: ADD_RECORD, DomainID: domain.ID, Record})
		})
	),

}

const initialState = {
	Domains: [],
	ByName: {},
	RecordsByID: {},
	WhoisByID: {},
	Query: '',
}

const buildDomainsState = (state, Domains) => {
	if (Domains === undefined) return state
	const ByName = Domains.reduce((map, i) => {map[i.Domain] = i; return map}, {})
	return {
		...state,
		Domains,
		ByName,
	}
}

const buildRecordsAndWhoisState = (state, data) => {
	const RecordsByID = {
		...state.RecordsByID,
		[data.Domain.ID]: data.Records,
	}
	const WhoisByID = {
		[data.Domain.ID]: data.Whois,
	}
	return {
		...state,
		RecordsByID,
		WhoisByID,
	}
}

export const reducer = createReducer(initialState, {
	[GET_ALL]: (state, action) => {
		const Domains = action.Domains
		return buildDomainsState(state, Domains)
	},

	[GET]: (state, action) => {
		const Domains = state.Domains.filter(i => i.ID !== action.data.Domain.ID).concat([action.data.Domain])
		return buildRecordsAndWhoisState(buildDomainsState(state, Domains), action.data)
	},

	[GET_RECORDS]: (state, action) => {
		const RecordsByID = {
			...state.RecordsByID,
			[action.DomainID]: action.Records,
		}

		return {
			...state,
			RecordsByID,
		}
	},

	[CREATE]: (state, action) => {
		const Domains = state.Domains.concat([action.Domain])
		return buildDomainsState(state, Domains)
	},

	[ADD_RECORD]: (state, action) => {
		const RecordsByID = {
			...state.RecordsByID,
			[action.DomainID]: state.RecordsByID[action.DomainID].concat([action.Record]),
		}

		return {
			...state,
			RecordsByID,
		}
	},

	[SEARCH_RECORDS]: (state, action) => ({...state, Query: action.query}),
	[RESET_RECORDS]: (state, action) => ({...state, Query: ''}),

})


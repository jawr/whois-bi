import { combineReducers } from 'redux'
import { connectRouter } from 'connected-react-router'
import { reducer as domainsReducer } from './domains'
import { reducer as loginReducer } from './login'
import { reducer as jobsReducer } from './jobs'

const createRootReducer = (history) => combineReducers({
	router: connectRouter(history),
	domains: domainsReducer,
	login: loginReducer,
	jobs: jobsReducer,
})

export default createRootReducer

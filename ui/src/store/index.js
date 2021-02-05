import { createHashHistory as createHistory } from 'history'
import { applyMiddleware, compose, createStore } from 'redux'
import { routerMiddleware } from 'connected-react-router'
import logger from 'redux-logger'
import thunk from 'redux-thunk'
import createRootReducer from './reducers'

export const history = createHistory()

export default function configureStore(preloadedState) {
	const store = createStore(
		createRootReducer(history), 
		preloadedState,
		compose(
			applyMiddleware(
				routerMiddleware(history), 
				logger,
				thunk,
			),
		),
	)

	return store
}

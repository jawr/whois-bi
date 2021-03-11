import React from 'react'
import ReactDOM from 'react-dom'

// react context providers
import Compose from './context/compose'
import LoginProvider from './context/login'
import DomainsProvider from './context/domains'
import RecordsProvider from './context/records'
import WhoisProvider from './context/whois'
import JobsProvider from './context/jobs'

// router
import { HashRouter as Router } from 'react-router-dom'

// routes
import Root from './routes/root'

import * as serviceWorker from './serviceWorker'

const components = [
	Router,
	LoginProvider,
	DomainsProvider,
	RecordsProvider,
	WhoisProvider,
	JobsProvider,
]

ReactDOM.render(
	<Compose components={components}>
		<Root />
	</Compose>,
	document.getElementById('root'),
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()

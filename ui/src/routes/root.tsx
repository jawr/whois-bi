import React from 'react'
import { Switch, Route, useHistory } from 'react-router-dom'

// components
import Wrapper from '../components/wrapper'
import Nav from '../components/nav'
import Footer from '../components/footer'

// routes
import Hero from './hero'
import Login from './login'

import Register, { Verify, Success } from './register'
import Dashboard from './dashboard'
import Domain from './domain'

import { useLogin } from '../context/login'

const Root = () => {
	const history = useHistory()
	const { loggedIn, checkLogin, postLogout } = useLogin()

	React.useEffect(() => {
		const load = async () => {
			await checkLogin()
		}
		load()
	}, [checkLogin])

	return (
		<>
			<Nav loggedIn={loggedIn} />
			<Wrapper>
				<section className="tc ph4 h5 vh-100 pv5 w-100 dt">
					<Switch>
						<Route 
							exact 
							path="/"
							render={() => (loggedIn ? <Dashboard /> : <Hero />)}
						/>

						<Route path="/login">
							<Login />
						</Route>

						<Route 
							path="/logout"
							render={
								() => {
									postLogout()
									history.push('/')
									return <Hero />
								}
							}
						/>

						<Route exact path="/dashboard">
							<Dashboard />
						</Route>
						<Route path="/dashboard/:name">
							<Domain />
						</Route>

						<Route path="/register">
							<Register />
						</Route>
						<Route path="/success">
							<Success />
						</Route>
						<Route path="/verify/:code">
							<Verify />
						</Route>

					</Switch>
				</section>
			</Wrapper>
			<Footer />
		</>
	)
}

export default Root

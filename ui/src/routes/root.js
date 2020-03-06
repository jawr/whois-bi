import React from 'react'
import {
	Switch,
	Route,
} from 'react-router-dom'

// components
import Wrapper from '../components/wrapper'
import Nav from '../components/nav'
import Footer from '../components/footer'

// routes
import Hero from './hero'
import Register from './register'
import Login from './login'
import Success from './success'
import Dashboard from './dashboard'

export default () => (
	<>
		<Nav />

		<Wrapper>
			<section className="tc ph4 h5 vh-100 pv5 w-100 dt">
				<Switch>
					<Route exact path="/">
						<Hero />
					</Route>
					<Route path="/register">
						<Register />
					</Route>
					<Route path="/login">
						<Login />
					</Route>
					<Route path="/success">
						<Success />
					</Route>
					<Route path="/dashboard">
						<Dashboard />
					</Route>
				</Switch>
			</section>
		</Wrapper>

		<Footer />
	</>
)

import React from 'react'
import {
	Switch,
	Route,
} from 'react-router-dom'

import { useSelector } from 'react-redux'

// components
import Wrapper from '../components/wrapper'
import Nav from '../components/nav'
import Footer from '../components/footer'

// routes
import Hero from './hero'
import Register, { Verify, Success } from './register'
import Login from './login'
import Dashboard from './dashboard'
import Domain from './domain'

const Root = () => {
  const loggedIn = useSelector(state => state.login.LoggedIn)

  return (
    <>
      <Nav />

      <Wrapper>
        <section className="tc ph4 h5 vh-100 pv5 w-100 dt">
          <Switch>
            <Route 
              exact 
              path="/"
              render={() => {
                return loggedIn ? <Dashboard /> : <Hero />
              }}
            />
            <Route path="/register">
              <Register />
            </Route>
            <Route path="/login">
              <Login />
            </Route>
            <Route path="/success">
              <Success />
            </Route>
            <Route exact path="/dashboard">
              <Dashboard />
            </Route>
            <Route path="/dashboard/:name">
              <Domain />
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

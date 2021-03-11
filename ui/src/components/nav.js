import React from 'react'
import Wrapper from './wrapper'
import { Link } from 'react-router-dom'

const Nav = ({ loggedIn }) => {
  return (
    <div className="w-100 border-box pa3 ph5-l bg-white absolute top-0 left-0 right-0">
      <Wrapper>
        <nav className="db dt-l w-100">
          <Link className="db dtc-l v-mid mid-gray link dim w-100 w-25-l tc tl-l mb2 mb0-l" to="/" title="Home">
            <h3 className="f3 lh-title">
              <span className="bg-light-green lh-copy white pa2 tracked-tight">
                Whois.
              </span>
            </h3>
          </Link>
          <div className="db dtc-l v-mid w-100 w-75-l tc tr-l">
            {loggedIn ? 
              <UserLinks /> :
              <AnonymousLinks />
            }
          </div>
        </nav>
      </Wrapper>
    </div>
  )
}

export default Nav

const AnonymousLinks = () => (
  <>
    <Link className="link dim dark-gray f6 f5-l dib mr3 mr4-l" to="/register" title="Register">Register</Link>
    <Link className="link dim dark-gray f6 f5-l dib" to="/login" title="Login">Login</Link>
  </>
)

const UserLinks = () => {
  return (
    <>
      <Link className="link dim dark-gray f6 f5-l dib mr3 mr4-l" to="/dashboard" title="Dashboard">Dashboard</Link>
      <Link className="link dim dark-gray f6 f5-l dib mr3 mr4-l" to="/logout" title="Logout">Logout</Link>
    </>
  )
}

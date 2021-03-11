import React from 'react'
import { Page } from '../components/wrapper'
import { useParams } from 'react-router-dom'
import { useLogin } from '../context/login'

const Register = () => {
  const { register } = useLogin()
  const [email, setEmail] = React.useState('')
  const [password, setPassword] = React.useState('')
  const [confirmPassword, setConfirmPassword] = React.useState('')

  const [idle, setIdle] = React.useState(true)
  const [error, setError] = React.useState('')

  const handleSubmit = async (e) => {
    setIdle(false)
    setError('')

    e.preventDefault()

    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    try {
      await register(email, password)
    } catch (error) {
      setTimeout(() => {
        setError('' + error)
        setIdle(true)
      }, 1000)
    }
  }

  return (
    <Page>
      <h1 className="f3 f2-m f1-l fw2 mv3">
        Register
      </h1>
      <p>Create an account and start monitoring immediately.</p>

      <form onSubmit={handleSubmit}>
        <fieldset id="sign_up" className="ba b--transparent ph0 mh0">
          <legend className="ph0 mh0 fw6 clip">Sign Up</legend>
          <div className="mt3">
            <label className="db fw4 lh-copy f6" htmlFor="email">Email address</label>
            <input 
              className="pa2 input-reset ba bg-transparent w-100 measure" 
              type="email"
              name="email" 
              value={email}
              onChange={e => setEmail(e.target.value)}
            />
          </div>

          <div className="mt3">
            <label className="db fw4 lh-copy f6" htmlFor="password">Password</label>
            <input 
              className="b pa2 input-reset ba bg-transparent w-100 measure" 
              type="password" 
              name="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
            />
          </div>

          <div className="mt3">
            <label className="db fw4 lh-copy f6" htmlFor="confirmPassword">Confirm Password</label>
            <input 
              className="b pa2 input-reset ba bg-transparent w-100 measure" 
              type="password" 
              name="confirmPassword" 
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
            />
          </div>
        </fieldset>

        {
          error.length > 0 && 
            <div className="mt3">
              <p className="red">{error}</p>
            </div>
        }

        <div className="mt3">
          { idle && <input className="pointer input-reset f4 dim br1 ph3 pv2 bn mb2 dib white bg-green grow" type="submit" value="Hit it!" /> }
          { !idle && <p>Pigeon enroute...</p> }
        </div>
      </form>
    </Page>
  )
}

export const Success = () => (
  <Page>
    <h1 className="f3 f2-m f1-l fw2 mv3">Success</h1>
    <p>Registration complete, please check your inbox (or spambox!) to complete the registration proccess</p>
  </Page>
)

export const Verify = () => {
  const { verify } = useLogin()
  const { code } = useParams()

  const [title, setTitle] = React.useState('Verifying')
  const [text, setText] = React.useState('Please bare with us whilst we discombobulate and rejig your account')

  React.useEffect(() => {
    const load = async () => {
      try {
        await verify(code)
        setTimeout(() => {
          setTitle('Excellent!')
          setText('Your email has been verified. Please login to explore your new account!')
        }, 1000)
      } catch (error) {
        setTimeout(() => {
          setTitle('Woops!')
          setText(error)
        }, 1000)
      }
    }
    load()
  }, [code, verify])

  return (
    <Page>
      <h1 className="f3 f2-m f1-l fw2 mv3">{title}</h1>
      <p>{text}</p>
    </Page>
  )
}

export default Register

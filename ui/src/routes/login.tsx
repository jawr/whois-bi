import React from 'react'
import { Link } from 'react-router-dom'
import { Page } from '../components/wrapper'
import { useLogin } from '../context/login'
import { useHistory } from 'react-router-dom'

const Login = () => {
	const { postLogin } = useLogin()
	const history = useHistory()

	const [email, setEmail] = React.useState('')
	const [password, setPassword] = React.useState('')
	const [idle, setIdle] = React.useState(true)
	const [error, setError] = React.useState('')

	const handleSubmit = async (e) => {
		setIdle(false)

		e.preventDefault()

		setError('')

		try {
			await postLogin(email, password)
			history.push('/')
		} catch (error) {
			setError('' + error)
			setTimeout(() => {
				setIdle(true)
			}, 1000)
		}
	}

	return (
		<Page>
			<h1 className="f3 f2-m f1-l fw2 mv3">
				Login
			</h1>
			<p>If you haven't already, <Link to="register" className="no-underline green">create an account</Link> to start monitoring.</p>

			<form onSubmit={handleSubmit}>
				<fieldset id="sign_up" className="ba b--transparent ph0 mh0">
					<legend className="ph0 mh0 fw6 clip">Sign Up</legend>
					<div className="mt3">
						<label className="db fw4 lh-copy f6" htmlFor="email-address">Email address</label>
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
				</fieldset>

				{
					error.length > 0 && 
						<div className="mt3">
							<p className="red">{error}</p>
						</div>
				}

				<div className="mt3">
					{ idle && <input className="pointer input-reset f4 dim br1 ph3 pv2 bn mb2 dib white bg-green grow" type="submit" value="Let me in!" /> }
					{ !idle && <p>Logging in...</p> }
				</div>
			</form>
		</Page>
	)
}

export default Login

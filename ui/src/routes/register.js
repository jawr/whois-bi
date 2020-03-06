import React, { useState } from 'react'
import { push } from 'connected-react-router'
import { useDispatch } from 'react-redux'
import { actions } from '../store/register'

export default () => {
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [confirmPassword, setConfirmPassword] = useState('')

	const [idle, setIdle] = useState(true)
	const [error, setError] = useState('')

	const dispatch = useDispatch()

	const handleSubmit = (e) => {
		setIdle(false)
		setError('')

		e.preventDefault()


		if (password !== confirmPassword) {
			setError('Passwords do not match')
			return
		}

		dispatch(actions.register(email, password))
			.then((result) => {
				dispatch(push('/success'))
			})
			.catch((error) => {
				setTimeout(() => {
					setError('' + error)
					setIdle(true)
				}, 1000)
			})
	}

	return (
		<article className="pa4">
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
		</article>
	)
}

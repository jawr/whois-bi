import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { Page } from '../components/wrapper'
import { actions } from '../store/domains'

export default () => {
	const [loading, setLoading] = useState(true)
	const domains = useSelector(state => state.domains.Domains)

	const dispatch = useDispatch()

	useEffect(() => {
		dispatch(actions.getAll()).finally(() => setLoading(false))
	}, [dispatch]);

	return (
		<Page loading={loading}>
			<h1 className="f3 f2-m f1-l fw2 mv3">Dashboard</h1>
			<p>Currently monitoring {domains.length} domains.</p>

			{domains.length > 0 && <List domains={domains} />}

			<Create />

		</Page>
	)
}

const Create = () => {
	const [domainName, setDomainName] = useState('')
	const [error, setError] = useState('')
	const [status, setStatus] = useState('')

	const dispatch = useDispatch()

	const handleSubmit = (e) => {
		e.preventDefault()

		setStatus('Creating... Beep Boop...')

		dispatch(actions.create(domainName))
			.then(d => {
				setStatus('Added. Please wait while we prepare things for you.')
				setTimeout(() => setStatus(''), 10000)
			})
			.catch(error => {
				setError(error)
				setStatus('')
			})
	}

	return (
		<div className="pa4-l">
			<form className="bg-washed-green mw7 center pa4 br2-ns ba b--black-10" onSubmit={handleSubmit}>
				<fieldset className="cf bn ma0 pa0">

					{status.length > 0 && <p className="pa0 f5 f4-ns mb3 black-80">{status}</p>}
					{status.length === 0 &&
						<>
							<legend className="pa0 f5 f4-ns mb3 black-80">Enter a Domain to start monitoring</legend>
							<div className="cf">
								<label className="clip" htmlFor="domain">Domain Name</label>
								<input
									className="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 w-75-m w-80-l br2-ns br--left-ns"
									placeholder="whois.bi"
									type="text"
									name="domain" 
									value={domainName}
									onChange={e => setDomainName(e.target.value)}
								/>
								<input
									className="f6 f5-l button-reset fl pv3 tc bn bg-animate bg-black-70 hover-bg-black white pointer w-100 w-25-m w-20-l br2-ns br--right-ns"
									type="submit"
									value="Start"
								/>
							</div>
						</>
					}
					{error.length > 0 && <p className="pa0 f5 f4-ns mb3 black-80">{error}</p>}
				</fieldset>
			</form>
		</div>
	)
}

const List = ({ domains }) => (
	<ul className="list pl0 mt6 mb6 center mw7">
		{domains.map(d => <Item key={d.ID} domain={d} />)}
	</ul>
)

const Item = ({ domain }) => (
	<li
		className="flex items-center lh-copy pa3 ph0-l bb b--black-10 mw7"
	>
		<div className="pl3 flex-auto">
			<span className="f6 db black-70">{domain.Domain}</span>
		</div>
		<div>
			<Link to={`/dashboard/${domain.Domain}`} className="f6 link blue hover-dark-gray">details</Link>
		</div>
	</li>
)

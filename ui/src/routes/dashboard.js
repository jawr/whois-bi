import React, { useState, useEffect, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { Page } from '../components/wrapper'
import { actions, selectors } from '../store/domains'
import { Table } from '../components/table'

export default () => {
	const [loading, setLoading] = useState(true)
	const domains = useSelector(selectors.filterDomains())

	const dispatch = useDispatch()

	useEffect(() => {
		dispatch(actions.getAll()).finally(() => setLoading(false))
	}, [dispatch]);

	const columns = useMemo(() => [
		{
			Header: 'Domain', accessor: 'Domain', 
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-60',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Records', accessor: 'Records', 
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Whois', accessor: 'Whois', 
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Updated', accessor: 'LastUpdatedAt',
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
			cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
		},
		{
			Header: '', id: 'Options',
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
			cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
			dontSort: true,
			Cell: ({ cell }) => <Link to={`/dashboard/${cell.row.original.Domain}`} className="f6 link blue hover-dark-gray">details</Link>
		}
	], [])

	return (
		<Page loading={loading}>
			<h1 className="f3 f2-m f1-l fw2 mv3">Dashboard</h1>
			<p>Currently monitoring {domains.length} domains.</p>

			<div className="mv5">
				<Search />
				<Table data={domains} columns={columns} />
			</div>

			<Create />

		</Page>
	)
}

const Create = () => {
	const [domainName, setDomainName] = useState('')
	const [error, setError] = useState('')
	const [status, setStatus] = useState('')
  const [disabled, setDisabled] = useState(false)

	const dispatch = useDispatch()

	const handleSubmit = (e) => {
		e.preventDefault()

		setStatus('Creating... Beep Boop...')
    setDisabled(true)

		dispatch(actions.create(domainName))
			.then(d => {
				setStatus('Added. Please wait while we prepare things for you.')
        setDomainName('')
				setTimeout(() => setStatus(''), 5000)
			})
			.catch(error => {
				setError(error)
				setStatus('')
			})
      .finally(() => {
        setDisabled(false)
      })
	}

	return (
		<div className="pa4-l">
			<form className="bg-washed-green mw8 center pa4 br2-ns ba b--black-10" onSubmit={handleSubmit}>
				<fieldset className="cf bn ma0 pa0">
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
              disabled={disabled}
            />
            <input
              className="f6 f5-l button-reset fl pv3 tc bn bg-animate bg-black-70 hover-bg-black white pointer w-100 w-25-m w-20-l br2-ns br--right-ns"
              type="submit"
              value="Start"
            />
          </div>
					{status.length > 0 && <p className="pa0 f5 f4-ns mb3 black-80">{status}</p>}
					{error.length > 0 && <p className="pa0 f5 f4-ns mb3 black-80">{error}</p>}
				</fieldset>
			</form>
		</div>
	)
}


const Search = () => {
	const dispatch = useDispatch()
	const query = useSelector(selectors.query())
	return (
		<div className="mw8 mb5 center">
			<form className="black-80">
				<small className="f6 black-60 db mb2">Filter records results</small>
				<input 
					className="input-reset ba b--black-20 pa2 mb2 db w-100" 
					type="text" 
					value={query}
					onChange={(e) => dispatch(actions.search(e.target.value))}
				/>
			</form>
		</div>
	)
}


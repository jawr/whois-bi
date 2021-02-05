import React, { useState, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Loader } from './wrapper'
import { actions, selectors } from '../store/jobs'

export const Jobs = ({ domain }) => {
	const jobs = useSelector(selectors.ByDomainID(domain.ID))
	const [loading, setLoading] = useState(false)

	const dispatch = useDispatch()

	useEffect(() => {
		if (domain.Domain !== undefined) {
			setLoading(true)
			dispatch(actions.get(domain)).finally(() => setLoading(false))
		}
	}, [dispatch, domain]);

	return (
		<div className="min-vh-100">
			{loading && <Loader />}
			{!loading && 
				<>
					<Table jobs={jobs} />

					<Add domain={domain} />
				</>
			}
		</div>
	)
}

const Add = ({ domain }) => {
	const [idle, setIdle] = useState(true)
	const [success, setSuccess] = useState('')
	const [error, setError] = useState('')

	const unfinishedJobs = useSelector(selectors.UnfinishedByDomainID(domain.ID))

	const dispatch = useDispatch()

	useEffect(() => {
		if (unfinishedJobs.length > 0) {
					setSuccess('Job queued. This page will auto update when the job is complete. Please grab a cup of tea.')
		} else {
					setSuccess('')
		}
	}, [unfinishedJobs])

	const handleClick = (e) => {
		e.preventDefault()

			setIdle(false)
			setError('')
			dispatch(actions.create(domain))
				.catch(error => {
					setTimeout(() => setError(error), 1000)
				})
				.finally(() => {
					setTimeout(() => setIdle(true), 1000)
				})
	}

	return (
		<div className="tr mt5">
			{ idle && success.length === 0 && <button className="pointer input-reset f4 dim br1 ph3 pv2 bn mb2 dib green bg-washed-green grow" onClick={handleClick}>Request Update</button> }
			{ !idle && <p></p> }
			{ success.length > 0 && <p>{success}</p> }
			{ idle && error.length > 0 && <p>{error}</p> }
		</div>
	)
}

const Table = ({ jobs }) => {
	if (jobs.length === 0) return <p>No jobs found</p>

		return (
				<table className="f6 mw8 center dt--fixed" cellSpacing="0">
					<thead>
						<tr>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-20">Domain</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Created</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Finished</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Additions</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Removals</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Whois Update</th>
						</tr>
					</thead>
					<tbody className="lh-copy tl">
						{jobs.map(job => <Row key={job.ID} job={job} />)}
					</tbody>
				</table>
		)
}

const Row = ({ job, historical }) => (
		<tr>
			<td className="pv3 ph3 bb b--black-20">{job.Domain.Domain}</td>
			<td className="pv3 ph3 bb b--black-20">{job.CreatedAt}</td>
			<td className="pv3 ph3 bb b--black-20">{job.FinishedAt}</td>
			<td className="pv3 ph3 bb b--black-20">{job.Additions}</td>
			<td className="pv3 ph3 bb b--black-20">{job.Removals}</td>
			<td className="pv3 ph3 bb b--black-20">{job.WhoisUpdated.toString()}</td>
		</tr>
	)

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
			dispatch(actions.getAll(domain)).finally(() => setLoading(false))
		}
	}, [dispatch, domain]);

	return (
		<div className="min-vh-100">
			{loading && <Loader />}
			{!loading && <Table jobs={jobs} />}
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

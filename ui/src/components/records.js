import React, { useState, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { actions, selectors } from '../store/domains'
import { Tabs, Menu, Tab, Panels, Panel } from '../components/tabs'

export const Records = ({ domain }) => {
	const records = useSelector(selectors.recordsByID(domain.ID))
	const historical = useSelector(selectors.historicalRecordsByID(domain.ID))
	const [loadingHistorical, setLoadingHistorical] = useState(false)

	const dispatch = useDispatch()

	useEffect(() => {
		if (domain.Domain !== undefined) {
			setLoadingHistorical(true)
			dispatch(actions.getRecords(domain)).finally(() => setLoadingHistorical(false))
		}
	}, [dispatch, domain]);

	return (
		<Tabs sub={true}>
			<Menu small={true}>
				<Tab sub={true}>Current</Tab>
				<Tab sub={true}>Historical</Tab>
			</Menu>

			<Panels small={true}>
				<Panel>
					<Table lastUpdatedAt={domain.LastUpdatedAt} records={records} />
					<Add domain={domain} />
				</Panel>
				<Panel loading={loadingHistorical}>
					<Table lastUpdatedAt={domain.LastUpdatedAt} records={historical} historical />
				</Panel>
			</Panels>
		</Tabs>
	)
}

const Table = ({ lastUpdatedAt, records, historical }) => {
	if (records.length === 0) return <p>No records found</p>

		return (
			<>
				<table className="f6 mw8 center dt--fixed" cellSpacing="0">
					<thead>
						<tr>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-20">Name</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">Type</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-40 dn dtc-ns">Fields</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10">TTL</th>
							<th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10 dn dtc-ns">Added</th>
							{historical && <th className="fw6 bb b--black-20 tl pb3 pr3 bg-white w-10 dn dtc-ns">Removed</th>}
						</tr>
					</thead>
					<tbody className="lh-copy tl">
						{records.map(record => <Row key={record.ID} record={record} historical={historical} />)}
					</tbody>
				</table>
				<p className="mt5">Records last updated: {lastUpdatedAt}</p>
			</>
		)
}

const Row = ({ record, historical }) => {
	const classes = (historical && record.RemovedAt.length > 0) ? "bg-washed-red" : ""
	return (
		<tr className={classes}>
			<td className="pv3 ph3 bb b--black-20">{record.Name}</td>
			<td className="pv3 ph3 bb b--black-20">{record.RRType}</td>
			<td className="pv3 ph3 bb b--black-20 dn dtc-ns">
				<div className="overflow-content">{record.Fields}</div>
			</td>
			<td className="pv3 ph3 bb b--black-20">{record.TTL}</td>
			<td className="pv3 ph3 bb b--black-20 dn dtc-ns">{record.AddedAt}</td>
			{historical && <td className="pv3 ph3 bb b--black-20 dn dtc-ns">{record.RemovedAt}</td>}
		</tr>
	)
}

const Add = ({ domain }) => {
	const [rawRecord, setRawRecord] = useState('')
	const [error, setError] = useState('')
	const [status, setStatus] = useState('')

	const dispatch = useDispatch()

	const handleSubmit = (e) => {
		e.preventDefault()

		setError('')
		setStatus('Adding... Beep Boop...')

		dispatch(actions.addRecord(domain, rawRecord))
			.then(record => {
				setStatus('Added.')
				setTimeout(() => setStatus(''), 1000)
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
							<legend className="pa0 f5 f4-ns mb3 black-80">Enter a raw Record to add to this domain</legend>
							<div className="cf">
								<label className="clip" htmlFor="rawRecord">Raw Rcord</label>
								<input
									className="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 w-75-m w-80-l br2-ns br--left-ns"
									placeholder="www.whois.bi	IN	A	123.45.67.89"
									type="text"
									name="rawRecord" 
									value={rawRecord}
									onChange={e => setRawRecord(e.target.value)}
								/>
								<input
									className="f6 f5-l button-reset fl pv3 tc bn bg-animate bg-black-70 hover-bg-black white pointer w-100 w-25-m w-20-l br2-ns br--right-ns"
									type="submit"
									value="Add"
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



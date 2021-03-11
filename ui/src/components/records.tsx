import React, { useMemo, useState, useEffect } from 'react'
import { Table } from './table'
import Search from './search'
import { Tabs, Menu, Tab, Panels, Panel } from '../components/tabs'
import { useRecords } from '../context/records'

export const Records = ({ domain }) => {

	const {
		records,
		historical,
		getRecords,
		filter,
		setQuery,
	} = useRecords()

	const [loading, setLoading] = useState(true)

	const columns = useMemo(() => RecordTableColumns, [])
	const historicalColumns = useMemo(() => HistoricalRecordTableColumns, [])

	useEffect(() => {
		const load = async () => {
			setLoading(true)
			getRecords(domain)
			setLoading(false)
		}
		load()
	}, [domain, getRecords]);

	return (
		<Tabs sub={true}>
			<Menu small={true}>
				<Tab sub={true}>Current</Tab>
				<Tab sub={true}>Historical</Tab>
			</Menu>

			<Panels small={true}>
				<Panel loading={loading}>
					<RecordsTable setQuery={setQuery} columns={columns} data={filter(records)}/>
					<Add domain={domain} />
				</Panel>
				<Panel loading={loading}>
					<RecordsTable setQuery={setQuery} columns={historicalColumns} data={filter(historical)} />
				</Panel>
			</Panels>
		</Tabs>
	)
}

const RecordTableColumns = [
	{
		Header: 'Name', accessor: 'Name', 
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-30',
		cellClassName: 'pv2 ph3 bb b--black-20',
	},
	{
		Header: 'Type', accessor: 'RRType', 
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-5',
		cellClassName: 'pv2 ph3 bb b--black-20',
	},
	{
		Header: 'Fields', accessor: 'Fields', 
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-40 dn dtc-ns',
		cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns mw1',
		Cell: ({ cell: { value } }) => <div className="word-wrap">{value}</div>,
	},
	{
		Header: 'TTL', accessor: 'TTL',
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-5',
		cellClassName: 'pv2 ph3 bb b--black-20',
	},
	{
		Header: 'Added', accessor: 'AddedAt',
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-15 dn dtc-ns',
		cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
	},
]

const HistoricalRecordTableColumns = [
	{
		Header: 'Name', accessor: 'Name', 
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-30',
		cellClassName: 'pv2 ph3 bb b--black-20',
	},
	{
		Header: 'Type', accessor: 'RRType', 
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-5',
		cellClassName: 'pv2 ph3 bb b--black-20',
	},
	{
		Header: 'Fields', accessor: 'Fields', 
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-40 dn dtc-ns',
		cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns mw1',
		Cell: ({ cell: { value } }) => <div className="word-wrap">{value}</div>,
	},
	{
		Header: 'TTL', accessor: 'TTL',
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-5',
		cellClassName: 'pv2 ph3 bb b--black-20',
	},
	{
		Header: 'Added', accessor: 'AddedAt',
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10 dn dtc-ns',
		cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
	},
	{
		Header: 'Removed', accessor: 'RemovedAt',
		headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white dn dtc-ns w-10',
		cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
	},
]

const RecordsTable = ({ setQuery, ...props }) => {
	useEffect(() => {
		// dispatch(actions.resetSearchRecords())
	}, []);

	return (
		<div className="mv5">
			<Search title="Filter records" onChange={setQuery} />
			<Table {...props} />
		</div>
	)
}

const Add = ({ domain }) => {
	const { createRecord } = useRecords()

	const [raw, setRaw] = useState('')
	const [errors, setErrors] = useState([])
	const [status, setStatus] = useState('')

	const handleSubmit = async (e) => {
		e.preventDefault()

		setErrors([])
		setStatus('Adding... Beep Boop...')

		const data = await createRecord(domain, raw)
			.catch(error => {
				setErrors([''+error])
				setStatus([])
			})

		if (data.Records.length > 0) {
			setStatus("Added " + data.Records.length + " records")
			setTimeout(() => setStatus(''), 1000)
		}

		if (data.Errors.length > 0) {
			setStatus([])
			setErrors(data.Errors)
		}
	}

	return (
		<form className="bg-washed-green mw9 center pa4 br2-ns ba b--black-10" onSubmit={handleSubmit}>
			<fieldset className="cf bn ma0 pa0">

				{status.length > 0 && <p className="pa0 f5 f4-ns mb3 black-80">{status}</p>}
				{status.length === 0 &&
					<>
						<legend className="pa0 f5 f4-ns mb3 black-80">Paste a raw record or Zone file</legend>
						<div className="cf">
							<label className="clip" htmlFor="raw">Raw Record</label>
							<textarea
								className="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 br2-ns br--left-ns"
								placeholder="www.whois.bi	IN	A	123.45.67.89"
								rows="8"
								name="raw" 
								value={raw}
								onChange={e => setRaw(e.target.value)}
							></textarea>
						</div>
						<div className="cf">
							<input
								className="mt3 f6 f5-l button-reset fr pv3 tc bn bg-animate bg-light-green hover-bg-green white pointer w-100 w-25-m w-20-l br2-ns br--right-ns"
								type="submit"
								value="Add"
							/>
						</div>
					</>
				}
				{errors.length > 0 && errors.map(e => <p className="pa0 f6 word-wrap black-80">{e}</p>)}
			</fieldset>
		</form>
	)
}

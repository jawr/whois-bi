import React, { useMemo, useState, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Table } from './table'
import { actions, selectors } from '../store/domains'
import { Tabs, Menu, Tab, Panels, Panel } from '../components/tabs'

export const Records = ({ domain }) => {
	const records = useSelector(selectors.recordsByID(domain.ID))
	const historical = useSelector(selectors.historicalRecordsByID(domain.ID))

	const [loadingHistorical, setLoadingHistorical] = useState(false)

	const dispatch = useDispatch()

	const columns = useMemo(() => [
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
			// cell wrapper needed
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
	], [])

	const historicalColumns = useMemo(() => [
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
			// cell wrapper needed
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
	], [])

	useEffect(() => {
		if (domain.Domain !== undefined) {
			setLoadingHistorical(true)
			dispatch(actions.getRecords(domain)).finally(() => setLoadingHistorical(false))
			dispatch(actions.resetSearchRecords())
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

					<RecordsTable columns={columns} data={records} />
					<Add domain={domain} />
				</Panel>
				<Panel loading={loadingHistorical}>
					<RecordsTable columns={historicalColumns} data={historical} />
				</Panel>
			</Panels>
		</Tabs>
	)
}

const RecordsTable = (props) => {
	const dispatch = useDispatch()

	useEffect(() => {
		dispatch(actions.resetSearchRecords())
	}, [dispatch]);

	return (
		<div className="mv5">
			<Search />
			<Table {...props} />
		</div>
	)
}

const Add = ({ domain }) => {
	const [rawRecord, setRawRecord] = useState('')
	const [errors, setErrors] = useState([])
	const [status, setStatus] = useState('')

	const dispatch = useDispatch()

	const handleSubmit = (e) => {
		e.preventDefault()

		setErrors([])
		setStatus('Adding... Beep Boop...')

		dispatch(actions.addRecord(domain, rawRecord))
			.then(data => {
				if (data.Records.length > 0) {
					setStatus("Added " + data.Records.length + " records")
					setTimeout(() => setStatus(''), 1000)
				}

				if (data.Errors.length > 0) {
					setStatus('')
					setErrors(data.Errors)
				}
			})
			.catch(err => {
				setErrors([''+err])
				setStatus('')
			})
	}

	return (
			<form className="bg-washed-green mw9 center pa4 br2-ns ba b--black-10" onSubmit={handleSubmit}>
				<fieldset className="cf bn ma0 pa0">

					{status.length > 0 && <p className="pa0 f5 f4-ns mb3 black-80">{status}</p>}
					{status.length === 0 &&
						<>
							<legend className="pa0 f5 f4-ns mb3 black-80">Paste a raw record or Zone file</legend>
							<div className="cf">
								<label className="clip" htmlFor="rawRecord">Raw Record</label>
								<textarea
									className="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 br2-ns br--left-ns"
									placeholder="www.whois.bi	IN	A	123.45.67.89"
									rows="8"
									name="rawRecord" 
									value={rawRecord}
									onChange={e => setRawRecord(e.target.value)}
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

const Search = () => {
	const dispatch = useDispatch()
	const query = useSelector(selectors.recordsQuery())
	return (
		<div className="w-100 mb5 center">
			<form className="black-80">
				<small className="f6 black-60 db mb2">Filter records results</small>
				<input 
					className="input-reset ba b--black-20 pa2 mb2 db w-100" 
					type="text" 
					value={query}
					onChange={(e) => dispatch(actions.searchRecords(e.target.value))}
				/>
			</form>
		</div>
	)
}

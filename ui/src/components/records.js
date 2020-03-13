import React, { useMemo, useState, useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useTable, useSortBy } from 'react-table'

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
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-20',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Type', accessor: 'RRType', 
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
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
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
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
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-15',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Type', accessor: 'RRType', 
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Fields', accessor: 'Fields', 
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-35 dn dtc-ns',
			cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns mw1',
			// cell wrapper needed
			Cell: ({ cell: { value } }) => <div className="word-wrap">{value}</div>,
		},
		{
			Header: 'TTL', accessor: 'TTL',
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-10',
			cellClassName: 'pv2 ph3 bb b--black-20',
		},
		{
			Header: 'Added', accessor: 'AddedAt',
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-15 dn dtc-ns',
			cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
		},
		{
			Header: 'Removed', accessor: 'RemovedAt',
			headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white dn dtc-ns w-15',
			cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
		},
	], [])

	useEffect(() => {
		if (domain.Domain !== undefined) {
			setLoadingHistorical(true)
			dispatch(actions.getRecords(domain)).finally(() => setLoadingHistorical(false))
			dispatch(actions.resetSearch())
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

					<Table columns={columns} data={records} />
					<Add domain={domain} />
				</Panel>
				<Panel loading={loadingHistorical}>
					<Table columns={historicalColumns} data={historical} />
				</Panel>
			</Panels>
		</Tabs>
	)
}

// expose more getters
// https://codesandbox.io/s/github/tannerlinsley/react-table/tree/master/examples/data-driven-classes-and-styles
const Table = ({ columns, data }) => {
	const dispatch = useDispatch()

	useEffect(() => {
			dispatch(actions.resetSearch())
	}, [dispatch]);

	const {
		getTableProps,
		getTableBodyProps,
		headerGroups,
		rows,
		prepareRow,
	} = useTable(
		{
			columns,
			data,
		},
		useSortBy
	)

	if (data.length === 0) {
		return <p>No records found</p>
	}

	return (
		<>
			<Search />

			<table {...getTableProps()} className="f6 mw8 center dt--fixed" cellSpacing="0">
				<thead>
					{headerGroups.map(headerGroup => (
						<tr {...headerGroup.getHeaderGroupProps()}>
							{headerGroup.headers.map(column => (
								<th {...column.getHeaderProps({
									...column.getSortByToggleProps(),
									className: column.headerClassName,
								})}
								>
										{column.isSorted
											? column.isSortedDesc
												? <span className="sort-by desc"></span>
												: <span className="sort-by asc"></span>
												: <span className="sort-by"></span>}
									{column.render('Header')}
								</th>
							))}
						</tr>
					))}
				</thead>
				<tbody {...getTableBodyProps()} className="lh-copy tl">
					{rows.map(
						(row, i) => {
							prepareRow(row);
							return (
								<tr {...row.getRowProps()}>
									{row.cells.map(cell => {
										return (
											<td {...cell.getCellProps({
												className: cell.column.cellClassName,
											})}>{cell.render('Cell')}</td>
										)
									})}
								</tr>
							)}
					)}
				</tbody>
			</table>
		</>
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

const Search = () => {
	const dispatch = useDispatch()
	const query = useSelector(selectors.query())
	return (
		<div className="w-100 mb5 center">
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

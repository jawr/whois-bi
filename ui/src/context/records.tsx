import React from 'react'
import { postJSON, getJSON } from './fetchJSON'

const Context = React.createContext()

const useRecords = () => {
	const context = React.useContext(Context)
	if (!context) {
		throw new Error(`useRecords must be used within a RecordsProvider`)
	}
	return context
}

const RecordsProvider = ({ children }) => {
	const [query, setQuery] = React.useState('')
	const [records, setRecords] = React.useState([])
	const [historical, setHistorical] = React.useState([])

	const filter = (records) => {
		if (query.length > 0) {
			const q = query.toLowerCase()
			return records.filter(i => (
				i.Name.toLowerCase().indexOf(q) > -1 
					|| i.Fields.toLowerCase().indexOf(q) > -1 
					|| i.TTL.toString() === q 
					|| i.RRType.toLowerCase() === q
			))
		}
		return [...records]
	}

	const createRecord = async (domain, raw) => {
		return postJSON(
			`/api/user/domain/${domain}/record`,
			{Raw: raw},
		)
	}

	const getRecords = React.useCallback( async (name) => {
		const records = await getJSON(`/api/user/domain/${name}/records`)

		setRecords(records.filter(i => i.RemovedAt.length === 0))
		setHistorical(records.filter(i => i.RemovedAt.length > 0))
	}, [setRecords, setHistorical])

	const values = {
		records,
		historical,
		getRecords,
		createRecord,
		filter,
		setQuery,
	}

	return (
		<Context.Provider value={values}>
			{children}
		</Context.Provider>
	)
}

export default RecordsProvider
export { useRecords }

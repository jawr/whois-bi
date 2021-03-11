import React from 'react'
import { postJSON, getJSON } from './fetchJSON'

const Context = React.createContext()

const useJobs = () => {
	const context = React.useContext(Context)
	if (!context) {
		throw new Error(`useJob must be used within a JobProvider`)
	}
	return context
}

const JobProvider = ({ children }) => {
	const [jobs, setJobs] = React.useState([])
	const [unfinishedJobs, setUnfinishedJobs] = React.useState([])

	const getJobs = React.useCallback( async (name) => {
		const jobs = await getJSON(`/api/user/jobs/${name}`)
		setJobs(jobs)

		const unfinished = jobs.filter(j => j.FinishedAt.length === 0)
		setUnfinishedJobs(unfinished)

	}, [setJobs, setUnfinishedJobs])

	const createJob = async (name) => {
		const job = await postJSON(`/api/user/jobs/${name}`)
		setUnfinishedJobs(prevJobs => [...prevJobs, job])
	}

	const values = {
		jobs,
		unfinishedJobs,
		getJobs,
		createJob,
	}

	return (
		<Context.Provider value={values}>
			{children}
		</Context.Provider>
	)
}

export default JobProvider
export { useJobs }

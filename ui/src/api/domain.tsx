import { getJSON, postJSON } from './fetchJSON'
import { 
	APIError,
	Domain,
	Record,
	Whois,
	Job,
} from './types'

export class DomainAPI {
	domain: Domain = {} as Domain
	records: Record[] = [] as Record[]
	whois: Whois[] = [] as Whois[]
	jobs: Job[] = [] as Job[]

	getDomain = async (name: string) => {
		this.domain = await getJSON(`/api/user/domain/${name}`)
		return this.domain
	}

	getRecords = async (name: string) => {
		this.records = await getJSON(`/api/user/domain/${name}/records`)
		return this.records
	}

	getWhois = async (name: string) => {
		this.whois = await getJSON(`/api/user/domain/${name}/whois`)
		return this.whois
	}

	getJobs = async (name: string) => {
		this.jobs = await getJSON(`/api/user/domain/${name}/jobs`)
		return this.jobs
	}
}

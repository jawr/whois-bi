export type APIError = {
	error: string
}

export type Domain = {
	id: number
	domain: string
	owner_id: number
	added_at: string
	deleted_at: string
	last_job_at: string
	last_updated_at: string
}

export type Record = {
	id: number
	domain_id: number
	record_source: number
	raw: string
	fields: string
	name: string
	rr_type: string
	rr_class: number
	ttl: number
	hash: number
	added_at: string
	removed_at: string
	deleteda_at: string
}

export type Whois = {
	id: number
	domain_id: number
	raw: string
	version: string
	created_date: string
	updated_date: string
	expiration_date: string
	date_errors: string[]
	added_at: string
	deleted_at: string
}

export type Job = {
	id: number
	domain_id: number
	errors: string[]
	additions: number
	removals: number
	whois_updated: boolean
	created_at: string
	started_at: string
	finished_at: string
}

import { writable } from 'svelte/store'
export const loggedIn = writable(false)
export const domain = writable({})
export const domains = writable([])
export const domainsQuery = writable('')
export const records = writable([])
export const recordsQuery = writable('')

async function http<T>(request: RequestInfo, options: RequestInit): Promise<T> {
	const response: Response = await fetch(request, options)
	const parsed: T = await response.json()

	if (!response.ok) {
		return Promise.reject(new Error(response.statusText))
	}

	return parsed
}

async function getJSON<T>(url: RequestInfo): Promise<T> {
	return http<T>(url, {})
}

async function postJSON<T>(url: RequestInfo, data: T): Promise<T> {
	return http<T>(url, {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify(data),
	})
}

export {
	http,
	getJSON,
	postJSON,
}

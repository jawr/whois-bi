const fetchJSON = async (url, options) => {
	const response = await fetch(url, options)
	const json = await response.json()

	if (!response.ok) {
		const contentType = response.headers.get("content-type")
		if (contentType && contentType.indexOf("application/json") !== -1) {
			throw new Error(json.Error)
		} else {
			throw new Error(response.status)
		}
	}

	return json
}

const postJSON = async (url, body = {}) => {
	const options = {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify(body),
	}

	return await fetchJSON(url, options)
}

const getJSON = async (url) => {
	return await fetchJSON(url)
}

export {
	fetchJSON,
	postJSON,
	getJSON,
}

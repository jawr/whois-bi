const fetchWrapper = (url, options={}) => fetch(
	url,
	options,
)
	.then(r => {
		if (r.status >= 200 && r.status < 300) {
			return Promise.resolve(r.json())
		} else {
			// check if it is json
			const contentType = r.headers.get("content-type")
			if (contentType && contentType.indexOf("application/json") !== -1) {
				return r.json()
			}

			return Promise.reject(r.statusText || r.status)
		}
	})
	.then(r => {
		if (Object.keys(r).length === 1 && r.hasOwnProperty('Error')) return Promise.reject(r.Error)
		return Promise.resolve(r)
	})

export const del = (url, options={method: 'DELETE'}) => fetchWrapper(url, options)

export const get = (url, options={}) => fetchWrapper(url, options)
export const post = (url, body) => fetchWrapper(url, {
	method: 'POST', 
	headers: { 
		'content-type': 'application/json',
	}, 
	body: JSON.stringify(body),
})

export default fetchWrapper


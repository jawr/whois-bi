export default (url, options={}) => fetch(
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

			let error = new Error(r.statusText || r.status)
			error.response = r
			return Promise.reject(error)
		}
	})
	.then(r => {
		if (r.hasOwnProperty('Error')) throw r.Error
		return Promise.resolve(r)
	})

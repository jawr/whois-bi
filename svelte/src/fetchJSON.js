export const fetchJSON = async (url, options) => {
  const response = await fetch(url, options)

  if (!response.ok) {
    let error = ''
    try {
      const parsed = await response.json()
      error = parsed.error
    } catch (err) {
      error = response.statusText
    }
    throw new Error(error)
  }

  return await response.json()
}

export const postJSON = async (url, data) => {
  return await fetchJSON(url, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify(data)
  })
}

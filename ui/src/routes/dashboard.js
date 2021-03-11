import React from 'react'
import { Link } from 'react-router-dom'
import { Page } from '../components/wrapper'
import { Table } from '../components/table'
import Search from '../components/search'
import { useDomains } from '../context/domains'

const Dashboard = () => {
  const { domains, getDomains, setQuery, filter } = useDomains()
  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    getDomains()
    setLoading(false)
  }, [getDomains])

  const columns = React.useMemo(() => TableColumns, [])

  return (
    <Page loading={loading}>
      <h1 className="f3 f2-m f1-l fw2 mv3">Dashboard</h1>
      <p>Currently monitoring {domains.length} domains.</p>

      <div className="mv5">
        <Search title="Filter domains" onChange={setQuery} />
        <Table data={filter(domains)} columns={columns} />
      </div>

      <Create />

    </Page>
  )
}

const TableColumns = [
    {
      Header: 'Domain', accessor: 'Domain', 
      headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-50',
      cellClassName: 'pv2 ph3 bb b--black-20',
      Cell: ({ cell }) => (
        <Link 
          to={`/dashboard/${cell.row.original.Domain}`} 
          className="f6 link blue hover-dark-gray"
        >{cell.row.original.Domain}</Link>
      ),
    },
    {
      Header: '# Records', accessor: 'Records', 
      headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-15',
      cellClassName: 'pv2 ph3 bb tc b--black-20',
    },
    {
      Header: '# Whois', accessor: 'Whois', 
      headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-15',
      cellClassName: 'pv2 ph3 bb tc b--black-20',
    },
    {
      Header: 'Updated', accessor: 'LastUpdatedAt',
      headerClassName: 'fw6 bb b--black-20 tl pb3 pr3 bg-white w-20',
      cellClassName: 'pv2 ph3 bb b--black-20 dn dtc-ns',
    },
]

const Create = () => {
  const { createDomain } = useDomains()

  const [domains, setDomains] = React.useState('')
  const [errors, setErrors] = React.useState([])
  const [status, setStatus] = React.useState([])
  const [disabled, setDisabled] = React.useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()

    setDisabled(true)
    setStatus(['Creating... Beep Boop...'])

    const all = domains.split(/\r?\n/)

    let messages = []
    let errorMessages = []

    await Promise.all(all.map(async (domain) => {
      await createDomain(domain)
        .then(() => {
          messages = [
            ...messages,
            `Added '${domain}'.`,
          ]
        })
        .catch(error => {
          errorMessages = [...errorMessages, ''+error]
        })
    }))

    setErrors(errorMessages)
    setStatus(messages)
    setDisabled(false)

    setTimeout(() => setStatus([]), 5000)
    setTimeout(() => setErrors([]), 10000)
  }

  return (
    <div className="pa4-l">
      <form className="bg-washed-green mw8 center pa4 br2-ns ba b--black-10" onSubmit={handleSubmit}>
        <fieldset className="cf bn ma0 pa0">
          <legend className="pa0 f5 f4-ns mb3 black-80">Enter a Domain per line to start monitoring</legend>
          <div className="cf">
            <label className="clip" htmlFor="domain">Domain Name</label>
            <textarea
              className="f6 f5-l input-reset bn fl black-80 bg-white pa3 lh-solid w-100 w-75-m w-80-l br2-ns br--left-ns"
              placeholder="whois.bi"
              name="domain" 
              rows="0"
              value={domains}
              onChange={e => setDomains(e.target.value)}
              disabled={disabled}
            />
            <input
              className="f6 f5-l button-reset fl pv3 tc bn bg-animate bg-black-70 hover-bg-black white pointer w-100 w-25-m w-20-l br2-ns br--right-ns"
              type="submit"
              value="Start"
            />
          </div>
          {status.length > 0 && 
              status.map((s, i) => <p key={i} className="pa0 f5 f4-ns mb1 black-80">{s}</p>)
          }
          {errors.length > 0 && 
              errors.map((e, i) => <p key={i} className="pa0 f5 f4-ns mb1 black-80">{e}</p>)
          }
        </fieldset>
      </form>
    </div>
  )
}

export default Dashboard

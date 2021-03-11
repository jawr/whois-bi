import React from 'react'
import { useParams, useHistory } from 'react-router-dom'
import { useDomains } from '../context/domains'
import { Page } from '../components/wrapper'
import { Tabs, Menu, Tab, Panels, Panel } from '../components/tabs'
import { Records } from '../components/records'
import { Whois } from '../components/whois'
import { Jobs } from '../components/jobs'

const Domain = () => {
  const { name } = useParams()
  const { domain, getDomain, deleteDomain } = useDomains()

  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    setLoading(true)
    getDomain(name)
    setLoading(false)
  }, [getDomain, name]);

  return (
    <Page loading={loading}>
      <div className="mw8 center">
        <h1 className="f3 f2-m f1-l fw2 mv3">Details</h1>
        <p>Look in depth at '{name}'.</p>
        {domain.LastUpdatedAt && <small>Last updated {domain.LastUpdatedAt}</small>}

        <Tabs>
          <Menu>
            <Tab>Records</Tab>
            <Tab>Whois</Tab>
            <Tab>Jobs</Tab>
          </Menu>

          <Panels>
            <Panel>
              <Records domain={name} />
            </Panel>

            <Panel>
              <Whois domain={name} />
            </Panel>

            <Panel>
              <Jobs domain={name} />
            </Panel>
          </Panels>
        </Tabs>

        <ConfirmDelete deleteDomain={deleteDomain} name={name} />
      </div>
    </Page>
  )
}

const ConfirmDelete = ({ deleteDomain, name }) => {
  const history = useHistory()

  const [confirmDelete, setConfirmDelete] = React.useState(false)
  const [confirmDeleteDomain, setConfirmDeleteDomain] = React.useState('')
  const [confirmDeleteErr, setConfirmDeleteErr] = React.useState('')

  const handleDelete = (e) => {
    setConfirmDeleteErr('')

    if (confirmDelete) {
      if (confirmDeleteDomain === name) {
        deleteDomain(name)
          .catch(err => setConfirmDeleteErr(err))
          .then(() => history.push('/'))
      }
    } else {
      setConfirmDelete(true)
    }
  }

  return (
    <>
      <div className="cf">
        <input
          className="mt6 f6 f5-l button-reset fr pv3 tc dib ba bw1 b--light-red  bg-animate bg-light-red hover-bg-red white pointer w-100 w-25-m w-20-l br--right br2"
          type="submit"
          value="Delete"
          onClick={handleDelete}
        />
        {confirmDelete && <input
          className="mt6 f6 f5-l button-reset fr pv3 dib ba bw1 b--light-red tc w-100 w-40-m w-30-l br--left br2"
          type="text"
          placeholder={`please type ${name}`}
          value={confirmDeleteDomain}
          onChange={e => setConfirmDeleteDomain(e.target.value)}
        />}
      </div>

      <div className="cf">
        <p className="mt3 pr3 fr light-red">{confirmDeleteErr}</p>
      </div>
    </>
  )
}

export default Domain

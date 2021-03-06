import React, { useState, useEffect } from 'react'
import { useParams, useHistory } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { Page } from '../components/wrapper'
import { actions, selectors } from '../store/domains'
import { Tabs, Menu, Tab, Panels, Panel } from '../components/tabs'
import { Records } from '../components/records'
import { Whois } from '../components/whois'
import { Jobs } from '../components/jobs'

export default () => {
	const [loading, setLoading] = useState(true)

  const [confirmDelete, setConfirmDelete] = useState(false)
  const [confirmDeleteDomain, setConfirmDeleteDomain] = useState('')
  const [confirmDeleteErr, setConfirmDeleteErr] = useState('')

	const { name } = useParams()
	const domain = useSelector(selectors.domainByName(name))

	const dispatch = useDispatch()
  const history = useHistory()

	useEffect(() => {
		dispatch(actions.get(name)).finally(() => setLoading(false))
	}, [dispatch, name]);

  const handleDelete = (e) => {
    setConfirmDeleteErr('')

    if (confirmDelete) {
      if (confirmDeleteDomain === name) {
        dispatch(actions.delete(name))
          .catch(err => setConfirmDeleteErr(err))
          .then(() => history.push('/'))
      }
    } else {
      setConfirmDelete(true)
    }
  }

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
            <Panel classes="">
              <Records domain={domain} />
            </Panel>

            <Panel>
              <Whois domain={domain} />
            </Panel>

            <Panel>
              <Jobs domain={domain} />
            </Panel>
          </Panels>
        </Tabs>

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
      </div>

		</Page>
	)
}

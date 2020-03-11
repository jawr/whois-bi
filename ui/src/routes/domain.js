import React, { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { Page } from '../components/wrapper'
import { actions, selectors } from '../store/domains'
import { Tabs, Menu, Tab, Panels, Panel } from '../components/tabs'
import { Records } from '../components/records'
import { Whois } from '../components/whois'
import { Jobs } from '../components/jobs'

export default () => {
	const [loading, setLoading] = useState(true)

	const { name } = useParams()
	const domain = useSelector(selectors.domainByName(name))

	const dispatch = useDispatch()

	useEffect(() => {
		dispatch(actions.get(name)).finally(() => setLoading(false))
	}, [dispatch, name]);

	return (
		<Page loading={loading}>
			<h1 className="f3 f2-m f1-l fw2 mv3">Details</h1>
			<p>Look in depth at '{name}'.</p>

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

		</Page>
	)
}

import React from 'react'
import { Link } from 'react-router-dom'
import { Page } from '../components/wrapper'

export default () => (
	<Page>
			<h1 className="f3 f2-m f1-l fw2 black-90 mv3">
				Simple.
			</h1>
			<h2 className="f5 f4-m f3-l fw2 black-50 mt0 lh-copy">
				DNS & Whois monitoring has never been so straight forward.
			</h2>
			<h3 className="f7 f6-m f5-l fw2 black-50 mt0 lh-copy">
				Add a domain, select monitoring frequency and when we detect any changes we will let you know.
			</h3>

			<br />
			<Link className="f4 link dim br1 ph3 pv2 mb2 dib white bg-green grow" to="/register">Start Monitoring</Link>
	</Page>
)


import React from 'react'
import Wrapper from './wrapper'

export default () => (
	<footer className="ph3 ph4-ns pv6 bt tr b--black-10 black-70 bg-light-green">
		<Wrapper>
			<a href="mailto:hi@whois.bi" className="link b f3 f2-ns dim black-70 lh-solid">hi@whois.bi</a>
			<p className="f6 db b ttu lh-solid">© 2020 PageUp Ltd</p>
			<div className="mt5">
				<a href="/privacy/"  title="Privacy" className="f6 dib pl2 mid-gray dim">Privacy</a>
			</div>
		</Wrapper>
	</footer>
)

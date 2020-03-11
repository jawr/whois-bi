import React from 'react'
import { 
	Tabs as LibTabs, 
	Panel as LibPanel,
	useTabState,
} from '@bumaga/tabs'
import { Loader } from './wrapper'

export const Tab = ({ children, sub }) => {
	const { isActive, onClick } = useTabState()

	let classes = ["tabs__menu-item", "w-50", "tc", "bg-animate", "pointer", "green"]

	if (sub === true) {
		classes.push("hover-bg-near-white")
		classes.push("pb2")
		classes.push("pt2")

		if (isActive) {
			classes.push("bg-near-white")
		}
	} else {
		classes.push("hover-bg-washed-green")
		classes.push("pb3")
		classes.push("pt3")

		if (isActive) {
			classes.push("bg-washed-green")
		}
	}

	return (
		<label
			onClick={onClick}
			className={classes.join(' ')}
		>{ children }</label>
	)
}

export const Panel = ({ children, loading, classes }) => {
	const classOverrides = (classes !== undefined) ? classes : "mt5"
	return (
	<LibPanel>
		<div className={classOverrides}>
			{loading && <Loader />}
			{!loading && children}
		</div>
	</LibPanel>
)
}

export const Menu = ({ children, small }) => {
	let classes = ["tabs__menu", "flex", "bb", "b--black-20"]

	return (
		<div className={classes.join(' ')}>
		{children}
	</div>
	)
}

export const Panels = ({ children }) => (
	<div>
		{children}
	</div>
)

export const Tabs = ({ children, sub }) => {
	let classes = ["tabs mw8 center"]
	if (sub !== true) classes.push("mt5")

	return (
	<LibTabs>
		<div className={classes.join(' ')}>
			{children}
		</div>
	</LibTabs>
)
}

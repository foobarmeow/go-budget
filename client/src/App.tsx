import React from 'react';
import Balances from './Balance'
import Items from './Items'
import './App.css';

const App = () => {
	// Get the budget information
	

	document.title = "bank"

	//const getBudget = async () => {
	//	try {
	//		const res = await fetch("http://localhost:8080/api/budget")
	//		console.log(res.json())
	//	} catch (e) {
	//		console.error(e)
	//	}
	//}

	const budget = {
		balances: [
			{
				name: "WF",
				balance: 256.52,
			},
			{
				name: "DF",
				balance: 789.52,
			},
		],
		items: [
			{
				name: 'Phone',
				amount: 200,
				active: true,
				date: 23,
			},
			{
				name: 'Car',
				amount: 200,
				active: true,
				date: 23,
			},
		],
	}

	const total = budget.balances.reduce((a, b) => {
		return a + b.balance
	}, 0)

	const earmarked = budget.items.reduce((a, b) => {
		return a + b.amount
	}, 0)

	const available = total - earmarked

	return (
		<div>
			<Balances balances={budget.balances} />
			<p>Total: ${total}</p>
			<p>Available: ${available}</p>
			<Items items={budget.items} />
		</div>
	);
}

export default App
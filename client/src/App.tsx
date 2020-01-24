import React, {useState, useEffect} from 'react';
import axios, { AxiosResponse, AxiosPromise } from 'axios'
import Balances from './Balance'
import ItemComponent, {Item} from './Item'
import './App.css';

interface Balance {
	name: string;
	balance: number;
}

interface Budget {
	balances: Balance[];
	items: Item[];
}

type Response<Success> = AxiosPromise<Success>

const path = "http://localhost:9042/api"

const App = () => {
	// Get the budget information
	document.title = "bank"

	const [budget, setBudget] = useState({balances: [], items: []} as Budget)

	const getBudget = async () => {
		const budget = await axios({
			url: `${path}/budget`,
			method: "GET",
			withCredentials: true,
		});
		setBudget(budget.data)
	}

	useEffect(() => {
		getBudget()
	}, [])

	const total = budget.balances.reduce((a, b) => {
		return a + b.balance
	}, 0)

	const earmarked = budget.items.reduce((a, b) => {
		return a + b.amount
	}, 0)

	const available = total - earmarked

	const items = budget.items.map(i => {
		return (
			<ItemComponent item={i}/>
		)
	})

	return (
		<div className="container">
			<Balances balances={budget.balances} />
			<p>Total: ${total}</p>
			<p>Available: ${available}</p>
			<p>Earmarked: ${earmarked}</p>
			{items}
		</div>
	);
}

export default App
	//const budget = {
	//	balances: [
	//		{
	//			name: "WF",
	//			balance: 256.52,
	//		},
	//		{
	//			name: "DF",
	//			balance: 789.52,
	//		},
	//	],
	//	items: [
	//		{
	//			_id: '1',
	//			name: 'Phone',
	//			amount: 200,
	//			active: true,
	//			date: 23,
	//		},
	//		{
	//			_id: '2',
	//			name: 'Car',
	//			amount: 200,
	//			active: true,
	//			date: 23,
	//		},
	//	],
	//}
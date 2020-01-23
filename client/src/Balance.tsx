import React from 'react';
import './App.css'

export interface Balance {
    name: string;
    balance: number;
}

export interface BalanceProps {
    balances: Balance[];
}

const BalanceComponent: React.FunctionComponent<BalanceProps>  = ({ balances }) => {
    const bs = balances.map(b => {
        return (
            <div className='balance'>{b.name} - ${b.balance}</div>
        );
    })
    return (
        <>
            {bs}
        </>
    )
}

export default BalanceComponent
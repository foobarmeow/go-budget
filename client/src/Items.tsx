import React from 'react';

export interface Item {
    name: string;
    amount: number;
    active: boolean;
    date: number;
}

export interface ItemProps {
    items: Item[];
}

const Items: React.FC<ItemProps> = ({ items }) => {

    const itemElements = items.map(({name, amount, active, date}) => {
        const className = active ? 'item' : 'inactiveItem'
        return (
            <div className='${item}'>{name} - {amount} - {date}</div>
        )
    })

    return (
        <>{itemElements}</>
    )

}

export default Items
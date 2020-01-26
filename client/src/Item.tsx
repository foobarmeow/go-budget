import React, {useState } from 'react';
import axios, { AxiosPromise } from 'axios'
import { Toggle } from 'react-toggle-component'
import Editable from './editable'


const path = "http://localhost:9042/api"

export interface Item {
    _id: string;
    name: string;
    amount: number;
    active: boolean;
    date: number;
}

export interface ItemProps {
    item: Item
    onChange: (item: Item) => void;
}

const ItemComponent: React.FC<ItemProps> = (props) => {
    const [item, setItem] = useState(props.item)

    const onToggle = (e: React.SyntheticEvent) => {
        const target = e.target as HTMLInputElement
        item.active = target.checked
        setItem(item)
        props.onChange(item)
    }
    
    const onEdit = (key: string, value: string | number) => {
        switch (key) {
            case "name":
                item.name = value as string;
                break
            case "amount":
                item.amount = value as number;
                break
            case "date":
                item.date = value as number;
                break
        }

        setItem(item)
        props.onChange(item)
    }

    return (
        <div className='item'>
            <Editable name={"name"} initialValue={item.name}  type="string" onChange={onEdit}/>
            <Editable name={"amount"} initialValue={item.amount} type="number" onChange={onEdit}/>
            <Editable name={"date"} initialValue={item.date} type="number" onChange={onEdit}/>
            <Toggle checked={item.active} onToggle={onToggle}/>
        </div>
    )
}

export default ItemComponent
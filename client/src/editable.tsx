import React, { useState } from 'react';


export interface EditableProps {
    name: string;
    type: string;
    initialValue: string | number;
    onChange: (name: string, value: string | number) => void;
}

const Editable: React.FC<EditableProps> = ({ name, type, initialValue, onChange }) => {
    const [value, setValue] = useState(initialValue)
    const [editing, setEditing] = useState(false)
    const focus = () => {
        setEditing(!editing)
    }

    const blur = (e: React.SyntheticEvent) => {
        var value: string | number;
        const target = e.target as HTMLInputElement;

        switch (type) {
            case "number":
                value = parseFloat(target.value)
                break;
            default: 
                value = target.value;
                break;
        }

        setValue(value)
        setEditing(false)
        onChange(name, value)
    }

    const keyDown = (e: React.KeyboardEvent) => {
        if (e.keyCode === 13) blur(e)
    }

    if (editing) {
        return (
            <input className="editable-write" type="text" defaultValue={value} onBlur={blur} onKeyDown={keyDown}/>
        )
    }

    return (
        <div className="editable-read" onClick={focus}>{value}</div>
    )
}

export default Editable
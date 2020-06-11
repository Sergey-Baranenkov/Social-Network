import React from "react";
import "../../scss/default_blocks.scss"

export function Line({children}) {
    return (
        <div className={"default_container default_line"}>
            {children}
        </div>
    )
}

export function Column({children}) {
    return (
        <div className={"default_container default_column"}>
            {children}
        </div>
    )
}

export function SmallTextField({header,...props}) {
    return (
        <div className={"default_field"}>
            <p className={"default_field__label"}>{header}</p>
            <input className={"default_field__input"} {...props}/>
        </div>
    )
}

export function SelectField({header, options=[], value = "", ...props}){
    return (
        <div className={"default_field"}>
            <p className={"default_field__label"}>{header}</p>
            <select className={"default_field__input"} value = {value} {...props}>
                {options.map(option=><option key = {option} >{option}</option>)}
            </select>
        </div>
    )
}

export function LargeTextField({header, ...props}) {
    return (
        <div className={"default_field"}>
            <p className={"default_field__label"}>{header}</p>
            <textarea className={"default_field__input"} {...props}/>
        </div>
    )
}

export const SmallFileField = React.forwardRef(({header, ...props}, ref)=>(
    <div className={"default_field"}>
        <p className={"default_field__label"}>{header}</p>
        <input className="default_field__input" type={"file"} {...props} ref = {ref}/>
    </div>
));
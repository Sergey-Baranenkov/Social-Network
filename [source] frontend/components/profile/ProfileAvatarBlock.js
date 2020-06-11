import React from "react";
import "./profile_avatar_block.scss"
export default function ProfileAvatarBlock(props){
    return (
        <div className={"profile_avatar_block"}>
            <img src={props.src} className={"avatar"} alt = " "/>
            <p className={"name"} style={{color:props.nameColor}}>{props.name}</p>
            <p className={"description"} style={{color:props.descriptionColor}}>{props.description}</p>
        </div>
    )
}
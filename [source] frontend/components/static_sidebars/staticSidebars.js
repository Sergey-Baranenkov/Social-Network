import "./menu__static.scss"
import "../../scss/default_blocks.scss"
import emblem from "../../images/emblem.jpg";
import ProfileAvatarBlock from "../profile/ProfileAvatarBlock";
import React, {useRef, useState, memo} from "react";
import {Link, useHistory} from "react-router-dom";
import HeaderAudioPlayer from "./player/HeaderAudioPlayer";
import PeopleSearchField from "./searchfield/PeopleSearchField";
import PathFromIdGenerator from "../../functools/PathFromIdGenerator";
import {PlayerContext} from "../../PlayerContext";
import getCookie from "../../functools/getCookie";
import {ADDR, HTTP} from "../../address";

const StaticSidebars = memo(({audioInfo, changeTrackIndex}) => {
    return (
        <>
            <HeaderTop audioInfo ={audioInfo} changeTrackIndex={changeTrackIndex}/>
            <LeftSidebar/>
            <RightSidebar/>
        </>
    )
});

function clearListCookies(){
    let cookies = document.cookie.split(";");
    for (let i = 0; i < cookies.length; i++){
        let spcook =  cookies[i].split("=");
        document.cookie = spcook[0] + "=;expires=Thu, 21 Sep 1979 00:00:01 UTC;";
    }
}


function HeaderTop() {
    const history = useHistory();
    const myId = +getCookie("userId");
    const firstName = getCookie("firstName");
    const lastName = getCookie("lastName");
    const exitButtonClickedHandler = ()=>{
        clearListCookies();
        history.push("/авторизация");
    }

    return (
        <div className={"header_top"}>
            <img className={"logo header__logo"} src={emblem} alt={"logo"}/>
            <div className={"header_top__container"}>
                <PeopleSearchField/>

                <PlayerContext.Consumer>
                    {HeaderAudioPlayer}
                </PlayerContext.Consumer>

                <ProfileAvatarBlock
                    src = {HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(myId)}/profile_avatar.jpg`}
                    name = {`${firstName} ${lastName}`}
                    nameColor = "whitesmoke"
                    description = "online"
                    descriptionColor = "cadetblue"
                />
            </div>
            <button className={"exit_button"} onClick={exitButtonClickedHandler}/>
        </div>
    )
}

function LeftSidebar() {
    const [left_sidebar_isOpened, change_left_sidebar_state] = useState(false)
    const leftSidebarRef = useRef(null)
    const myId = getCookie("userId");
    const showMenu = ()=>{
        left_sidebar_isOpened ?
            leftSidebarRef.current.classList.remove("show_menu__effect"):
            leftSidebarRef.current.classList.add("show_menu__effect");
        change_left_sidebar_state(!left_sidebar_isOpened)
    };

    return(
        <div className="sidebar_left" ref={leftSidebarRef}>
            <button className={"sidebar_left__item"} onClick={showMenu}>
                <i className={"show_menu__button default_img"}/>
                <span>Скрыть меню</span>
            </button>

            <Link
                className="sidebar_left__item"
                role="button"
                to={"/профиль/" + myId}
            >
                <i className={"show_my_profile__button default_img"}/>
                <span>Мой профиль</span>
            </Link>

            <Link
                className="sidebar_left__item"
                role="button"
                to= {"/фотографии/" + myId}
            >
                <i className={"show_photos__button default_img"}/>
                <span>Фотографии</span>
            </Link>

            <Link
                className="sidebar_left__item"
                role="button"
                to="/сообщения/"
            >
                <i className={"show_messages__button default_img"}/>
                <span>Сообщения</span>
            </Link>

            <Link
                className="sidebar_left__item"
                role="button"
                to={`/связи/${myId}/друзья`}
            >
                <i className={"show_friends__button default_img"}/>
                <span>Друзья</span>
            </Link>

            <Link
                className="sidebar_left__item"
                role="button"
                to={"/музыка/" + myId}
            >
                <i className={"show_music__button default_img"}/>
                <span>Музыка</span>
            </Link>

            <Link
                className="sidebar_left__item"
                role="button"
                to={"/видео/"+ myId}
            >
                <i className={"show_video__button default_img"}/>
                <span>Видео</span>
            </Link>

            <Link
                className="sidebar_left__item"
                role="button"
                to="/погода"
            >
                <i className={"show_weather__button default_img"}/>
                <span>Погода</span>
            </Link>


            <Link
                className="sidebar_left__item"
                role="button"
                to="/настройки"
            >
                <i className={"show_settings__button default_img"}/>
                <span>Настройки</span>
            </Link>

        </div>
    )
}

function RightSidebar() {
    return (
        <div className="sidebar_right"/>
    )
}

export default StaticSidebars;
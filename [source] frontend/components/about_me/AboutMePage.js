import React from "react";
import {withRouter} from "react-router-dom";
import "../../scss/page.scss";
import "./about_me_page.scss";
import "../../scss/default_blocks.scss";
import Fetcher from "../../functools/Fetcher";
import {femaleStatuses, maleStatuses} from "./sexStatuses"
import {HTTP, ADDR} from "../../address";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";

class AboutMePage extends React.Component{
    state = {error: null}
    handleError = handleError.bind(this);

    componentDidMount() {
        this.fetchData();
    }

    fetchData = async ()=>{
        const [error, response] = await Fetcher(
            HTTP + ADDR + "/about_me/select_extended_user_info",
        {userId: this.props.match.params.id},
            "GET"
        );
        if (error === null)
            this.setState(response);
        else
            this.handleError("невозможно загрузить данные");
    }
    render() {
        const s = this.state;
        return (
            <div className={"page__container"}>
                <div className={"page__header background_pic__city background_pic__city_black"}>
                    <h1 style={{color: "white"}}>Обо мне</h1>
                    <p style={{color:  "white"}}>Здесь вы можете найти информацию обо мне</p>
                    <div className={"default_img about_me_foreground"}/>
                </div>
                <h1 style={{marginBottom: '50px'}}>{s.first_name} {s.last_name}</h1>
                <div style={{display:'grid',gridTemplateColumns:'1fr 1fr', width:'100%', gridGap:'20px'}}>
                    <div className={"default_block"}>
                        <h1 className={"default_block__header"}>Общая информация</h1>
                        <div className={"default_block__item"}>
                            <h2>Страна проживания:</h2>
                            <p>{s.country || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>Город проживания:</h2>
                            <p>{s.city || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>День рождения:</h2>
                            <p>{s.birthday || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>Телефон:</h2>
                            <p>{s.tel || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>Пол:</h2>
                            <p>{s.sex === "M" ? "Мужской":"Женский"}</p>
                        </div>

                        <div className={"default_block__item"}>
                            <h2>Семейное положение:</h2>
                            <p>{s.sex === "M" ? maleStatuses[s.status] : femaleStatuses[s.status]}</p>
                        </div>
                    </div>

                    <div className={"default_block"}>
                        <h1 className={"default_block__header"}>Хобби и интересы:</h1>
                        <div className={"default_block__item"}>
                            <h2>Хобби:</h2>
                            <p>{s.hobby || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>Любимые книги:</h2>
                            <p>{s.fav_books || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>Любимые фильмы:</h2>
                            <p>{s.fav_films || "Не указано"}</p>
                        </div>
                        <div className={"default_block__item"}>
                            <h2>Любимые игры:</h2>
                            <p>{s.fav_games || "Не указано"}</p>
                        </div>

                        <div className={"default_block__item"}>
                            <h2>Любимая музыка:</h2>
                            <p>{s.fav_music || "Не указано"}</p>
                        </div>

                        <div className={"default_block__item"}>
                            <h2>Другие интересы:</h2>
                            <p>{s.other_interests || "Не указано"}</p>
                        </div>
                    </div>
                </div>

                <div className={"default_block"}>
                    <h1 className={"default_block__header"}>Образование и работа</h1>
                    {Array.isArray(s.edu_and_emp_info) && s.edu_and_emp_info.map((info, i)=>(
                        <div className={"default_block__item"} key={i}>
                            <h2>Название:</h2>
                            <p>{info.title}</p>

                            <h2>Период:</h2>
                            <p>{info.period}</p>

                            <h2>Описание:</h2>
                            <p>{info.description}</p>
                        </div>
                    ))}
                </div>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
    handleClose = handleClose.bind(this);
}

export default withRouter(AboutMePage);
import React from "react";
import {Line, LargeTextField} from "../SettingsLayout";
import Loading from "../../loader/LoadingPage";
import Fetcher from "../../../functools/Fetcher";
import {HTTP, ADDR} from "../../../address";
import InfoPopup, {handleError, handleClose} from "../../infoPopup/infoPopup";
export default class HobbyAndInterestsPage extends React.Component{
    state ={
        isFetching: true,
        error: null,
        hobby: null,
        fav_music: null,
        fav_films: null,
        fav_books: null,
        fav_games: null,
        other_interests: null,
    };
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    fetchData = async ()=>{
        this.setState({isFetching: true});
        const [error, response] = await Fetcher(HTTP + ADDR + "/settings/hobbies");
        if (error === null){
            this.setState({...response});
        }else{
            this.handleError("Невозможно получить данные с сервера");
        }
        this.setState({isFetching: false});
    }

    componentDidMount() {
        this.fetchData();
    }


    onChangeHandler = ({target})=>{
        this.setState({[target.name] : target.value} )
    };


    sendFormHandler = async (e)=> {
        e.preventDefault();

        const copy = {...this.state};
        delete copy.error;
        delete copy.isFetching;
        const [error] = await Fetcher(HTTP + ADDR + "/settings/update_hobbies",
            {}, "POST", "text", JSON.stringify(copy));
        if (error !== null){
            this.handleError("Невозможно обновить данные");
        }else{
            this.handleError("Успешно!")
        }
    };

    reload = (e)=>{
        e.preventDefault();
        window.location.reload();
    }

    render() {
        if (this.state.isFetching){
            return <Loading/>
        }else {
            return (
                <div className={"default_block"}>
                    <div className={"default_block__header"}>Хобби и интересы</div>
                    <form onSubmit={this.sendFormHandler}>
                        <Line>
                            <LargeTextField header={"Хобби"}
                                            value={this.state.hobby}
                                            name="hobby"
                                            onChange={this.onChangeHandler}
                            />

                            <LargeTextField header={"Любимая музыка"}
                                            value={this.state.fav_music}
                                            name={"fav_music"}
                                            onChange={this.onChangeHandler}
                            />
                        </Line>

                        <Line>
                            <LargeTextField
                                header={"Любимые фильмы"}
                                value={this.state.fav_films}
                                name={"fav_films"}
                                onChange={this.onChangeHandler}
                            />
                            <LargeTextField
                                header={"Любимые книги"}
                                value={this.state.fav_books}
                                name={"fav_books"}
                                onChange={this.onChangeHandler}
                            />
                        </Line>

                        <Line>
                            <LargeTextField
                                header={"Любимые игры"}
                                value={this.state.fav_games}
                                name={"fav_games"}
                                onChange={this.onChangeHandler}
                            />
                            <LargeTextField
                                header={"Другие занятия"}
                                value={this.state.other_interests}
                                name={"other_interests"}
                                onChange={this.onChangeHandler}
                            />
                        </Line>

                        <Line>
                            <button className={"default_button cancel_button"} onClick={this.reload}>Отменить</button>
                            <button className={"default_button confirm_button"}>Сохранить изменения</button>
                        </Line>

                    </form>
                    {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
                </div>
            )
        }
    }
}
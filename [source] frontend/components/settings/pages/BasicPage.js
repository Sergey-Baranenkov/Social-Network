import React from "react";
import {countries} from "../Countries";
import {Line, SelectField, SmallTextField, SmallFileField} from "../SettingsLayout";
import {femaleStatuses, maleStatuses} from "../../about_me/sexStatuses";
import Loading from "../../loader/LoadingPage";
import Fetcher from "../../../functools/Fetcher";
import {HTTP, ADDR} from "../../../address";
import InfoPopup, {handleError, handleClose} from "../../infoPopup/infoPopup";
export default class BasicPage extends React.Component{
    state ={
        isFetching: false,
        error: null,
        sex: null,
        status: null,
        birthday: null,
        tel: 0,
        country: null,
        city: null,
    };
    profile_avatar = React.createRef();
    profile_bg = React.createRef();

    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    onChangeHandler = ({target})=>{
        this.setState({[target.name] : target.value} )
    };

    fetchData = async ()=>{
        const [error, response] = await Fetcher(HTTP + ADDR + '/settings/get_basic_info');
        if (error === null){
            if (response.sex === "M"){
                response.sex = 'Мужской';
                response.status = maleStatuses[+response.status];
            }else{
                response.sex = 'Женский';
                response.status = femaleStatuses[+response.status];
            }
            console.log(response);
            this.setState(response);
        }else{
            this.handleError("Невозможно получить данные с сервера")
        }
    }
    componentDidMount() {
        this.fetchData();
    }

    changeAvatar = async () => {
        if (this.profile_avatar.current.files.length){
            const data = new FormData();
            data.append("photo", this.profile_avatar.current.files[0])
            const [error] = await Fetcher(
                HTTP + ADDR + '/settings/update_basic_info/profile_avatar',
                null,
                "post",
                "text",
                data
                );
            if (error !== null){
                this.handleError("Невозможно обновить аватар")
            }
        }
    };

    changeBg = async () => {
        if (this.profile_bg.current.files.length){
            const data = new FormData();
            data.append("photo", this.profile_bg.current.files[0])
            const [error] = await Fetcher(
                HTTP + ADDR + '/settings/update_basic_info/profile_bg',
                null,
                "post",
                "text",
                data
            );
            if (error !== null){
                this.handleError("Невозможно обновить задний план")
            }
        }
    };

    changeTextInfo = async ()=> {
        const sendData = {...this.state}
        delete sendData.error;
        delete sendData.isFetching;
        sendData.tel = +sendData.tel;
        if (isNaN(sendData.tel)){
            this.handleError("Номер телефона должен полностью состоять из цифр!");
            return;
        }
        if (sendData.sex === "Мужской"){
            sendData.sex = 'M';
            sendData.status = maleStatuses.indexOf(sendData.status);
        }else{
            sendData.sex = 'F';
            sendData.status = femaleStatuses.indexOf(sendData.status);
        }

        const [error] = await Fetcher(
            HTTP + ADDR + '/settings/update_basic_info/text_data',
            null,
            "post",
            "text",
            JSON.stringify(sendData)
        );
        if (error !== null){
            this.handleError("Невозможно обновить базовую информацию профиля")
        }else {
            this.handleError("Успешно!")
        }
    };


    sendFormHandler = (e)=>{
        e.preventDefault();
        this.changeAvatar();
        this.changeBg();
        this.changeTextInfo();
    };

    reload = (e)=>{
        e.preventDefault();
        window.location.reload();
    }
    render() {
        if (this.state.isFetching){
            return <Loading/>
        }else{
            return (
                <div className={"default_block"}>
                    <div className={"default_block__header"}>Настройка основной информации</div>

                    <form onSubmit={this.sendFormHandler}>
                        <Line>
                            <SmallFileField header={"Аватар"}
                                            accept={".jpg"}
                                            ref = {this.profile_avatar}/>

                            <SmallFileField header={"Изображение заднего плана"}
                                            accept={".jpg"}
                                            ref = {this.profile_bg}/>
                        </Line>

                        <Line>
                            <SelectField header="Пол"
                                         onChange={this.onChangeHandler}
                                         options={["Мужской","Женский"]}
                                         name ="sex"
                                         value = {this.state.sex}/>

                            <SelectField header="Статус"
                                         name = "status"
                                         onChange={this.onChangeHandler}
                                         options={
                                             this.state.sex === "Мужской"?
                                             maleStatuses :
                                             femaleStatuses}

                                         value={this.state.status}
                            />
                        </Line>

                        <Line>
                            <SmallTextField type="date"
                                            name="birthday"
                                            onChange={this.onChangeHandler}
                                            header = "Дата рождения"
                                            value = {this.state.birthday}
                            />
                            <SmallTextField type="tel"
                                            name="tel"
                                            onChange={this.onChangeHandler}
                                            header = "Телефон"
                                            value = {this.state.tel}
                            />
                        </Line>

                        <Line>
                            <SelectField header="Страна"
                                         name="country"
                                         onChange={this.onChangeHandler}
                                         options={countries}
                                         value = {this.state.country}
                            />
                            <SmallTextField type="text"
                                            name="city"
                                            onChange={this.onChangeHandler}
                                            header = "Город"
                                            value = {this.state.city}
                            />
                        </Line>

                        <Line>
                            <button className={"default_button cancel_button"} type={"button"} onClick={this.reload}>Отменить</button>
                            <button className={"default_button confirm_button"} type="submit">Сохранить изменения</button>
                        </Line>

                    </form>
                    {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
                </div>
            )
        }

    }
}
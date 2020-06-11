import React from "react";
import {SmallTextField, Line} from "../SettingsLayout";
import Fetcher from "../../../functools/Fetcher";
import {HTTP, ADDR} from "../../../address";
import InfoPopup, {handleError, handleClose} from "../../infoPopup/infoPopup";
export default class PasswordPage extends React.Component{
    state ={
        old_password: "",
        new_password: "",
        confirm_new_password: "",
        error: null,
    };
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    onChangeHandler = ({target})=>{
        this.setState({[target.name] : target.value} )
    };

    onFormSubmit = async (e) =>{
        e.preventDefault();
        if (!this.state.old_password.length || !this.state.new_password.length){
            this.handleError("Поля не заполнены");
        }else if (this.state.new_password !== this.state.confirm_new_password){
            this.handleError("Пароли не совпадают");
        }else{
            const old_password = this.state.old_password;
            const new_password = this.state.new_password;
            const [error] = await Fetcher(HTTP + ADDR + "/settings/update_password",
                null, "POST", "text", JSON.stringify({old_password, new_password}))
            if (error === null){
                this.handleError("Пароль успешно изменен");
            }else{
                this.handleError("Невозможно изменить пароль, возможно вы ввели неправильный старый пароль");
            }
        }
    };

    render() {
        return (
                <div className={"default_block"}>
                    <div className={"default_block__header"}>Сменить пароль</div>
                    <form onSubmit={this.onFormSubmit}>
                        <Line>
                            <SmallTextField
                                type = "password"
                                header={"Подтвердить текущий пароль"}
                                name = "old_password"
                                onChange = {this.onChangeHandler}
                            />
                        </Line>

                        <Line>
                            <SmallTextField type = "password"
                                            header={"Новый пароль"}
                                            name = "new_password"
                                            onChange = {this.onChangeHandler}
                            />
                            <SmallTextField type = "password"
                                            header={"Подтвердить новый пароль"}
                                            name = "confirm_new_password"
                                            onChange = {this.onChangeHandler}
                            />
                        </Line>

                        <Line>
                            <button className={"default_button confirm_button"}>Подтвердить смену пароля</button>
                        </Line>
                    </form>
                    {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
                </div>
        )
    }
}
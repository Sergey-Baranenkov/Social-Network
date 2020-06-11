import React from "react";
import "./auth_page.scss";
import Fetcher from "../../functools/Fetcher";
import {ADDR, HTTP} from "../../address";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
import { withRouter } from 'react-router-dom';

function mixHandleChange ({target}) {
    this.setState({[target.name]: target.value});
}

export default class AuthPage extends React.Component {
    state = { signIn: false, error: null }

    changeFunc = () => {
        this.setState((state) => ({ signIn: !state.signIn }));
    };

    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    render() {
        return (
            <div className="auth_page__container">
                <div className="reg_form">
                    <div className="auth_switcher">
                        <label>
                            <input
                                defaultChecked="true"
                                type="radio"
                                name="auth_switcher"
                                onChange={this.changeFunc}
                            />
                            <div>Sign Up</div>
                        </label>
                        <label>
                            <input
                                type="radio"
                                name="auth_switcher"
                                onChange={this.changeFunc}
                            />
                            <div>Log In</div>
                        </label>
                    </div>
                    {this.state.signIn ?
                        <Login handleError={this.handleError}/>
                    :
                        <Registration handleError={this.handleError}/>
                    }
                </div>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        );
    }
}

function validMail(mail)
{
    return /^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>().,;\s@"]+\.?)+([^<>().,;:\s@"]{2,}|[\d.]+))$/.test(mail);
}

const Registration = withRouter(class extends React.Component {
    state = {first_name: "", last_name:"", sex: "Мужской", email: "", password:""}
    addr = "/registration"

    handleChange = mixHandleChange.bind(this);
    handleSubmit = async (e) => {
        e.preventDefault();
        const s = {...this.state};
        if (s.first_name.length && s.last_name.length && validMail(s.email) && s.password.length){
            s.sex = s.sex === "Мужской" ? "M" : "F";
            const [error] = await Fetcher(
                HTTP + ADDR + this.addr,
                null,
                "post",
                "text",
                JSON.stringify(s)
            )
            if (error!== null)
                this.props.handleError("пользователь уже существует");
            else
                this.props.history.push("/сообщения")
        }else{
            this.props.handleError("не все поля корректно заполнены");
        }
    }

    render() {
        return (
            <div>
                <h1 className="auth_header">Регистрация</h1>
                <form className="auth_form" onSubmit={this.handleSubmit}>
                    <input
                        type="text"
                        name="first_name"
                        placeholder="Имя"
                        autoComplete="off"
                        onChange={this.handleChange}
                        value={this.state.first_name}
                    />
                    <input
                        type="text"
                        name="last_name"
                        placeholder="Фамилия"
                        autoComplete="off"
                        onChange={this.handleChange}
                        value={this.state.last_name}
                    />
                    <input
                        type="text"
                        name="email"
                        placeholder="Email"
                        autoComplete="off"
                        onChange={this.handleChange}
                        value={this.state.email}
                    />
                    <select name="sex" placeholder="Пол" onChange={this.handleChange} value={this.state.sex}>
                        <option>Мужской</option>
                        <option>Женский</option>
                    </select>
                    <input
                        type="password"
                        name="password"
                        placeholder="Пароль"
                        autoComplete="off"
                        className="full_width"
                        onChange={this.handleChange}
                        value={this.state.password}
                    />
                    <input
                        type="submit"
                        className="full_width"
                        id="signup_submit_button"
                        value="Зарегистрироваться"
                    />
                </form>
            </div>
        );
    }
})

const Login = withRouter(class extends React.Component {
    state = {email: "", password:""};
    addr = "/login";

    handleChange = mixHandleChange.bind(this);
    handleSubmit = async (e) => {
        e.preventDefault();
        const s = {...this.state};
        if (validMail(s.email) && s.password.length){
            const [error] = await Fetcher(
                HTTP + ADDR + this.addr,
                null,
                "post",
                "text",
                JSON.stringify(s)
            )
            if (error !== null)
                this.props.handleError("некорректные имя пользователя или пароль");
            else
                this.props.history.push("/сообщения")
        }else{
            this.props.handleError("не все поля корректно заполнены");
        }
    }


    render() {
        return (
            <div>
                <h1 className="auth_header">С возвращением!</h1>
                <form className="auth_form" onSubmit={this.handleSubmit}>
                    <input
                        type="text"
                        name="email"
                        placeholder="Email"
                        autoComplete="off"
                        className="full_width"
                        value={this.state.email}
                        onChange={this.handleChange}
                    />
                    <input
                        type="password"
                        name="password"
                        placeholder="Пароль"
                        autoComplete="off"
                        className="full_width"
                        value={this.state.password}
                        onChange={this.handleChange}
                    />
                    <input
                        type="submit"
                        className="full_width"
                        id="signup_submit_button"
                        value="Авторизоваться"
                    />
                </form>
            </div>
        );
    }
})
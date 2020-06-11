import React, {lazy, Suspense} from "react";
import {BrowserRouter as Router, Link, Route, Switch} from 'react-router-dom';
import Loading from "../loader/LoadingPage";
import "./settings.scss"
import "../../scss/default_blocks.scss"
import "../../scss/page.scss"

const BasicPage = lazy(()=> import("./pages/BasicPage"));
const HobbyAndInterestsPage = lazy(()=> import("./pages/HobbyAndInterestsPage"));
const EducationPage = lazy(()=> import("./pages/EducationPage"));
const PasswordPage = lazy(()=> import("./pages/PasswordPage"));

function SettingsPage(){
    return (
        <Router>
            <div className={"page__container"}>
                <div className={"page__header background_pic__city background_pic__city_orange"}>
                    <h1 style={{color: "white"}} >Настройки</h1>
                    <p  style={{color: "white"}}>Здесь вы можете настроить свою страницу</p>
                    <div className={"settings_foreground_pic default_img"}/>
                </div>

                <div className={"settings__main"}>
                    <div className={"default_block"}>
                        <div className={"default_block__header"}>
                            Настройки профиля
                        </div>
                        <Link className={"default_link"} to={"/настройки/основные"}>Основная информация</Link>
                        <Link className={"default_link"} to={"/настройки/хобби"}>Хобби и интересы</Link>
                        <Link className={"default_link"} to={"/настройки/образование"}>Образование и работа</Link>
                        <Link className={"default_link"} to={"/настройки/пароль"}>Сменить пароль</Link>

                    </div>

                    <Switch>
                        <Route path="/настройки/основные">
                            <Suspense fallback={<Loading/>}>
                                <BasicPage/>
                            </Suspense>
                        </Route>

                        <Route path="/настройки/хобби">
                            <Suspense fallback={<Loading/>}>
                                <HobbyAndInterestsPage/>
                            </Suspense>
                        </Route>

                        <Route path="/настройки/образование">
                            <Suspense fallback={<Loading/>}>
                                <EducationPage/>
                            </Suspense>
                        </Route>


                        <Route path="/настройки/пароль">
                            <Suspense fallback={<Loading/>}>
                                <PasswordPage/>
                            </Suspense>
                        </Route>

                    </Switch>
                </div>
            </div>
        </Router>
    )
}

export default SettingsPage
import React from "react";
import {WEATHER_SYMBOL, WWO_CODE} from "./weather_icon_constants";
import "./weather_page.scss";
import Debounce from "../../functools/Debounce";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";

export default class WeatherPage extends React.Component{
    state = {weatherInfo: null, error: null}
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);

    debounced = new Debounce((event)=> {
        this.getWeather(event.target.value);
    }, 1000);

    onChange = (event) => {
        event.persist()
        this.debounced(event);
    }

    componentDidMount() {
        if (!localStorage.getItem("weather_city")){
            localStorage.setItem("weather_city","Нижний Новгород")
        }
        this.getWeather(localStorage.getItem("weather_city"));
    }

    getWeather = async (city)=>{
        localStorage.setItem("weather_city", city)
        // cors за что :( 
        const url = `https://cors-anywhere.herokuapp.com/http://wttr.in/${city.replace(" ", "+")}?format=j1`;
        const options = {method: 'get', headers: {'Accept-Language':'ru-RU, ru', 'Access-Control-Allow-Origin': 'wttr.in'}};
        const response = await fetch(url, options);
        if (response.ok){
            this.setState({weatherInfo: await response.json()});
        }else {
            this.handleError("Сервер погоды не отвечает");
        }
    }

    render() {
        return (
            <div className={"page__container"}>
                <WeatherTopBlock weatherInfo={this.state.weatherInfo}/>
                <input type={"text"}
                       className={"weather_city__input"}
                       placeholder={"Введите новый город"}
                       onChange={this.onChange}
                       defaultValue={localStorage.getItem("weather_city")}
                />
                <ExtendedForecastBlock weather={this.state.weatherInfo?.weather}/>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
}


const days_of_week = ["Воскресенье", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота" ]
function ExtendedForecastBlock({weather}) {
    return (
        <div className={"extended_forecast__block"}>
            {
                weather?.map((day)=>(
                        <div key = {day.date} className={"extended_forecast__item"}>
                            <div className={"extended_forecast_date"}>{days_of_week[new Date(day.date).getDay()]}</div>
                            <div className={"extended_forecast_hours_row"}>
                                {
                                    day.hourly?.map((hour)=>(
                                        <div key={hour.time} className={"extended_forecast_hrow__item"}>
                                            <div className={"ex_forecast__hour"}>{hour.time / 100}:00</div>
                                            <div className={"default_img weather_icon"}>{WEATHER_SYMBOL[WWO_CODE[hour.weatherCode]]} </div>
                                            <div style={{fontSize:'40px'}}>{hour.tempC}°C</div>
                                            <div style={{fontSize:'20px', marginTop: "10px"}}>{hour.lang_ru[0]?.value}</div>
                                        </div>
                                    ))
                                }
                            </div>
                        </div>
                    )
                )
            }
        </div>
    )
}

function WeatherTopBlock({weatherInfo}) {
    console.log(weatherInfo);
    const getCurrentDate = ()=>(
        new Date().toLocaleString("ru", {
            month: 'long',
            day: 'numeric',
            weekday: 'long',
            timezone: 'UTC',
        })
    );
    return (
        <div className={"weather_top_block background_pic__city background_pic__city_black"}>
            <div className={"weather_date_and_place"}>
                <div className={"weather__date"}>{getCurrentDate()}</div>
                <div className={"weather__place"}>{weatherInfo?.nearest_area[0]?.areaName[0]?.value}</div>
            </div>

            <div className={"weather_content"}>
                <div className={"default_img weather_icon"}>{WEATHER_SYMBOL[WWO_CODE[weatherInfo?.current_condition[0]?.weatherCode]]}</div>
                <div className={"weather__temp"}>{weatherInfo?.current_condition[0]?.temp_C}°C</div>
                <div className={"temp__minmax"}>
                    <span>Мин: {weatherInfo?.weather[0]?.mintempC}°</span>
                    <span>Макс: {weatherInfo?.weather[0]?.maxtempC}°</span>
                </div>
                <div className={"weather__description"}>
                    {weatherInfo?.current_condition[0].lang_ru[0].value}
                </div>

                <div className={"weather__additional_info"}>
                    <div className={"weather__additional_info__block"}>
                        <div className={"default_img real_temp_feel_icon"}/>
                        <div>Ощущается как</div>
                        <div>{weatherInfo?.current_condition[0]?.FeelsLikeC}°</div>
                    </div>

                    <div className={"weather__additional_info__block"}>
                        <div className={"default_img humidity_icon"}/>
                        <div>Влажность</div>
                        <div>{weatherInfo?.current_condition[0]?.humidity}%</div>
                    </div>

                    <div className={"weather__additional_info__block"}>
                        <div className={"default_img pressure_icon"}/>
                        <div>Давление</div>
                        <div>{
                            weatherInfo?.current_condition[0]?.pressure
                            &&
                            (weatherInfo?.current_condition[0]?.pressure / 1.33333).toFixed()
                        }

                        </div>
                    </div>
                    <div className={"weather__additional_info__block"}>
                        <div className={"default_img wind_icon"}/>
                        <div>Скорость ветра</div>
                        <div>{
                            weatherInfo?.current_condition[0]?.windspeedKmph
                            &&
                            (weatherInfo.current_condition[0].windspeedKmph / 3.6).toFixed(1) } м/с
                        </div>
                    </div>
                </div>
            </div>

            <div className={"weather_last_observation_time"}>
                Обновлено в: {weatherInfo?.current_condition[0]?.localObsDateTime}
            </div>
        </div>
    )
}
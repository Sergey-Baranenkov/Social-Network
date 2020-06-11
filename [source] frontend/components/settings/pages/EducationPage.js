import React from "react";
import {Column, LargeTextField, Line, SmallTextField} from "../SettingsLayout";
import Loading from "../../loader/LoadingPage";
import Fetcher from "../../../functools/Fetcher";
import {HTTP, ADDR} from "../../../address";
import InfoPopup, {handleError, handleClose} from "../../infoPopup/infoPopup";
export default class EducationSettings extends React.Component{
    state ={
        isFetching: false,
        error: null,
        data: [],
    };
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    fetchData = async () =>{
        this.setState({isFetching: true});
        const [error, response] = await Fetcher(HTTP + ADDR + "/settings/get_edu_emp");
        if (error === null){
            console.log(response);
            this.setState({data: response});
        }else {
            this.handleError("Невозможно получить данные с сервера");
            this.setState({data: []});
        }
        this.setState({isFetching: false});
    }

    componentDidMount() {
        this.fetchData();
    }

    sendFormHandler = async (e)=> {
        e.preventDefault();
        const [error] = await Fetcher(HTTP + ADDR + "/settings/post_edu_emp",
            {}, "POST", "text", JSON.stringify(this.state.data.filter((d)=>d !== undefined)));
        if (error !== null){
            this.handleError("Невозможно обновить данные");
        }else {
            this.handleError("Успешно!")
        }
    };


    addField = ()=>{
        this.setState(({data}) => ({data: [...data, {title: '', period: '', description: ''}]}));
    };

    changeField = (key, name, value)=>{
        this.setState(({data}) => ({data: [
                ...data.slice(0,key),
                Object.assign({}, data[key], {[name]:value}),
                ...data.slice(key + 1)
            ]}));
    };

    deleteField = key => {
        this.setState(({data}) => ({data: [
                ...data.slice(0,key),
                undefined,
                ...data.slice(key + 1)
            ]}));
    };

    reload = (e)=>{
        e.preventDefault();
        window.location.reload();
    }

    render() {
        if (this.state.isFetching){
            return <Loading/>
        }else {
            return(
                <div className={"default_block"}>
                    <div className={"default_block__header"}>Образование и работа</div>
                    <form onSubmit={this.sendFormHandler}>
                        {
                            this.state.data?.map((field, i)=>
                                field && <EducationFieldBlock key={i}
                                                              formKey={i}
                                                              onChange = {this.changeField}
                                                              onDelete = {this.deleteField}
                                                              {...field}
                                />
                            )

                        }

                        <Line>
                            <button type={"button"}
                                    className={"default_button add_field__button"}
                                    onClick={this.addField}>+ Добавить новое поле
                            </button>
                        </Line>

                        <Line>
                            <button className={"default_button cancel_button"} onClick={this.reload}>Отменить</button>
                            <button className={"default_button confirm_button"} type={"submit"}>Сохранить изменения</button>
                        </Line>
                    </form>
                    {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
                </div>
            )
        }
    }
}

function EducationFieldBlock(props) {

    const onChange = ({target})=> props.onChange(props.formKey, target.name, target.value);

    const onDelete = () => {
        props.onDelete(props.formKey);
    };

    return (
        <Column>
            <Line>
                <SmallTextField header={"Название"}
                                type="text"
                                value={props.title}
                                name = "title"
                                onChange={onChange}
                />
                <button type={"button"} className={"default_button education__delete_button"} onClick={onDelete}>x</button>
                <SmallTextField header={"Период времени"}
                                type="text"
                                value={props.period}
                                name = "period"
                                onChange={onChange}
                />
            </Line>
            <Line>
                <LargeTextField header={"Описание"}
                                value={props.description}
                                name = "description"
                                onChange={onChange}
                />
            </Line>
        </Column>
    )
}
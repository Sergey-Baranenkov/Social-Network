import React, {memo} from "react";
import "./people_search_field.scss"
import {useHistory} from "react-router-dom";
import PathFromIdGenerator from "../../../functools/PathFromIdGenerator";
import Throttle from "../../../functools/Trottle";
import Fetcher from "../../../functools/Fetcher";
import Debounce from "../../../functools/Debounce";
import {HTTP, ADDR} from "../../../address";
import InfoPopup, {handleError, handleClose} from "../../infoPopup/infoPopup";
export default class PeopleSearchField extends React.Component{
    state = {people: [], focus: false, hoverFlag: false, value:"", offset: 0, Done: false, isFetching: false, error: null};
    limit = 10;
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    searchPeople = async ()=> {
        this.setState({isFetching: true});
        const {value, offset} = {...this.state};
        const limit = this.limit;
        console.log(value,offset,limit)
        const [error, response] = await Fetcher(
            HTTP + ADDR + "/search_people",
            {value, offset, limit},
            'GET',
            "json"
        )
        if (error === null){
            this.setState(s => ({
                people: [...s.people, ...response.People],
                offset: s.offset + this.limit,
                Done: response.Done
            }));
        }else{
            this.handleError("невозможно получить данные с сервера");
        }
        this.setState({isFetching: false});
    }

    onRequestChange = ({target}) => {
        this.setState({value: target.value}, ()=>this._onRequestChange());
    }

    _onRequestChange = Debounce(() =>{
        this.setState({Done: false, offset: 0, people: []}, ()=>this.searchPeople());
    }, 1000);

    _handleScroll = Throttle((difference)=> {
        if (difference < 10 &&
            !this.state.Done &&
            !this.state.isFetching) {
            this.searchPeople();
        }
        }, 500)

    handleScroll = ({target})=>{
        const difference = target.scrollTopMax - target.scrollTop;
        this._handleScroll(difference);
    }


    setFalseFocus = ()=> this.setState({focus:false})
    setTrueFocus = ()=> this.setState({focus:true})
    setTrueHoverFlag = () => this.setState({hoverFlag: true})
    setFalseHoverFlag = ()=> this.setState({hoverFlag: false})

    onBlur = () =>{
        if (!this.state.hoverFlag){
            this.setFalseFocus();
        }
    }

    render(){
        return (
            <div className={"search__container"}
                 onFocus={this.setTrueFocus}
                 onBlur={this.onBlur}
            >
                <div className={"search_input__block"}>
                    <input type="text"
                           placeholder="Поиск друзей"
                           className={"search__input"}
                           onChange={this.onRequestChange}
                    />
                </div>
                <div className={"searched_people__container"} onScroll={this.handleScroll}>
                    {
                        this.state.focus
                        &&
                        <div onMouseEnter={this.setTrueHoverFlag}
                             onMouseLeave={this.setFalseHoverFlag}
                             tabIndex={1}
                        >
                            {
                                this.state.people.map( person =>
                                    <FoundPerson {...person} key = {person.user_id}/>
                                )
                            }
                            {
                                this.state.isFetching && <span>Загрузка...</span>
                            }
                        </div>
                    }
                </div>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
}

const FoundPerson = memo(({user_id,first_name, last_name })=>{
    const history = useHistory();
    const avatar_ref = HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(user_id)}/profile_avatar.jpg`;
    const onProfileClick = () => {
        history.push(`/профиль/${user_id}`);
    }

    const onSendMessageClick = (e) =>{
        e.stopPropagation();
        history.push({
            pathname: '/сообщения/',
            user_id: user_id,
            first_name: first_name,
            last_name: last_name,
        })
    }

    return (
        <div className={"searched_person__container"} onClick={onProfileClick}>
            <img className={"default_img searched_user__avatar"} src={avatar_ref} alt={" "}/>
            <span className={"searched_user__fullname"}>{`${first_name} ${last_name}`}</span>
            <button className={"searched_user_send_mes__button default_img"} onClick={onSendMessageClick}/>
        </div>
    )
})
import React from "react";
import "../../scss/page.scss";
import "./friends_page.scss";
import PathFromIdGenerator from "../../functools/PathFromIdGenerator";
import "../../scss/default_blocks.scss";
import {HTTP, ADDR} from "../../address";
import {Link, Route, Switch, useRouteMatch, withRouter, useHistory} from 'react-router-dom';
import Fetcher from "../../functools/Fetcher";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
import Throttle from "../../functools/Trottle";
import getCookie from "../../functools/getCookie";

class FriendsPage extends React.Component{
    state = {error: null};
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);

    render(){
        return (
            <div className={"page__container"}>
                <div className={"page__header background_pic__city background_pic__city_green"}>
                    <h1 style={{color: "white"}}>Связи</h1>
                    <p style={{color:  "white"}}>Здесь находится список связей</p>
                    <div className={"friends_foreground_pic default_img"}/>
                </div>
                <div className={"friends_classifier__container"}>
                    <RelationClassifierLink
                        to={`/связи/${this.props.match.params.id}/друзья`}
                        label={"Друзья"}
                    />
                    <RelationClassifierLink
                        to={`/связи/${this.props.match.params.id}/подписчики`}
                        label={"Подписчики"}
                    />
                    <RelationClassifierLink
                        to={`/связи/${this.props.match.params.id}/подписки`}
                        label={"Подписки"}
                    />
                </div>
                <Switch>
                    <Route path="/связи/:id/друзья" render={props => <Relationships type={3} handleError={this.handleError} {...props}/>}/>
                    <Route path="/связи/:id/подписчики" render={props => <Relationships type = {2} handleError={this.handleError} {...props}/>}/>
                    <Route path="/связи/:id/подписки" render = {props => <Relationships type ={1} handleError={this.handleError} {...props}/>}/>
                </Switch>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
}

function RelationClassifierLink({ label, to, activeOnlyWhenExact}) {
    const match = useRouteMatch({
        path: to,
        exact: activeOnlyWhenExact
    });
    const style = {color: "white"}
    if (match) style.backgroundColor = "orange";
    return (
        <Link style={style}
              className={`default_link friends_classifier__item`}
              to={to}>
            {label}
        </Link>
    );
}



class Relationships extends React.Component{
    state = {
        people: [],
        offset: 0,
        isFetching: false,
        Done: false,
    }
    limit = 15;
    myId = +getCookie("userId");
    isPageMine = +this.props.match.params.id === this.myId;

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps.type !== this.props.type || this.props.match.params.id !== prevProps.match.params.id){
            this.isPageMine = +this.props.match.params.id ===  this.myId;
            this.setState({people:[], offset:0, Done: false}, ()=>{
                this.fetchData();
            })
        }
    }

    componentDidMount() {
        this.fetchData();
        window.addEventListener('scroll', this.handleScrollThrottled, true);
    }

    componentWillUnmount() {
        window.removeEventListener('scroll', this.handleScrollThrottled);
    }

    handleScrollThrottled = Throttle(()=>{
        if ((Math.abs(window.scrollY + window.innerHeight - document.documentElement.scrollHeight) < 10)
            &&
            !this.state.Done
            &&
            !this.state.isFetching
        ){
            this.fetchData();
        }
    }, 1000 )

    fetchData = async ()=>{
        this.setState({isFetching: true});
        const [error, response] = await Fetcher(HTTP + ADDR + "/relations/get_relations",
            {userId:this.props.match.params.id, mode: this.props.type, limit: this.limit, offset: this.state.offset}

            );
        if (error === null){
            this.setState(s => ({people: [...s.people, ...response.People], offset: s.offset + this.limit, Done: response.Done}));
        }else{
            this.props.handleError("невозможно получить данные с сервера");
        }
        this.setState({isFetching: false});
    }

    onAction = async (userId) => {
        const query = HTTP + ADDR + "/relations/update_relationship";
        const [error] = await Fetcher(query, {prevRelType : this.props.type, userId}, "POST","text");
        if (error === null){
            this.setState(s => ({people: s.people.filter(person => person.user_id!==userId), offset: s.offset - 1}));
        }else {
            this.props.handleError("невозможно обновить связь");
        }
    }

    render() {
        return (
            <>
                {this.state.people.map( person =>
                    <Relationship
                        key = {person.user_id}
                        {...person}
                        type = {this.props.type}
                        onAction={this.onAction}
                        isPageMine = {this.isPageMine}
                    />
                )}
                {this.state.isFetching && <span>Загрузка...</span>}
            </>
        )
    }
}

function Relationship(props) {
    const history = useHistory();

    const onAction = ()=>{
        props.onAction (props.user_id);
    }

    const goToProfile = ()=>{
        history.push("/профиль/" + props.user_id);
    }

    return (
        <div className={"rel__container"}>
            <img className={"default_img rel_user__avatar"} alt = " " src={
                HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(props.user_id)}/profile_avatar.jpg`
            }
            />
            <div className={"rel_user__vertical"}>
                <span className={"rel_user__fullname"} style={{cursor:"pointer"}} onClick={goToProfile}>
                    {`${props.first_name} ${props.last_name}`}
                </span>
                <Link className={"rel_user__send_mes"} to = {{
                    pathname: '/сообщения/',
                    user_id: props.user_id,
                    first_name: props.first_name,
                    last_name: props.last_name
                }}>Написать сообщение</Link>

            </div>
            <div className={"rel__functional"}>
                {
                    props.isPageMine &&
                        (
                        (props.type === 3 || props.type === 1) ?
                            <button className={"rel__button delete_friend__button"} onClick={onAction}>❌</button>
                            :
                            <button className={"rel__button add_friend__button"} onClick={onAction}>+</button>
                        )
                }
            </div>
        </div>
    )
}

export default withRouter(FriendsPage);
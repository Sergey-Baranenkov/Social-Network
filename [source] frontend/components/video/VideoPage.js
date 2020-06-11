import React, {createRef, useState} from "react";
import "./video_page.scss"
import "../../scss/page.scss";
import "../../scss/video_audio.scss";
import "../../scss/hidden_input.scss"
import "../../scss/default_popup.scss"

import {withRouter, NavLink as Link} from "react-router-dom"
import Fetcher from "../../functools/Fetcher";
import Debounce from "../../functools/Debounce";
import Throttle from "../../functools/Trottle";
import getCookie from "../../functools/getCookie";
import PathFromIdGenerator from "../../functools/PathFromIdGenerator";
import {HTTP, ADDR} from "../../address";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
class VideoPage extends React.Component{
    state = { UserVideos: [], AllVideos : [], offset: 0, searchInputValue: "", isFetching: false, popupOpen: false, Done: false, error: null}
    pageId = + this.props.match.params.id;
    isPageMine = +getCookie("userId") === this.pageId;
    fetchLimit = 16;
    addVideoRef = createRef()
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    fetchVideo = async () => {
        const offset = this.state.offset;
        const val = this.state.searchInputValue;
        const address = HTTP + ADDR + `/video/${val === ""? "get_user_video" : "get_combined_video" }`;
        const params = {userId: this.pageId, offset: offset, withValue: val, limit: this.fetchLimit}
        this.setState({isFetching: true});
        const [error, response] = await Fetcher(
            address,
            params
        )
        if (error === null){
            this.setState(state => ({
                offset: offset + this.fetchLimit,
                UserVideos: [...state.UserVideos, ...response.UserVideos],
                AllVideos: [...state.AllVideos, ...response.AllVideos],
                Done: response.Done
            }));
        }else {
            this.handleError("невозможно получить видео с сервера");
        }
        this.setState({isFetching: false});
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps !== this.props){
            this.pageId = + this.props.match.params.id;
            this.isPageMine = +getCookie("userId") === this.pageId;

            this.setState({ UserVideos: [],
                AllVideos : [],
                offset: 0,
                searchInputValue: "",
                isFetching: false,
                popupOpen: false,
                Done: false,
                error: null
            }, ()=> this.fetchInitialVideos());
        }
    }

    fetchInitialVideos = ()=>{
        this.setState(
            {offset: 0, Done: false, UserVideos: [], AllVideos : []},
            () => this.fetchVideo()
        );
    }

    componentDidMount(){
        this.fetchInitialVideos();
        window.addEventListener('scroll', this.handleScrollThrottled, true);
    }

    componentWillUnmount() {
        window.removeEventListener('scroll', this.handleScrollThrottled);
    }

    handleScrollThrottled = Throttle(() => {
        if ((Math.abs(window.scrollY + window.innerHeight - document.documentElement.scrollHeight) < 10)
            &&
            !this.state.Done
            &&
            !this.state.popupOpen
            &&
            !this.state.isFetching
        ){
            this.fetchVideo();
        }
    }, 1000);

    
    fetchInitialVideosDebounced = Debounce(()=>{
        this.setState(
            {UserVideos: [], AllVideos : [], offset: 0},
            () => this.fetchInitialVideos()
        );
    }, 500);

    
    onPopupClose = ()=>{
        this.setState({popupOpen: false});
    }

    onFileDownloaded = ()=>{
        if (this.addVideoRef.current.files.length) {
            this.setState({popupOpen: true});
        }
    }

    postVideo = async (title)=>{
        this.onPopupClose();
        const data = new FormData();
        data.append("video", this.addVideoRef.current.files[0])
        const [error, response] = await Fetcher(
            HTTP + ADDR + '/video/post_video',
            {title},
            'POST',
            "text",
            data);

        if (error === null){
            this.setState(state => (
                {UserVideos:
                    [
                        {
                            video_id: +response,
                            adder_id: this.myId,
                            name: title
                        },
                        ...state.UserVideos
                    ],
                    offset: state.offset + 1
            }))
        }else{
            this.handleError("невозможно отправить видео на сервер");
        }
    }

    addToPlaylistHandler = async (videoId) => {
        const [error] = await Fetcher(
            HTTP + ADDR + "/video/add_to_playlist",
            {videoId},
            "POST",
            "text"
        )
        if (error !== null){
            this.handleError("невозможно добавить видео в плейлист");
        }
    }

    deleteVideo = async(videoId) =>{
        const [error] = await Fetcher(
            HTTP + ADDR + "/video/remove_video",
            {videoId},
            "GET",
            "text"
        )
        if (error === null){
            this.setState(state => ({
                    offset: this.offset - 1,
                    UserVideos: state.UserVideos.filter((video)=>video.video_id !== videoId)}
            ));
        }else {
            this.handleError("невозможно удалить видео из плейлиста");
        }
    }


    onInputChange = (event)=> {
        this.setState({searchInputValue: event.target.value});
        this.fetchInitialVideosDebounced()
    }

    render() {
        return (
            <div className={"page__container"}>
                <div className={"page__header background_pic__city background_pic__city_blue"}>
                    <h1 style={{color: "white"}}>Видео</h1>
                    <p style={{color:  "white"}}>Здесь вы можете посмотреть видео</p>
                    <div className={"video_foreground_pic default_img"}/>
                </div>
                <div className={"VA_func_block"}>
                    <input className={"default_search_input"}
                           placeholder={"Поиск видео"}
                           onChange={this.onInputChange}

                    />
                    <input
                        type={"file"}
                        className={"hidden_input"}
                        id = {"select_video_file__input"}
                        accept={"video/*"}
                        onChange={this.onFileDownloaded}
                        ref = {this.addVideoRef}
                    />
                    {
                        this.isPageMine &&
                        <label htmlFor="select_video_file__input">+</label>
                    }
                </div>
                <div className={"videos__container"}>
                    {
                        this.state.UserVideos.map((about )=> (
                            <Video key = {about.video_id}
                                   about = {about}
                                   isPageMine = {this.isPageMine}
                                   videoId = {about.video_id}
                                   deleteMusicHandler = {this.deleteVideo}
                            />
                        ))
                    }
                </div>
                {Boolean(this.state.AllVideos.length) && <span>Все видео</span>}
                <div className={"videos__container"}>
                    {
                        this.state.AllVideos.map((about )=> (
                            <Video key = {about.video_id}
                                   about = {about}
                                   videoId = {about.video_id}
                                   addToPlaylistHandler = {this.addToPlaylistHandler}
                            />
                        ))
                    }
                </div>
                {this.state.isFetching && <span>Загрузка...</span>}
                {this.state.popupOpen && <AddVideoPopup onClose={this.onPopupClose} _onSend={this.postVideo}/>}
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
}

function Video({about, isPageMine, videoId, addToPlaylistHandler, deleteMusicHandler}) {
    const [buttonActive, updateButtonActive] = useState(true);
    const [hover, updateHover] = useState(false);

    const onFocus = ()=>{
        updateHover(true);
    }

    const onAdd = (e)=>{
        e.stopPropagation();
        addToPlaylistHandler(videoId);
        updateButtonActive(false);
    }

    const onDelete = (e) =>{
        e.stopPropagation();
        deleteMusicHandler(videoId);
    }

    const onBlur = ()=>{
        updateHover(false);
    }

    return (
        <div className={"video__container"}>
            <div className={"video__header"}>
                <span className={"video_name"}>{about.name}</span>
            </div>
            <video
                tabIndex={0}
                onFocus={onFocus}
                onBlur={onBlur}
                onMouseEnter={onFocus}
                onMouseLeave={onBlur}
                className={"video"}
                src={HTTP + ADDR + `/video_storage${PathFromIdGenerator(about.video_id)}/video.mp4`}
                preload={"metadata"}
                controls = {hover}
            />
            <span className={"video_functional"}>
                {
                    isPageMine ?
                        <button onClick={onDelete}>-</button>
                        :
                        <button
                            disabled={!buttonActive}
                            onClick={onAdd}
                        >{buttonActive ? '+' : '✓'}
                        </button>
                }
                <Link to = {`/профиль/${about.adder_id}`}>Страница автора ⇢</Link>
            </span>
        </div>
    )
}

function AddVideoPopup({onClose, _onSend}) {
    const [title, _updateTitle] = useState("");

    const onSend = () =>{
        _onSend(title);
    }

    const updateTitle = ({target})=>{
        _updateTitle(target.value);
    }
    return (
        <div className="b-popup">
            <div className="b-popup-content">
                <button className={"popup_button__close"} onClick={onClose}>×</button>
                <input placeholder={"Название"} value={title} onChange={updateTitle}/>
                <button className={"popup_button__action"} onClick={onSend}>Загрузить</button>
            </div>
        </div>
    )
}

export default withRouter(VideoPage);
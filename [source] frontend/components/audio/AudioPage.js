import React, {createRef, useState, memo} from "react";
import "./audio.scss"
import "../../scss/video_audio.scss";
import "../../scss/default_popup.scss"
import Debounce from "../../functools/Debounce";
import Throttle from "../../functools/Trottle";
import Fetcher from "../../functools/Fetcher";
import "../../scss/hidden_input.scss"
import {withRouter} from "react-router-dom"
import track_cover from "../../images/track_cover.ico"
import {HTTP,  ADDR} from "../../address";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
import getCookie from "../../functools/getCookie";

class AudioPage extends React.Component{
    state = {
        UserMusic: [],
        AllMusic : [],
        offset: 0,
        searchInputValue: "",
        isFetching: false,
        popupOpen: false,
        Done: false,
        error: null
    }
    pageId = + this.props.match.params.id;
    myId = + getCookie("userId");
    isPageMine = this.myId === this.pageId;

    fetchLimit = 20;
    addMusicRef = createRef();
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    fetchMusic = async () => {
        const offset = this.state.offset;
        const val = this.state.searchInputValue;
        const address = HTTP + ADDR + `/music/${val === ""? "get_user_music" : "get_combined_music" }`;
        const params = {userId: this.pageId, offset: offset, withValue: val, limit: this.fetchLimit}
        this.setState({isFetching: true});
        const [error, response] = await Fetcher(
            address,
            params
        )
        if (error === null){
            this.setState(state => ({
                offset: offset + this.fetchLimit,
                UserMusic: [...state.UserMusic, ...response.UserMusic],
                AllMusic: [...state.AllMusic, ...response.AllMusic],
                Done: response.Done
            }));
        }else {
            this.handleError("ошибка при получении данных с сервера");
        }
        this.setState({isFetching: false});
    }

    fetchInitialMusic = ()=> {
        this.setState(
            {offset: 0, Done: false, UserMusic: [], AllMusic : []},
            () => this.fetchMusic()
        );
    }

    componentDidMount(){
        this.fetchInitialMusic();
        window.addEventListener('scroll', this.handleScrollThrottled, true);
    }

    componentWillUnmount() {
        window.removeEventListener('scroll', this.handleScrollThrottled);
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps.match.params.id !== this.props.match.params.id){
            this.pageId = + this.props.match.params.id;
            this.isPageMine = this.myId === this.pageId;
            this.setState({UserMusic: [],
                AllMusic : [],
                offset: 0,
                searchInputValue: "",
                isFetching: false,
                popupOpen: false,
                Done: false,
                error: null}, ()=> this.fetchInitialMusic());
        }
    }

    handleScrollThrottled = Throttle(()=>{
        console.log(this.state.Done, this.state.popupOpen, this.state.isFetching);
        if ((Math.abs(window.scrollY + window.innerHeight - document.documentElement.scrollHeight) < 10)
            &&
            !this.state.Done
            &&
            !this.state.popupOpen
            &&
            !this.state.isFetching
        ){
            this.fetchMusic();
        }
    }, 1000 )


    fetchInitialMusicDebounced = Debounce(()=>{
        this.setState(
            {UserMusic: [], AllMusic : [], offset: 0},
            () =>this.fetchInitialMusic()
        );
    }, 500);


    onPopupClose = ()=>{
        this.setState({popupOpen: false});
    }

    onFileDownloaded = ()=>{
        if (this.addMusicRef.current.files.length) {
            this.setState({popupOpen: true});
        }
    }

    postMusic = async (author, title)=>{
        this.onPopupClose();
        const data = new FormData();
        data.append("audio", this.addMusicRef.current.files[0])
        const [error, response] = await Fetcher(
            HTTP + ADDR + "/music/post_music",
            {author, title},
            'POST',
            "text",
            data);

        if (error === null){
            this.setState(state => (
                {UserMusic:
                    [
                        {
                            music_id: response,
                            adder_id: this.myId,
                            author: author,
                            name: title
                        },

                        ...state.UserMusic
                    ],
                offset: state.offset + 1
            }))
        }else{
            this.handleError("ошибка отправки данных на сервер");
        }
    }

    addToPlaylistHandler = async (music_id) => {
       const [error] = await Fetcher(
            HTTP + ADDR + "/music/add_to_playlist",
            {musicId: music_id},
            "POST",
            "text"
        )
        if (error !== null){
            this.handleError("ошибка при добавлении музыки в плейлист");
        }
    }

    deleteMusic = async(musicId) =>{
        const [error] = await Fetcher(
            HTTP + ADDR + "/music/remove_music",
            {musicId},
            "GET",
            "text"
        )
        if (error === null){
            this.setState(state => ({
                offset: this.offset - 1,
                UserMusic: state.UserMusic.filter((track)=>track.music_id !== musicId)}
                ));
        }else {
            this.handleError("ошибка при удалении музыки");
        }
    }

    playUserMusicHandler = (offset) => {
        this.props.updatePlayerStoreState({
            playlist: this.state.UserMusic,
            trackIndex: offset,
            done: false,
            playlistName: "UserMusic",
            withValue: this.state.searchInputValue,
            userId: this.pageId
        })
    }

    playAllMusicHandler = (offset) => {
        this.props.updatePlayerStoreState({
            playlist: this.state.AllMusic,
            trackIndex: offset,
            done: false,
            playlistName: "AllMusic",
            withValue: this.state.searchInputValue,
            userId: this.pageId
        })
    }

    onInputChange = (event)=> {
        this.setState({searchInputValue: event.target.value});
        this.fetchInitialMusicDebounced()
    }

    render() {
        return (
            <div className={"page__container"}>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
                <div className={"page__header background_pic__city background_pic__city_green"}>
                    <h1 style={{color: "white"}}>Музыка</h1>
                    <p style={{color:  "white"}}>Здесь вы можете послушать музыку</p>
                    <div className={"audio_foreground_pic default_img"}/>
                </div>
                <div className={"VA_func_block"}>
                    <input className={"default_search_input"}
                           placeholder={"Поиск музыки"}
                           value={this.state.searchInputValue}
                           onChange={this.onInputChange}/>
                    {
                        this.isPageMine &&
                        <label htmlFor="select_audio_file__input">+</label>
                    }
                    <input
                        type={"file"}
                        className={"hidden_input"}
                        id = {"select_audio_file__input"}
                        accept={"audio/*"}
                        ref = {this.addMusicRef}
                        onChange={this.onFileDownloaded}
                    />

                </div>
                {
                    this.state.UserMusic.map( (track, offset) =>
                        <UserTrack
                            key = {track.music_id}
                            offset = {offset}
                            playMusicHandler = {this.playUserMusicHandler}
                            deleteMusicHandler = {this.deleteMusic}
                            name = {track.name}
                            author = {track.author}
                            music_id = {track.music_id}
                            isPageMine = {this.isPageMine}
                        />
                    )
                }
                {this.state.AllMusic.length ? <div>Все аудиозаписи</div> : null}
                {
                    this.state.AllMusic.map((track,offset) =>
                        <UserTrack
                            key = {track.music_id}
                            offset = {offset}
                            name = {track.name}
                            author = {track.author}
                            music_id = {track.music_id}
                            playMusicHandler = {this.playAllMusicHandler}
                            addToPlaylistHandler = {this.addToPlaylistHandler}
                        />
                    )
                }
                {this.state.isFetching && <span>Загрузка...</span>}
                {this.state.popupOpen && <AddAudioPopup onClose={this.onPopupClose} _onSend={this.postMusic}/>}
            </div>

        )
    }
}

const UserTrack = memo(({isPageMine, playMusicHandler, deleteMusicHandler, offset, name, author, music_id, addToPlaylistHandler})=>{
    const [buttonActive, updateButtonActive] = useState(true);

    const onPlay = () =>{
        playMusicHandler(offset)
    }

    const onDelete = (e) =>{
        e.stopPropagation();
        deleteMusicHandler(music_id);
    }

    const onAdd = (e)=>{
        e.stopPropagation();
        addToPlaylistHandler(music_id);
        updateButtonActive(false);
    }


    return (
        <div className={"audio_track"} onClick={onPlay}>
            <div className={"audio_info"}>
                <img className={"audio_cover"} src={track_cover} alt = {" "}/>
                <span className={"audio_track__name"}>{name}</span>
                <span className={"audio_track__author"}>{author}</span>
            </div>
            {isPageMine ?
                <div>
                    <button className={"audio__button audio__button-delete"}
                            onClick={onDelete}>❌</button>
                </div>
                :
                <div>
                    <button className={"audio__button audio__button-add"}
                            onClick={onAdd}
                            disabled={!buttonActive}
                    >{buttonActive ? '+' : '✓'}</button>
                </div>

            }
        </div>
        )
});

function AddAudioPopup({onClose, _onSend}) {
    const [author, _updateAuthor] = useState("");
    const [title, _updateTitle] = useState("");

    const onSend = () =>{
        _onSend(author, title);
    }
    const updateAuthor = ({target})=>{
        _updateAuthor(target.value);
    }
    const updateTitle = ({target})=>{
        _updateTitle(target.value);
    }
    return (
        <div className="b-popup">
            <div className="b-popup-content">
                <button className={"popup_button__close"} onClick={onClose}>×</button>
                <input placeholder={"Автор"} value={author} onChange={updateAuthor}/>
                <input placeholder={"Название"} value={title} onChange={updateTitle}/>
                <button className={"popup_button__action"} onClick={onSend}>Загрузить</button>
            </div>
        </div>
    )
}

export default withRouter(AudioPage);
import React, {memo, useState, useRef} from "react";
import ProfileAvatarBlock from "./ProfileAvatarBlock";
import {withRouter, NavLink as Link} from "react-router-dom"
import AddNewMessageArea from "../addNewMessageArea/addNewMessageArea";
import Fetcher from "../../functools/Fetcher";
import "./profile_page.scss"
import "../../scss/default_blocks.scss"
import "../../scss/page.scss"
import PathFromIdGenerator from "../../functools/PathFromIdGenerator";
import Debounce from "../../functools/Debounce";
import getCookie from "../../functools/getCookie";
import {oldNewRelMatchDict} from "../../oldNewRelMatchDict";
import f_add from "../../images/f_add.ico"
import f_f_ico from "../../images/f_friends.ico"
import f_from_ico from "../../images/f_req_from.ico"
import f_to_ico from "../../images/f_req_to.ico"
import Throttle from "../../functools/Trottle";
import {HTTP, ADDR} from "../../address";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
const relationshipsIcons = [f_add, f_to_ico, f_from_ico, f_f_ico];

class ProfilePage extends React.Component {
    state = {posts: [], profileInfo: {}, offset: 0, error: null}
    myId = +getCookie("userId");
    postsFetchLimit = 10;
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);
    relationButtonClickedHandler = async () => {
        const [error] = await Fetcher (
            HTTP + ADDR + "/relations/update_relationship",
            {userId: this.props.match.params.id, prevRelType: this.state.profileInfo.rel},
            "POST",
            "text"
            );
        if (error === null){
            const profileInfo = {...this.state.profileInfo};
            profileInfo.rel = oldNewRelMatchDict[profileInfo.rel];
            this.setState({profileInfo});
        }else {
            this.handleError("Запрос отклонен сервером")
        }
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
            this.fetchPosts();
        }
    }, 1000);

    fetchProfileInfo = async () =>{
        let [error, profileInfo] = await Fetcher(
            HTTP + ADDR + `/profile/page_info`,
            {userId: this.props.match.params.id});
        if (error == null){
            try{
                const  x = {profileInfo}
                this.setState(x);

            }catch (e) {
                console.log(profileInfo);
            }

        }else{
            this.handleError("Невозможно загрузить информацию о профиле")
        }
    }

    fetchPosts = async () => {
        const offset = this.state.offset;
        const params = {userId: this.props.match.params.id, offset: offset, limit: this.postsFetchLimit}
        this.setState(() => ({isFetching: true}));

        const [error, response] = await Fetcher(HTTP + ADDR + "/profile/get_posts", params);
        if (error === null){
            this.setState(state => ({
                offset: offset + this.postsFetchLimit,
                posts: [...state.posts, ...response.Posts],
                Done: response.Done
            }));
        }else {
            this.handleError("Невозможно загрузить информацию посты")
        }
        this.setState({isFetching: false});
    }

    componentDidMount = () => {
        window.addEventListener('scroll', this.handleScrollThrottled, true);
        this.fetchProfileInfo();
        this.fetchPosts();
    }

    componentWillUnmount() {
        window.removeEventListener('scroll', this.handleScrollThrottled);
    }

    componentDidUpdate = async (prevProps) => {
        if (this.props.match.params.id !== prevProps.match.params.id){
            this.setState({posts: [], profileInfo: {}, offset: 0}, ()=>{
                this.fetchProfileInfo();
                this.fetchPosts();
            })
        }
    }

    addNewPost = async (text)=>{
        const [error, response] = await Fetcher(HTTP + ADDR + "/profile/add_new_object",
            {text: text},
            "POST",
            "json",
        )
        if (error === null){
            const newPost = {
                text: text,
                auth_id : this.myId,
                path: response.path,
                num_likes: 0,
                creation_time: response.creation_time,
                modification_time: null,
                first_name: getCookie("firstName"),
                last_name: getCookie("lastName"),
                num_comments: 0,
                me_liked: false
            };
            this.setState(s => ({posts: [newPost, ...s.posts], offset: s.offset + 1}));
        }else{
            this.handleError("Невозможно добавить новый пост")
        }
    }


    postNewLikeState = async (meLiked, path) => {
        const [error] = await Fetcher(
            HTTP + ADDR + "/profile/update_like",
            {meLiked: meLiked, path: path},
            "POST",
            "text"
        );
        return error === null;
    }

    deletePost = (path)=>{
        this.setState(s => ({posts: s.posts.filter(post => post.path !== path), offset: s.offset - 1}));
    }

    updatePostLike = new Debounce(async (meLiked, path)=>{
        if (await this.postNewLikeState(meLiked, path)){
            const copy = this.state.posts.map( post =>{
                if (post.path === path){
                    post.num_likes += post.me_liked ? -1 : 1;
                    post.me_liked = !post.me_liked;
                }
                return post;
            });
            this.setState({posts: copy});
        }
    }, 100)

    render() {
        return (
            <div className={"page__container"}>
                <div className={"profile__header_pics"}>
                    <img className={"default_img profile__bg_pic"}
                         src={HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(this.props.match.params.id)}/profile_bg.jpg`}
                         alt = {" "}
                    />
                    <img className={"default_img avatar profile__avatar"}
                         src={HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(this.props.match.params.id)}/profile_avatar.jpg`}
                         alt = {" "}
                    />
                </div>

                <div className={"navigator"}>
                    <div className={"navigator__item"}>
                        <Link to = {`/about_me/${this.props.match.params.id}`}>Обо мне</Link>
                    </div>

                    <div className={"navigator__item"}>
                        <Link to = {`/связи/${this.props.match.params.id}/друзья`}>
                            Мои друзья
                        </Link>
                    </div>

                    <div className={"navigator__item"}>
                        <Link to = {`/фотографии/${this.props.match.params.id}`}>
                            Мои фото
                        </Link>
                    </div>

                    <div className={"navigator__item profile__descr"}>
                        <p className={"profile_name"}>{`${this.state.profileInfo.first_name} ${this.state.profileInfo.last_name}`}</p>
                        <p className={"profile_country"}>{this.state.profileInfo.country}</p>
                    </div>

                    <div className={"navigator__item"}>
                        <Link to = {`/музыка/${this.props.match.params.id}`}>
                            Моя музыка
                        </Link>
                    </div>

                    <div className={"navigator__item"}>
                        <Link to = {`/видео/${this.props.match.params.id}`}>
                            Мои видео
                        </Link>
                    </div>

                    <div className={"navigator__item"}>
                        <div className={"navigator__functional"}>
                            <Link to = {{
                                pathname: '/сообщения/',
                                user_id: this.props.match.params.id,
                                first_name: this.state.profileInfo.first_name,
                                last_name: this.state.profileInfo.last_name,
                            }}>
                                <div className={"default_img open_conversation_img navigator__functional_item"}/>
                            </Link>
                            <button type={"button"}
                                    className={"default_img relations_with_me__button navigator__functional_item"}
                                    onClick={this.relationButtonClickedHandler}
                                    style={{
                                        backgroundImage: `url(${relationshipsIcons[this.state.profileInfo.rel]})`,
                                        display: `${this.state.profileInfo.rel === undefined ? 'none': 'block'}`
                                    }}
                            />
                        </div>
                    </div>
                </div>

                <div className={"profile_info__container"}>
                    <div className={"profile__sidebar__left"}>
                        <div className={"default_block"}>
                            <h1 className={"default_block__header"}>Кратко обо мне</h1>
                            <div className={"default_block__item"}>
                                <h2>Город проживания:</h2>
                                <p>{this.state.profileInfo.city  || "Не указано"}</p>
                            </div>
                            <div className={"default_block__item"}>
                                <h2>День рождения:</h2>
                                <p>{this.state.profileInfo.birthday  || "Не указано"}</p>
                            </div>
                            <div className={"default_block__item"}>
                                <h2>Телефон:</h2>
                                <p>{this.state.profileInfo.tel  || "Не указано"}</p>
                            </div>
                        </div>
                    </div>
                    <div className={"profile__posts"}>
                        {+this.props.match.params.id === this.myId &&
                            <AddNewMessageArea
                                buttonMessage = {"Опубликовать"}
                                placeholder = {"Что у Вас нового?"}
                                onSend = {this.addNewPost}
                            />
                        }

                        {Array.isArray(this.state.posts) && this.state.posts.map((post)=>
                            <>
                                <Post
                                    key = {post.path}
                                    updatePostLike = {this.updatePostLike}
                                    postNewLikeState = {this.postNewLikeState}
                                    deletePost = {this.deletePost}
                                    myId = {this.myId}
                                    handleError = {this.handleError}
                                    {...post}
                                />
                            </>
                            )
                        }
                    </div>


                    <div className={"profile__sidebar__right"}>
                        <div className={"default_block"}>
                            <h1 className={"default_block__header"}>Последние фото</h1>
                            <div className={"profile_gallery__container"}>
                                {
                                    Array.isArray(this.state.profileInfo.images_list)
                                    && this.state.profileInfo.images_list.map( (imageId, idx)=>
                                        <img
                                            key = {idx}
                                            className={"profile_gallery__item"}
                                            src={HTTP + ADDR + `/gallery_storage${PathFromIdGenerator(imageId)}/img.jpg`}
                                            alt = {" "}
                                        >
                                        </img>
                                    )
                                }
                            </div>
                        </div>
                    </div>
                </div>
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }

}

const Post = memo(({updatePostLike, postNewLikeState, deletePost, handleError, ...props}) => {
    const [comments, updateComments] = useState([]);
    const [postEditable, updatePostEditable] = useState(false);

    const textRef = useRef();
    const _updatePostLike = () => {
        updatePostLike(props.me_liked, props.path);
    }

    const addNewComment = async (text, path = props.path)=>{
        const [error, response] = await Fetcher(
            HTTP + ADDR + "/profile/add_new_object",
            {text: text, path: path},
            "POST"
        );
        if (error === null) {
            const newComment = {
                creation_time: response.creation_time,
                first_name: getCookie("firstName"),
                last_name: getCookie("lastName"),
                me_liked: false,
                auth_id: props.myId,
                modification_time: null,
                num_likes: 0,
                path: response.path,
                text: text,
                children: []
            }

            let copy;
            path = path.split('.').slice(1).map((el)=> el - 1);
            if (path.length === 0){
                copy = [...comments, newComment]
            }else{
                copy = [...comments];
                let pointer = copy[path[0]];
                for (let i = 1; i < path.length; i++){
                    pointer = pointer.children[path[i]];
                }
                pointer.children.push(newComment);
            }
            updateComments(copy);
        }else{
            handleError("Невозможно отправить комментарий")
        }
    }

    const updateCommentLike = new Debounce(async(meLiked, path)=>{
        if (await postNewLikeState(meLiked, path)){
            let newComments;
            const splitPath = path.split('.');
            if (splitPath.length === 2){
                newComments = comments.map(comment => {
                    if (comment.path === path){
                        comment.num_likes += comment.me_liked ? -1 : 1;
                        comment.me_liked = !comment.me_liked;
                    }
                    return comment;
                })
            }else{
                newComments = [...comments];
                let pointer = newComments.find((comment) => comment.path === splitPath.slice(0, 2).join('.'));
                for (let i = 2; i < splitPath.length; i++){
                    pointer = pointer.children.find((comment) => comment?.path === splitPath.slice(0, i+1).join('.'));
                }
                pointer.num_likes += pointer.me_liked ? -1 : 1;
                pointer.me_liked = !pointer.me_liked;
            }
            updateComments(newComments);
        }
    }, 100)

    const getComments = async () =>{
        try{
            const result = await fetch(HTTP + ADDR + `/profile/get_comments?path=${props.path}&lim=${10}`,{method: "get"});
            const json = await result.json();
            await updateComments(json);
        }catch(e)
        {
            console.log(e);
        }
    };

    const onPostDelete = async ()=>{
        const [error] = await Fetcher(
            HTTP + ADDR + "/profile/delete_object",
            null,
            "POST",
            "text",
            props.path
        )
        if (error === null){
            deletePost(props.path);
        }else{
            handleError("Невозможно удалить пост");
        }
    }

    const onPostEdit = async ()=>{
        if (postEditable){
            const [error] = await Fetcher(
                HTTP + ADDR + "/profile/update_object_text",
                null,
                "POST",
                "text",
                JSON.stringify({Path: props.path, Text:textRef.current.innerText})
                )
            if (error !== null){
                handleError("Невозможно изменить пост")
            }
        }
        updatePostEditable(!postEditable);
    }

    const deleteComment = (path) =>{
        let newComments;
        const splitPath = path.split('.');
        if (splitPath.length === 2){
            newComments = comments.filter(comment => {
                return comment.path !== path
            });

        }else{
            newComments = [...comments];
            let pointer = newComments.find((comment) => comment.path === splitPath.slice(0, 2).join('.'));
            for (let i = 2; i < splitPath.length - 1; i++){
                pointer = pointer.children.find((comment) => comment?.path === splitPath.slice(0, i+1).join('.'));
            }
            pointer.children = pointer.children.filter(comment => {
                return comment.path !== path
            });
        }
        updateComments(newComments);
    }
    return (
        <div className={"record post"}>
            <div className={"record__header"}>
                <ProfileAvatarBlock
                    src = {HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(props.auth_id)}/profile_avatar.jpg`}
                    name = {props.first_name + " " + props.last_name}
                    description = {`Добавлен: ${props.creation_time}` + (props.modification_time ? ` изменен: ${props.modification_time}`:'')}
                />
            </div>

            <div className={"record__text"} contentEditable={postEditable} ref = {textRef}>
                {props.text}
            </div>

            <div className={"record__footer"}>
                <div className={"footer_item likes"}>
                    <button type={"button"}
                            className={`footer_icon default_img ${props.me_liked? 'liked_pic':'not_liked_pic'}`}
                            onClick={_updatePostLike}
                    />
                    <p>{props.num_likes}</p>
                </div>

                <div className={"footer_item comments"}>
                    <button type={"button"} onClick={getComments} className={"footer_icon default_img comments_pic"}/>
                    <p>{props.num_comments}+</p>
                </div>

                {props.myId === props.auth_id
                 &&
                <div className={"footer_item post_comment__functional"}>
                    <button type={"button"}
                            className={"footer_icon default_img post_comment_button__edit"}
                            onClick={onPostEdit}

                    >{postEditable ? '✓' : '🖉'}</button>
                    <button type={"button"}
                            className={"footer_icon default_img post_comment_button__delete"}
                            onClick={onPostDelete}
                    >❌</button>
                </div>
                }

            </div>
            {Array.isArray(comments) && comments.map((c_props)=>
                <Comment
                    key={c_props.path}
                    updateCommentLike = {updateCommentLike}
                    addNewComment={addNewComment}
                    deleteComment = {deleteComment}
                    myId = {props.myId}
                    handleError={handleError}
                    {...c_props}
                />)
            }
            <AddNewMessageArea placeholder={"Оставить комментарий"} buttonMessage={"Комментировать"} onSend={addNewComment}/>
        </div>
    )
})

function Comment({children = [], addNewComment, updateCommentLike, deleteComment, handleError, ...props}) {
    const [showReplyArea, updateShowReplyArea] = useState(false);
    const [commentEditable, updateCommentEditable] = useState(false);
    const textRef = useRef()

    const onCommentDelete = async ()=>{
        const [error] = await Fetcher(
            HTTP + ADDR + "/profile/delete_object",
            null,
            "POST",
            "text",
            props.path
        )

        if (error === null){
            deleteComment(props.path);
        }else{
            handleError("Невозможно удалить комментарий")
        }
    }

    const _addNewComment = (text)=>{
        updateShowReplyArea(false);
        addNewComment(text, props.path);
    }

    const updateLike = () => {
        updateCommentLike(props.me_liked, props.path);
    };

    const onReplyClickHandler = () =>{
        updateShowReplyArea(true);
    }

    const onCommentEdit = async ()=>{
        if (commentEditable){
            const [error] = await Fetcher(
                HTTP + ADDR + "/profile/update_object_text",
                null,
                "POST",
                "text",
                JSON.stringify({Path: props.path, Text:textRef.current.innerText})
            )
            if (error !== null){
                handleError("Невозможно изменить комментарий")
            }
        }
        updateCommentEditable(!commentEditable);
    }

    return (
        <div className={"comment_container"} style={{paddingLeft: '10px'}}>
            <div className={"record comment"}>
                <div className={"record__header"}>
                    <ProfileAvatarBlock
                        src = {HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(props.auth_id)}/profile_avatar.jpg`}
                        name = {props.first_name + " " + props.last_name}
                        description = {`Добавлен: ${props.creation_time}` + (props.modification_time ? ` изменен: ${props.modification_time}`:'')}
                    />
                </div>

                <div className={"record__text"} ref = {textRef} contentEditable={commentEditable}>
                    {props.text}
                </div>

                <div className={"record__footer"}>
                    <div className={"footer_item likes"}>
                        <button type={"button"} onClick={updateLike} className={`footer_icon default_img ${props.me_liked? 'liked_pic':'not_liked_pic'}`}/>
                        <p>{props.num_likes}</p>
                    </div>

                    <div className={"footer_item reply"} onClick={onReplyClickHandler}>
                        <button type={"button"} className={"footer_icon default_img reply_pic"}/>
                    </div>
                    {props.myId === props.auth_id
                    &&
                    <div className={"footer_item post_comment__functional"}>
                        <button type={"button"}
                                className={"footer_icon default_img post_comment_button__edit"}
                                onClick={onCommentEdit}

                        >{commentEditable ? '✓' : '🖉'}</button>
                        <button type={"button"}
                                className={"footer_icon default_img post_comment_button__delete"}
                                onClick={onCommentDelete}
                        >❌</button>
                    </div>
                    }
                </div>
                {showReplyArea && <AddNewMessageArea buttonMessage={"Ответить"} onSend={_addNewComment}/>}
                {children.map((c_props)=> c_props && <Comment updateCommentLike = {updateCommentLike}
                                                              key={c_props.path}
                                                              addNewComment={addNewComment}
                                                              deleteComment = {deleteComment}
                                                              handleError={handleError}
                                                              myId = {props.myId}
                                                              {...c_props}
                />)}
            </div>
        </div>
    )
}


export default withRouter(ProfilePage)
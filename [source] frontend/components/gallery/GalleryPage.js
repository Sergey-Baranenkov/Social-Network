import React, {createRef, memo, useState} from "react";
import "./gallery_page.scss"
import Fetcher from "../../functools/Fetcher";
import Throttle from "../../functools/Trottle";
import "../../scss/default_blocks.scss";
import PathFromIdGenerator from "../../functools/PathFromIdGenerator";
import {withRouter} from "react-router-dom"
import {HTTP, ADDR} from "../../address";
import "../../scss/hidden_input.scss"
import "../../scss/hidden_input.scss"
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
import getCookie from "../../functools/getCookie";

class GalleryPage extends React.Component{
    addImageRef = createRef();
    fetchLimit = 20;
    pageId = + this.props.match.params.id;
    myId = + getCookie("userId");
    state = {ImagesList: [], offset: 0, isFetching: false, Done: false, error: null}
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);

    componentDidMount() {
        const fl = this.pageId === this.myId ? this.fetchLimit - 1: this.fetchLimit;
        this.fetchImages(fl);
        window.addEventListener('scroll', this.handleScrollThrottled, true);
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps !== this.props){
            const fl = this.pageId === this.myId ? this.fetchLimit - 1: this.fetchLimit;
            this.pageId = + this.props.match.params.id;
            this.setState({ImagesList: [], offset: 0, isFetching: false, Done: false, error: null}, ()=> this.fetchImages(fl));
        }
    }


    componentWillUnmount() {
        window.removeEventListener('scroll', this.handleScrollThrottled);
    }

    fetchImages = async (fetchLimit) => {
        this.setState(() => ({isFetching: true}));
        const [error, response] = await Fetcher(
            HTTP + ADDR + "/gallery/get_images",
            {limit: fetchLimit, startFrom: this.state.offset, userId : this.pageId}
        )

        if (error === null){
            this.setState(state => ({
                offset: state.offset + fetchLimit,
                ImagesList: [...state.ImagesList, ...response.ImagesList],
                Done: response.Done
            }));
        }else {
            this.handleError("Невозможно получить данные с сервера")
        }
        this.setState(() => ({isFetching: false}));
    };
    

    handleScrollThrottled = Throttle(()=>{
        if (Math.round(window.scrollY + window.innerHeight) === document.documentElement.scrollHeight
            &&
            !this.state.Done
            &&
            !this.state.isFetching
        ){
            this.fetchImages(this.fetchLimit);
        }
    }, 1000)

    postImageHandler = async ()=>{
        if (this.addImageRef.current.files.length){
            const data = new FormData();
            data.append("image", this.addImageRef.current.files[0])

            const [error, response] = await Fetcher(
                HTTP + ADDR + '/gallery/post_image',
                {},
                "post",
                "text",
                data
                );

            if (error === null){
                const imageId = +response;
                this.setState(state => ({ImagesList:
                        [
                            {image_id:
                                imageId,
                                adder_id: this.myId,
                            },
                            ...state.ImagesList
                        ],
                    offset: state.offset + 1
                }))
            }else{
                this.handleError("Невозможно добавить изображение")
            }

        }
    }
    addImageToMyGalleryHandler = async (imageId)=>{
        const address = HTTP + ADDR + "/gallery/add_to_my_gallery"
        const params = {imageId};
        const [error] = await Fetcher(
            address,
            params,
            "Post",
            "text"
        );
        if (error !== null){
            this.handleError("Невозможно добавить изображение в вашу галерею")
        }

    }
    deleteImageHandler = async (imageId)=>{
        const address = HTTP + ADDR + "/gallery/delete_image"
        const params = {imageId}

        const [error] = await Fetcher(
            address,
            params,
            "GET",
            "text"
        )

        if (error === null){
            this.setState(state => (
                        {
                            offset: state.offset - 1,
                            ImagesList: state.ImagesList.filter((image)=>image.image_id !== imageId)
                        }
                    )
            );
        }else {
            this.handleError("Невозможно удалить изображение")
        }
    }

    render() {
        return (
            <div className={"page__container"}>
                <div className={"page__header background_pic__city background_pic__city_purple"}>
                    <h1 style={{color: "white"}}>Фотографии</h1>
                    <p style={{color:  "white"}}>Здесь вы можете посмотреть фотографии</p>
                    <div className={"default_img"}/>
                </div>
                <div className={"gallery__container"}>
                    {this.pageId === this.myId &&
                            <div className={"gallery__item"}>
                                <input
                                    type={"file"}
                                    className={"hidden_input"}
                                    id = {"select_image_file__input"}
                                    accept={"image/jpeg"}
                                    ref = {this.addImageRef}
                                    onChange={this.postImageHandler}
                                />
                                <label className={"gallery__button_add"}
                                       htmlFor="select_image_file__input">
                                    +
                                </label>
                            </div>
                    }
                    {
                        this.state.ImagesList.map((image) =>
                            <PhotoBlock
                                key = {image.image_id}
                                imageId={image.image_id}
                                myId = {this.myId}
                                pageId = {this.pageId}
                                deleteImageHandler = {this.deleteImageHandler}
                                addImageHandler = {this.addImageToMyGalleryHandler}
                            />
                        )
                    }
                </div>
                {this.state.isFetching && <span>Загрузка...</span>}
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
}
const PhotoBlock = memo(({imageId, deleteImageHandler, addImageHandler, myId, pageId}) =>{
    const address = HTTP + ADDR + "/gallery_storage" + PathFromIdGenerator(imageId) + "/img.jpg";
    const [isExpanded, update] = useState(false);
    const updateSize = ()=>{
        update(!isExpanded);
    }

    const onDeleteImage = () =>{
        update(false);
        deleteImageHandler(imageId);
    }

    const onAddImage = () =>{
        addImageHandler(imageId);
    }

    return (
        <div className={"gallery__item"}>
            <img className={`gallery__photo ${isExpanded ? 'photo__expanded': null}`}
                 src={address}
                 onClick={updateSize}
                 alt = {" "}
            />
            {isExpanded
            &&
            (
                pageId === myId
                    ?
                    <button className={"gallery__button gallery_delete_photo_button"}
                            onClick={onDeleteImage}
                    >
                        Удалить
                    </button>
                    :
                    <button className={"gallery__button gallery_delete_photo_button"}
                            onClick={onAddImage}
                    >
                        Добавить
                    </button>
            )
            }
        </div>
    )
})

export default withRouter(GalleryPage)
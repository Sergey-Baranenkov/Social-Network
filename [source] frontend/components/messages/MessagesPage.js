import React, {memo, useEffect, useRef} from "react";
import "./messages_page.scss";
import {withRouter} from "react-router-dom";
import ProfileAvatarBlock from "../profile/ProfileAvatarBlock";
import getCookie from "../../functools/getCookie";
import Fetcher from "../../functools/Fetcher";
import undefined_avatar_pic from "../../images/undefined_avatar.png"
import AddNewMessageArea from "../addNewMessageArea/addNewMessageArea";
import "../../scss/default_blocks.scss"
import Throttle from "../../functools/Trottle";
import PathFromIdGenerator from "../../functools/PathFromIdGenerator";
import {WS, ADDR, HTTP} from "../../address";
import InfoPopup, {handleError, handleClose} from "../infoPopup/infoPopup";
class MessagesPage extends React.Component{
    state = {
        conversations: [],
        messengerMessagesList: [],
        messengerConversationId: undefined,
        messengerPartnerInfo : {userId: undefined, first_name: undefined, last_name: undefined, profile_avatar: undefined},
        messengerOffset: 0,
        conversationsOffset : 0,
        conversationsDone : false,
        messengerMessagesDone : false,
        areConversationsFetching : false,
        messengerAreMessagesFetching: false,
        error: null
    }
    conversationsLimit = 20;
    messagesLimit = 20;

    socket = new WebSocket(WS + ADDR + "/messenger/");
    myId = +getCookie("userId");
    handleError = handleError.bind(this);
    handleClose = handleClose.bind(this);

    openConversationFromOtherPage = ()=>{
        const params = this.props.location;
        this.openDialog(params.user_id, params.first_name, params.last_name);
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps.location.user_id !== this.props.location.user_id){ // searchbar message when on messages page
            this.openConversationFromOtherPage()
        }
    }

    componentDidMount () {
        if (this.props.location.user_id !== undefined){
            this.openConversationFromOtherPage()
        }

        this.fetchConversations()

        this.socket.onopen = () => {
            console.log("[open] Соединение установлено");
        }
        this.socket.onerror = () => {
            console.log("[error] Ошибка соединения");
        }

        this.socket.onmessage = async (event) => {
            const json = JSON.parse(event.data);
            let unknownConversationFlag = true;

            let newConversations = this.state.conversations.map( conversation => {
                if (conversation.conversation_id === json.conversation_id) {
                    unknownConversationFlag = false;
                    conversation.message_text = json.message_text;
                    return conversation;
                } else {
                    return conversation
                }
            });

            if (unknownConversationFlag === true) {
                const [error, response] = await Fetcher(
                    HTTP + ADDR + "/messenger/get_short_profile_info",
                    {conversationId: json.conversation_id}
                    );
                if (error === null){
                    newConversations = [{...response, ...json}, ...newConversations];
                    this.setState(s => ({conversationsOffset : s.conversationsOffset + 1}))
                }else{
                    this.handleError("Невозможно загрузить информацию профиля")
                }
            }
            this.setState({conversations: newConversations});
            if (this.state.messengerConversationId === 0 && this.state.messengerPartnerInfo.userId === json.message_to){
                this.setState(()=> ({messengerConversationId: json.conversation_id}));
            }

            if (json.conversation_id === this.state.messengerConversationId){
                this.setState(state => ({messengerMessagesList: [...state.messengerMessagesList, json]}))
            }
        };

    }

    handleConversationListScroll = ({target}) => {
        const difference = target.scrollTopMax - target.scrollTop;
        this._handleConversationListScroll(difference);
    }

    _handleConversationListScroll = Throttle((difference)=> {
            if (difference < 10 &&
                !this.state.conversationsDone &&
                !this.state.areConversationsFetching) {
                this.fetchConversations();
            }
        }, 1000
    )


    onMessengerTopBoundaryReached = Throttle((difference)=>{
        if (difference < 10 &&
            !this.state.messengerMessagesDone &&
            !this.state.messengerAreMessagesFetching) {
            this.fetchMessages();
        }
    }, 1000)


    componentWillUnmount() {
        this.socket.close();
    }

    fetchMessages = async ()=>{
        const offset = this.state.messengerOffset;
        const [error, response] =  await Fetcher(HTTP + ADDR + "/messenger/conversation_messages",
            {
                userId2: this.state.messengerPartnerInfo.userId,
                limit: this.messagesLimit,
                offset: offset
            }
        );

        if (error === null){
            this.setState(s => ({
                messengerOffset: offset + this.messagesLimit,
                Done: response.Done,
                messengerMessagesList: [...response.MessagesList.reverse(), ...s.messengerMessagesList],
                messengerConversationId: response.ConversationId
            }))
        }else{
            this.handleError("Невозможно загрузить сообщения с сервера")
        }
        return [error, response];
    }

    openDialog = (userId, first_name, last_name)=>{
        this.setState(
            {
                messengerPartnerInfo : {userId, first_name, last_name, avatar_ref: HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(userId)}/profile_avatar.jpg`},
                messengerMessagesList: [],
                messengerConversationId: undefined,
                Done: false,
                messengerOffset : 0
            },
            ()=>{this.fetchMessages()}
            );
    }

    fetchConversations = async (offset) =>{
        if (offset === undefined) offset = this.state.conversationsOffset;
        const [error, response] = await Fetcher(HTTP + ADDR + "/messenger/conversation_list",{limit: this.conversationsLimit, offset: offset});
        if (error === null){
            this.setState(s => ({conversations: [...s.conversations, ...response.Conversations], conversationsOffset: offset + this.conversationsLimit, conversationsDone: response.Done}));
        }else{
            this.handleError("Невозможно загрузить список диалогов")
        }
    }

    render() {
        return (
            <div className={"message_page"}>
                <div className={"conversations_list"} onScroll={this.handleConversationListScroll}>
                    {
                        this.state.conversations.map((conversation)=>(
                            <Conversation key = {conversation.conversation_id}
                                          onChoose = {this.openDialog}
                                          chosenId = {this.state.messengerPartnerInfo.userId}
                                          {...conversation}
                            />
                        ))
                    }
                </div>

                <Messenger myId = {this.myId}
                           messagesList = {this.state.messengerMessagesList}
                           partnerInfo = {this.state.messengerPartnerInfo}
                           onMessengerTopBoundaryReached = {this.onMessengerTopBoundaryReached}
                           socket = {this.socket}
                />
                {this.state.error && <InfoPopup text = {this.state.error} handleClose = {this.handleClose}/>}
            </div>
        )
    }
}

const Conversation = memo(({onChoose, chosenId, partner_id, first_name, last_name, message_text})=>{
    const handleChoose = () => onChoose(partner_id, first_name, last_name);
    return (
      <div className={"rel__container"}
           style={{backgroundColor: `${chosenId === partner_id ? "lightgoldenrodyellow" : "white"}`, cursor: 'pointer'}}
           tabIndex={1}
           onClick={handleChoose}
      >
          <img className={"default_img rel_user__avatar"}
               src={HTTP + ADDR + `/profile_bgs${PathFromIdGenerator(partner_id)}/profile_avatar.jpg`}
               alt = {" "}
          />
          <div className={"rel_user__vertical"}>
              <span className={"rel_user__fullname"}>{`${first_name} ${last_name}`}</span>
              <span className={"rel_user__last_message"}>{`${message_text.length > 24 ? message_text.slice(0,24) + "..." : message_text}`}</span>
          </div>
      </div>
    )
})

const Messenger = memo (({messagesList, myId, partnerInfo, onMessengerTopBoundaryReached, socket}) => {

    const handleMessagesListScroll = ({target}) => {
        const difference = target.scrollTop;
        onMessengerTopBoundaryReached(difference);
    }

    const messagesBoxRef = useRef(null);
    useEffect(() => {
        /*
        * 200 is the difference between new bottom and old scroll
        * when user is not scrolling old messages, but ready for new ones - scroll to bottom
        */
        if (messagesBoxRef.current.scrollTopMax - messagesBoxRef.current.scrollTop < 200){
            messagesBoxRef.current.scrollTop = messagesBoxRef.current.scrollTopMax;
        }
    });

    const pushMessage = (text) => {
        socket.send(JSON.stringify({messageTo: +partnerInfo.userId, messageText: text}))
    };

    return(
        <div className={"messenger"} style={{pointerEvents: `${partnerInfo.userId === undefined?'none':'auto'}`}}>
            <div className={"messenger__header"}>
                <ProfileAvatarBlock
                    src = {partnerInfo.avatar_ref || undefined_avatar_pic}
                    name = {`${partnerInfo.first_name || ""} ${partnerInfo.last_name || ""}`}
                    nameColor = "white"
                    description = ""
                    descriptionColor = "cadetblue"
                />
            </div>

            <div className={"messenger__messages-box"} ref = {messagesBoxRef} onScroll={handleMessagesListScroll}>
                {
                    messagesList.map((message) =>
                        <Message text = {message.message_text}
                                 myId = {myId}
                                 senderId = {message.message_from}
                                 key = {message.message_id}
                        />
                    )
                }
            </div>

            <div className={"messenger__send_new_message"}>
                {partnerInfo.userId===undefined && <div className={"tumbleweed"}/>}
                <AddNewMessageArea placeholder = {"Отправить сообщение"}
                                   buttonMessage = {"Отправить"}
                                   onSend = {pushMessage}
                />
            </div>
        </div>
    )
})

const Message = memo(({text = "", myId, senderId}) => {
    return (
            <div className={`message ${senderId === myId ? 'message__from_me' : 'message__from_friend'}`}>
                {text}
            </div>
        )
});

export default withRouter(MessagesPage);
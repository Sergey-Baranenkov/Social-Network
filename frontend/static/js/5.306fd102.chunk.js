(this["webpackJsonpmy-app"]=this["webpackJsonpmy-app"]||[]).push([[5],{101:function(e,t,s){"use strict";s.r(t);var n=s(33),a=s(16),r=s.n(a),o=s(22),c=s(17),i=s(20),l=s(11),m=s(12),u=s(14),f=s(13),g=s(0),d=s.n(g),v=(s(81),s(18)),_=s(36),p=s(23),h=s(26),b=s(82),O=s.n(b),j=s(76),I=(s(55),s(31)),E=s(21),y=s(4),M=function(e){Object(u.a)(s,e);var t=Object(f.a)(s);function s(){var e;Object(l.a)(this,s);for(var n=arguments.length,a=new Array(n),m=0;m<n;m++)a[m]=arguments[m];return(e=t.call.apply(t,[this].concat(a))).state={conversations:[],messengerMessagesList:[],messengerConversationId:void 0,messengerPartnerInfo:{userId:void 0,first_name:void 0,last_name:void 0,profile_avatar:void 0},messengerOffset:0,conversationsOffset:0,conversationsDone:!1,messengerMessagesDone:!1,areConversationsFetching:!1,messengerAreMessagesFetching:!1},e.conversationsLimit=20,e.messagesLimit=20,e.socket=new WebSocket(y.c+y.a+"/messenger/"),e.myId=+Object(p.a)("userId"),e.handleConversationListScroll=function(t){var s=t.target,n=s.scrollTopMax-s.scrollTop;e._handleConversationListScroll(n)},e._handleConversationListScroll=Object(I.a)((function(t){t<10&&!e.state.conversationsDone&&!e.state.areConversationsFetching&&e.fetchConversations()}),1e3),e.onMessengerTopBoundaryReached=Object(I.a)((function(t){t<10&&!e.state.messengerMessagesDone&&!e.state.messengerAreMessagesFetching&&e.fetchMessages()}),1e3),e.fetchMessages=Object(i.a)(r.a.mark((function t(){var s,n,a,i,l;return r.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return s=e.state.messengerOffset,t.next=3,Object(h.a)(y.b+y.a+"/messenger/conversation_messages",{userId2:e.state.messengerPartnerInfo.userId,limit:e.messagesLimit,offset:s});case 3:return n=t.sent,a=Object(c.a)(n,2),i=a[0],l=a[1],null===i&&e.setState((function(t){return{messengerOffset:s+e.messagesLimit,Done:l.Done,messengerMessagesList:[].concat(Object(o.a)(l.MessagesList.reverse()),Object(o.a)(t.messengerMessagesList)),messengerConversationId:l.ConversationId}})),t.abrupt("return",[i,l]);case 9:case"end":return t.stop()}}),t)}))),e.openDialog=function(t,s,n){e.setState({messengerPartnerInfo:{userId:t,first_name:s,last_name:n,avatar_ref:y.b+y.a+"/profile_bgs".concat(Object(E.a)(t),"/profile_avatar.jpg")},messengerMessagesList:[],messengerConversationId:null,Done:!1,messengerOffset:0},(function(){e.fetchMessages()}))},e.fetchConversations=function(){var t=Object(i.a)(r.a.mark((function t(s){var n,a,i,l;return r.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return void 0===s&&(s=e.state.conversationsOffset),t.next=3,Object(h.a)(y.b+y.a+"/messenger/conversation_list",{limit:e.conversationsLimit,offset:s});case 3:n=t.sent,a=Object(c.a)(n,2),i=a[0],l=a[1],null===i&&e.setState((function(t){return{conversations:[].concat(Object(o.a)(t.conversations),Object(o.a)(l.Conversations)),conversationsOffset:s+e.conversationsLimit,conversationsDone:l.Done}}));case 8:case"end":return t.stop()}}),t)})));return function(e){return t.apply(this,arguments)}}(),e}return Object(m.a)(s,[{key:"componentDidMount",value:function(){var e=this;if(void 0!==this.props.location.user_id){var t=this.props.location;this.openDialog(t.user_id,t.first_name,t.last_name)}this.fetchConversations(),this.socket.onopen=function(){console.log("[open] \u0421\u043e\u0435\u0434\u0438\u043d\u0435\u043d\u0438\u0435 \u0443\u0441\u0442\u0430\u043d\u043e\u0432\u043b\u0435\u043d\u043e")},this.socket.onerror=function(){console.log("[error] \u041e\u0448\u0438\u0431\u043a\u0430 \u0441\u043e\u0435\u0434\u0438\u043d\u0435\u043d\u0438\u044f")},this.socket.onmessage=function(){var t=Object(i.a)(r.a.mark((function t(s){var a,i,l,m,u,f,g;return r.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:if(a=JSON.parse(s.data),i=!0,l=e.state.conversations.map((function(e){return e.conversation_id===a.conversation_id?(i=!1,e.message_text=a.message_text,e):e})),!0!==i){t.next=11;break}return t.next=6,Object(h.a)(y.b+y.a+"/messenger/get_short_profile_info",{conversationId:a.conversation_id});case 6:m=t.sent,u=Object(c.a)(m,2),f=u[0],g=u[1],null===f?(l=[Object(n.a)(Object(n.a)({},g),a)].concat(Object(o.a)(l)),e.setState((function(e){return{conversationsOffset:e.conversationsOffset+1}}))):console.error("unhandled error");case 11:e.setState({conversations:l}),a.conversation_id===e.state.messengerConversationId&&e.setState((function(e){return{messengerMessagesList:[].concat(Object(o.a)(e.messengerMessagesList),[a])}}));case 13:case"end":return t.stop()}}),t)})));return function(e){return t.apply(this,arguments)}}()}},{key:"componentWillUnmount",value:function(){this.socket.close()}},{key:"render",value:function(){var e=this;return d.a.createElement("div",{className:"message_page"},d.a.createElement("div",{className:"conversations_list",onScroll:this.handleConversationListScroll},this.state.conversations.map((function(t){return d.a.createElement(C,Object.assign({key:t.conversation_id,onChoose:e.openDialog,chosenId:e.state.messengerPartnerInfo.userId},t))}))),d.a.createElement(x,{myId:this.myId,messagesList:this.state.messengerMessagesList,partnerInfo:this.state.messengerPartnerInfo,onMessengerTopBoundaryReached:this.onMessengerTopBoundaryReached}))}}]),s}(d.a.Component),C=Object(g.memo)((function(e){var t=e.onChoose,s=e.chosenId,n=e.partner_id,a=e.first_name,r=e.last_name,o=e.message_text;return d.a.createElement("div",{className:"rel__container",style:{backgroundColor:"".concat(s===n?"lightgoldenrodyellow":"white"),cursor:"pointer"},tabIndex:1,onClick:function(){return t(n,a,r)}},d.a.createElement("img",{className:"default_img rel_user__avatar",src:y.b+y.a+"/profile_bgs".concat(Object(E.a)(n),"/profile_avatar.jpg"),alt:""}),d.a.createElement("div",{className:"rel_user__vertical"},d.a.createElement("span",{className:"rel_user__fullname"},"".concat(a," ").concat(r)),d.a.createElement("span",{className:"rel_user__last_message"},"".concat(o.length>24?o.slice(0,24)+"...":o))))})),x=Object(g.memo)((function(e){var t=e.messagesList,s=e.myId,n=e.partnerInfo,a=e.onMessengerTopBoundaryReached,r=Object(g.useRef)(null);Object(g.useEffect)((function(){r.current.scrollTopMax-r.current.scrollTop<200&&(r.current.scrollTop=r.current.scrollTopMax)}));return d.a.createElement("div",{className:"messenger",style:{pointerEvents:"".concat(void 0===n.userId?"none":"auto")}},d.a.createElement("div",{className:"messenger__header"},d.a.createElement(_.a,{src:n.avatar_ref||O.a,name:"".concat(n.first_name||""," ").concat(n.last_name||""),nameColor:"white",description:"",descriptionColor:"cadetblue"})),d.a.createElement("div",{className:"messenger__messages-box",ref:r,onScroll:function(e){var t=e.target.scrollTop;a(t)}},t.map((function(e){return d.a.createElement(k,{text:e.message_text,myId:s,senderId:e.message_from,key:e.message_id})}))),d.a.createElement("div",{className:"messenger__send_new_message"},void 0===n.userId&&d.a.createElement("div",{className:"tumbleweed"}),d.a.createElement(j.a,{placeholder:"\u041e\u0442\u043f\u0440\u0430\u0432\u0438\u0442\u044c \u0441\u043e\u043e\u0431\u0449\u0435\u043d\u0438\u0435",buttonMessage:"\u041e\u0442\u043f\u0440\u0430\u0432\u0438\u0442\u044c",onSend:function(e){Object(h.a)(y.b+y.a+"/messenger/push_message",{messageTo:n.userId,messageText:e},"POST","text")}})))})),k=Object(g.memo)((function(e){var t=e.text,s=void 0===t?"":t,n=e.myId,a=e.senderId;return d.a.createElement("div",{className:"message ".concat(a===n?"message__from_me":"message__from_friend")},s)}));t.default=Object(v.i)(M)},76:function(e,t,s){"use strict";var n=s(0),a=s.n(n),r=(s(77),s(55),Object(n.memo)((function(e){var t=e.placeholder,s=e.buttonMessage,r=e.onSend,o=Object(n.useRef)(),c=function(e){var t=e.target,s=o.current.value,n=o.current.selectionStart;o.current.value=s.slice(0,n)+t.innerText+s.slice(n,s.length),o.current.focus(),o.current.setSelectionRange(n+2,n+2)};return a.a.createElement("div",{className:"send_new_message__container"},a.a.createElement("textarea",{className:"default_search_input",placeholder:t,ref:o}),a.a.createElement("div",{className:"send_new_message__functional_area"},a.a.createElement("div",{className:"smile_area"},["\ud83d\ude10","\ud83d\ude2b","\ud83d\ude0e","\ud83d\ude02","\ud83d\ude21","\ud83d\ude2d","\ud83d\ude00","\ud83d\ude17","\ud83d\ude32","\ud83d\ude2c"].map((function(e){return a.a.createElement("span",{className:"smile",key:e,role:"img",onClick:c},e)}))),a.a.createElement("button",{className:"add_record__button",onClick:function(){r(o.current.value),o.current.value=""}},s)))})));t.a=r},77:function(e,t,s){},81:function(e,t,s){},82:function(e,t,s){e.exports=s.p+"static/media/undefined_avatar.e630516b.png"}}]);
//# sourceMappingURL=5.306fd102.chunk.js.map
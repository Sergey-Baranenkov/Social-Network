(this["webpackJsonpmy-app"]=this["webpackJsonpmy-app"]||[]).push([[3],{111:function(e,t,a){"use strict";a.r(t);var n=a(37),r=a(25),c=a(5),o=a.n(c),i=a(27),s=a(15),l=a(19),m=a(9),u=a(10),p=a(17),f=a(12),d=a(11),_=a(0),h=a.n(_),b=a(39),v=a(6),E=a(21),j=a(80),k=a(22),O=(a(94),a(57),a(58),a(26)),g=a(38),N=a(23),y={0:1,1:0,2:3,3:2},x=a(95),w=a.n(x),I=a(96),S=a.n(I),P=a(97),C=a.n(P),L=a(98),T=a.n(L),A=a(36),F=a(3),M=a(20),D=[w.a,T.a,C.a,S.a],R=function(e){Object(f.a)(a,e);var t=Object(d.a)(a);function a(){var e;Object(m.a)(this,a);for(var n=arguments.length,c=new Array(n),u=0;u<n;u++)c[u]=arguments[u];return(e=t.call.apply(t,[this].concat(c))).state={posts:[],profileInfo:{},offset:0,error:null},e.myId=+Object(N.a)("userId"),e.postsFetchLimit=10,e.handleError=M.c.bind(Object(p.a)(e)),e.handleClose=M.b.bind(Object(p.a)(e)),e.relationButtonClickedHandler=Object(l.a)(o.a.mark((function t(){var a,n,r;return o.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,Object(k.a)(F.b+F.a+"/relations/update_relationship",{userId:e.props.match.params.id,prevRelType:e.state.profileInfo.rel},"POST","text");case 2:a=t.sent,n=Object(s.a)(a,1),null===n[0]?((r=Object(i.a)({},e.state.profileInfo)).rel=y[r.rel],e.setState({profileInfo:r})):e.handleError("\u0417\u0430\u043f\u0440\u043e\u0441 \u043e\u0442\u043a\u043b\u043e\u043d\u0435\u043d \u0441\u0435\u0440\u0432\u0435\u0440\u043e\u043c");case 6:case"end":return t.stop()}}),t)}))),e.handleScrollThrottled=Object(A.a)((function(){Math.abs(window.scrollY+window.innerHeight-document.documentElement.scrollHeight)<10&&!e.state.Done&&!e.state.popupOpen&&!e.state.isFetching&&e.fetchPosts()}),1e3),e.fetchProfileInfo=Object(l.a)(o.a.mark((function t(){var a,n,r,c,i;return o.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,Object(k.a)(F.b+F.a+"/profile/page_info",{userId:e.props.match.params.id});case 2:if(a=t.sent,n=Object(s.a)(a,2),r=n[0],c=n[1],null==r)try{i={profileInfo:c},e.setState(i)}catch(o){console.log(c)}else e.handleError("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0437\u0430\u0433\u0440\u0443\u0437\u0438\u0442\u044c \u0438\u043d\u0444\u043e\u0440\u043c\u0430\u0446\u0438\u044e \u043e \u043f\u0440\u043e\u0444\u0438\u043b\u0435");case 7:case"end":return t.stop()}}),t)}))),e.fetchPosts=Object(l.a)(o.a.mark((function t(){var a,n,c,i,l,m;return o.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return a=e.state.offset,n={userId:e.props.match.params.id,offset:a,limit:e.postsFetchLimit},e.setState((function(){return{isFetching:!0}})),t.next=5,Object(k.a)(F.b+F.a+"/profile/get_posts",n);case 5:c=t.sent,i=Object(s.a)(c,2),l=i[0],m=i[1],null===l?e.setState((function(t){return{offset:a+e.postsFetchLimit,posts:[].concat(Object(r.a)(t.posts),Object(r.a)(m.Posts)),Done:m.Done}})):e.handleError("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0437\u0430\u0433\u0440\u0443\u0437\u0438\u0442\u044c \u0438\u043d\u0444\u043e\u0440\u043c\u0430\u0446\u0438\u044e \u043f\u043e\u0441\u0442\u044b"),e.setState({isFetching:!1});case 11:case"end":return t.stop()}}),t)}))),e.componentDidMount=function(){window.addEventListener("scroll",e.handleScrollThrottled,!0),e.fetchProfileInfo(),e.fetchPosts()},e.componentDidUpdate=function(){var t=Object(l.a)(o.a.mark((function t(a){return o.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:e.props.match.params.id!==a.match.params.id&&e.setState({posts:[],profileInfo:{},offset:0},(function(){e.fetchProfileInfo(),e.fetchPosts()}));case 1:case"end":return t.stop()}}),t)})));return function(e){return t.apply(this,arguments)}}(),e.addNewPost=function(){var t=Object(l.a)(o.a.mark((function t(a){var n,c,i,l,m;return o.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,Object(k.a)(F.b+F.a+"/profile/add_new_object",{text:a},"POST","json");case 2:n=t.sent,c=Object(s.a)(n,2),i=c[0],l=c[1],null===i?(m={text:a,auth_id:e.myId,path:l.path,num_likes:0,creation_time:l.creation_time,modification_time:null,first_name:Object(N.a)("firstName"),last_name:Object(N.a)("lastName"),num_comments:0,me_liked:!1},e.setState((function(e){return{posts:[m].concat(Object(r.a)(e.posts)),offset:e.offset+1}}))):e.handleError("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0434\u043e\u0431\u0430\u0432\u0438\u0442\u044c \u043d\u043e\u0432\u044b\u0439 \u043f\u043e\u0441\u0442");case 7:case"end":return t.stop()}}),t)})));return function(e){return t.apply(this,arguments)}}(),e.postNewLikeState=function(){var e=Object(l.a)(o.a.mark((function e(t,a){var n,r,c;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,Object(k.a)(F.b+F.a+"/profile/update_like",{meLiked:t,path:a},"POST","text");case 2:return n=e.sent,r=Object(s.a)(n,1),c=r[0],e.abrupt("return",null===c);case 6:case"end":return e.stop()}}),e)})));return function(t,a){return e.apply(this,arguments)}}(),e.deletePost=function(t){e.setState((function(e){return{posts:e.posts.filter((function(e){return e.path!==t})),offset:e.offset-1}}))},e.updatePostLike=new g.a(function(){var t=Object(l.a)(o.a.mark((function t(a,n){var r;return o.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,e.postNewLikeState(a,n);case 2:if(!t.sent){t.next=5;break}r=e.state.posts.map((function(e){return e.path===n&&(e.num_likes+=e.me_liked?-1:1,e.me_liked=!e.me_liked),e})),e.setState({posts:r});case 5:case"end":return t.stop()}}),t)})));return function(e,a){return t.apply(this,arguments)}}(),100),e}return Object(u.a)(a,[{key:"componentWillUnmount",value:function(){window.removeEventListener("scroll",this.handleScrollThrottled)}},{key:"render",value:function(){var e=this;return h.a.createElement("div",{className:"page__container"},h.a.createElement("div",{className:"profile__header_pics"},h.a.createElement("img",{className:"default_img profile__bg_pic",src:F.b+F.a+"/profile_bgs".concat(Object(O.a)(this.props.match.params.id),"/profile_bg.jpg"),alt:" "}),h.a.createElement("img",{className:"default_img avatar profile__avatar",src:F.b+F.a+"/profile_bgs".concat(Object(O.a)(this.props.match.params.id),"/profile_avatar.jpg"),alt:" "})),h.a.createElement("div",{className:"navigator"},h.a.createElement("div",{className:"navigator__item"},h.a.createElement(v.c,{to:"/about_me/".concat(this.props.match.params.id)},"\u041e\u0431\u043e \u043c\u043d\u0435")),h.a.createElement("div",{className:"navigator__item"},h.a.createElement(v.c,{to:"/\u0441\u0432\u044f\u0437\u0438/".concat(this.props.match.params.id,"/\u0434\u0440\u0443\u0437\u044c\u044f")},"\u041c\u043e\u0438 \u0434\u0440\u0443\u0437\u044c\u044f")),h.a.createElement("div",{className:"navigator__item"},h.a.createElement(v.c,{to:"/\u0444\u043e\u0442\u043e\u0433\u0440\u0430\u0444\u0438\u0438/".concat(this.props.match.params.id)},"\u041c\u043e\u0438 \u0444\u043e\u0442\u043e")),h.a.createElement("div",{className:"navigator__item profile__descr"},h.a.createElement("p",{className:"profile_name"},"".concat(this.state.profileInfo.first_name," ").concat(this.state.profileInfo.last_name)),h.a.createElement("p",{className:"profile_country"},this.state.profileInfo.country)),h.a.createElement("div",{className:"navigator__item"},h.a.createElement(v.c,{to:"/\u043c\u0443\u0437\u044b\u043a\u0430/".concat(this.props.match.params.id)},"\u041c\u043e\u044f \u043c\u0443\u0437\u044b\u043a\u0430")),h.a.createElement("div",{className:"navigator__item"},h.a.createElement(v.c,{to:"/\u0432\u0438\u0434\u0435\u043e/".concat(this.props.match.params.id)},"\u041c\u043e\u0438 \u0432\u0438\u0434\u0435\u043e")),h.a.createElement("div",{className:"navigator__item"},h.a.createElement("div",{className:"navigator__functional"},h.a.createElement(v.c,{to:{pathname:"/\u0441\u043e\u043e\u0431\u0449\u0435\u043d\u0438\u044f/",user_id:this.props.match.params.id,first_name:this.state.profileInfo.first_name,last_name:this.state.profileInfo.last_name}},h.a.createElement("div",{className:"default_img open_conversation_img navigator__functional_item"})),h.a.createElement("button",{type:"button",className:"default_img relations_with_me__button navigator__functional_item",onClick:this.relationButtonClickedHandler,style:{backgroundImage:"url(".concat(D[this.state.profileInfo.rel],")"),display:"".concat(void 0===this.state.profileInfo.rel?"none":"block")}})))),h.a.createElement("div",{className:"profile_info__container"},h.a.createElement("div",{className:"profile__sidebar__left"},h.a.createElement("div",{className:"default_block"},h.a.createElement("h1",{className:"default_block__header"},"\u041a\u0440\u0430\u0442\u043a\u043e \u043e\u0431\u043e \u043c\u043d\u0435"),h.a.createElement("div",{className:"default_block__item"},h.a.createElement("h2",null,"\u0413\u043e\u0440\u043e\u0434 \u043f\u0440\u043e\u0436\u0438\u0432\u0430\u043d\u0438\u044f:"),h.a.createElement("p",null,this.state.profileInfo.city||"\u041d\u0435 \u0443\u043a\u0430\u0437\u0430\u043d\u043e")),h.a.createElement("div",{className:"default_block__item"},h.a.createElement("h2",null,"\u0414\u0435\u043d\u044c \u0440\u043e\u0436\u0434\u0435\u043d\u0438\u044f:"),h.a.createElement("p",null,this.state.profileInfo.birthday||"\u041d\u0435 \u0443\u043a\u0430\u0437\u0430\u043d\u043e")),h.a.createElement("div",{className:"default_block__item"},h.a.createElement("h2",null,"\u0422\u0435\u043b\u0435\u0444\u043e\u043d:"),h.a.createElement("p",null,this.state.profileInfo.tel||"\u041d\u0435 \u0443\u043a\u0430\u0437\u0430\u043d\u043e")))),h.a.createElement("div",{className:"profile__posts"},+this.props.match.params.id===this.myId&&h.a.createElement(j.a,{buttonMessage:"\u041e\u043f\u0443\u0431\u043b\u0438\u043a\u043e\u0432\u0430\u0442\u044c",placeholder:"\u0427\u0442\u043e \u0443 \u0412\u0430\u0441 \u043d\u043e\u0432\u043e\u0433\u043e?",onSend:this.addNewPost}),Array.isArray(this.state.posts)&&this.state.posts.map((function(t){return h.a.createElement(h.a.Fragment,null,h.a.createElement(H,Object.assign({key:t.path,updatePostLike:e.updatePostLike,postNewLikeState:e.postNewLikeState,deletePost:e.deletePost,myId:e.myId,handleError:e.handleError},t)))}))),h.a.createElement("div",{className:"profile__sidebar__right"},h.a.createElement("div",{className:"default_block"},h.a.createElement("h1",{className:"default_block__header"},"\u041f\u043e\u0441\u043b\u0435\u0434\u043d\u0438\u0435 \u0444\u043e\u0442\u043e"),h.a.createElement("div",{className:"profile_gallery__container"},Array.isArray(this.state.profileInfo.images_list)&&this.state.profileInfo.images_list.map((function(e,t){return h.a.createElement("img",{key:t,className:"profile_gallery__item",src:F.b+F.a+"/gallery_storage".concat(Object(O.a)(e),"/img.jpg"),alt:" "})})))))),this.state.error&&h.a.createElement(M.a,{text:this.state.error,handleClose:this.handleClose}))}}]),a}(h.a.Component),H=Object(_.memo)((function(e){var t=e.updatePostLike,a=e.postNewLikeState,c=e.deletePost,i=e.handleError,m=Object(n.a)(e,["updatePostLike","postNewLikeState","deletePost","handleError"]),u=Object(_.useState)([]),p=Object(s.a)(u,2),f=p[0],d=p[1],v=Object(_.useState)(!1),E=Object(s.a)(v,2),y=E[0],x=E[1],w=Object(_.useRef)(),I=function(){var e=Object(l.a)(o.a.mark((function e(t){var a,n,c,l,u,p,_,h,b,v=arguments;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return a=v.length>1&&void 0!==v[1]?v[1]:m.path,e.next=3,Object(k.a)(F.b+F.a+"/profile/add_new_object",{text:t,path:a},"POST");case 3:if(n=e.sent,c=Object(s.a)(n,2),l=c[0],u=c[1],null===l){if(p={creation_time:u.creation_time,first_name:Object(N.a)("firstName"),last_name:Object(N.a)("lastName"),me_liked:!1,auth_id:m.myId,modification_time:null,num_likes:0,path:u.path,text:t,children:[]},0===(a=a.split(".").slice(1).map((function(e){return e-1}))).length)_=[].concat(Object(r.a)(f),[p]);else{for(_=Object(r.a)(f),h=_[a[0]],b=1;b<a.length;b++)h=h.children[a[b]];h.children.push(p)}d(_)}else i("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u043e\u0442\u043f\u0440\u0430\u0432\u0438\u0442\u044c \u043a\u043e\u043c\u043c\u0435\u043d\u0442\u0430\u0440\u0438\u0439");case 8:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}(),S=new g.a(function(){var e=Object(l.a)(o.a.mark((function e(t,n){return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,a(t,n);case 2:if(!e.sent){e.next=4;break}!function(){var e,t=n.split(".");if(2===t.length)e=f.map((function(e){return e.path===n&&(e.num_likes+=e.me_liked?-1:1,e.me_liked=!e.me_liked),e}));else{for(var a=(e=Object(r.a)(f)).find((function(e){return e.path===t.slice(0,2).join(".")})),c=function(e){a=a.children.find((function(a){return(null===a||void 0===a?void 0:a.path)===t.slice(0,e+1).join(".")}))},o=2;o<t.length;o++)c(o);a.num_likes+=a.me_liked?-1:1,a.me_liked=!a.me_liked}d(e)}();case 4:case"end":return e.stop()}}),e)})));return function(t,a){return e.apply(this,arguments)}}(),100),P=function(){var e=Object(l.a)(o.a.mark((function e(){var t,a;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,fetch(F.b+F.a+"/profile/get_comments?path=".concat(m.path,"&lim=",10),{method:"get"});case 3:return t=e.sent,e.next=6,t.json();case 6:return a=e.sent,e.next=9,d(a);case 9:e.next=14;break;case 11:e.prev=11,e.t0=e.catch(0),console.log(e.t0);case 14:case"end":return e.stop()}}),e,null,[[0,11]])})));return function(){return e.apply(this,arguments)}}(),C=function(){var e=Object(l.a)(o.a.mark((function e(){var t,a;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,Object(k.a)(F.b+F.a+"/profile/delete_object",null,"POST","text",m.path);case 2:t=e.sent,a=Object(s.a)(t,1),null===a[0]?c(m.path):i("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0443\u0434\u0430\u043b\u0438\u0442\u044c \u043f\u043e\u0441\u0442");case 6:case"end":return e.stop()}}),e)})));return function(){return e.apply(this,arguments)}}(),L=function(){var e=Object(l.a)(o.a.mark((function e(){var t,a;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!y){e.next=7;break}return e.next=3,Object(k.a)(F.b+F.a+"/profile/update_object_text",null,"POST","text",JSON.stringify({Path:m.path,Text:w.current.innerText}));case 3:t=e.sent,a=Object(s.a)(t,1),null!==a[0]&&i("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0438\u0437\u043c\u0435\u043d\u0438\u0442\u044c \u043f\u043e\u0441\u0442");case 7:x(!y);case 8:case"end":return e.stop()}}),e)})));return function(){return e.apply(this,arguments)}}(),T=function(e){var t,a=e.split(".");if(2===a.length)t=f.filter((function(t){return t.path!==e}));else{for(var n=(t=Object(r.a)(f)).find((function(e){return e.path===a.slice(0,2).join(".")})),c=function(e){n=n.children.find((function(t){return(null===t||void 0===t?void 0:t.path)===a.slice(0,e+1).join(".")}))},o=2;o<a.length-1;o++)c(o);n.children=n.children.filter((function(t){return t.path!==e}))}d(t)};return h.a.createElement("div",{className:"record post"},h.a.createElement("div",{className:"record__header"},h.a.createElement(b.a,{src:F.b+F.a+"/profile_bgs".concat(Object(O.a)(m.auth_id),"/profile_avatar.jpg"),name:m.first_name+" "+m.last_name,description:"\u0414\u043e\u0431\u0430\u0432\u043b\u0435\u043d: ".concat(m.creation_time)+(m.modification_time?" \u0438\u0437\u043c\u0435\u043d\u0435\u043d: ".concat(m.modification_time):"")})),h.a.createElement("div",{className:"record__text",contentEditable:y,ref:w},m.text),h.a.createElement("div",{className:"record__footer"},h.a.createElement("div",{className:"footer_item likes"},h.a.createElement("button",{type:"button",className:"footer_icon default_img ".concat(m.me_liked?"liked_pic":"not_liked_pic"),onClick:function(){t(m.me_liked,m.path)}}),h.a.createElement("p",null,m.num_likes)),h.a.createElement("div",{className:"footer_item comments"},h.a.createElement("button",{type:"button",onClick:P,className:"footer_icon default_img comments_pic"}),h.a.createElement("p",null,m.num_comments,"+")),m.myId===m.auth_id&&h.a.createElement("div",{className:"footer_item post_comment__functional"},h.a.createElement("button",{type:"button",className:"footer_icon default_img post_comment_button__edit",onClick:L},y?"\u2713":"\ud83d\udd89"),h.a.createElement("button",{type:"button",className:"footer_icon default_img post_comment_button__delete",onClick:C},"\u274c"))),Array.isArray(f)&&f.map((function(e){return h.a.createElement(J,Object.assign({key:e.path,updateCommentLike:S,addNewComment:I,deleteComment:T,myId:m.myId,handleError:i},e))})),h.a.createElement(j.a,{placeholder:"\u041e\u0441\u0442\u0430\u0432\u0438\u0442\u044c \u043a\u043e\u043c\u043c\u0435\u043d\u0442\u0430\u0440\u0438\u0439",buttonMessage:"\u041a\u043e\u043c\u043c\u0435\u043d\u0442\u0438\u0440\u043e\u0432\u0430\u0442\u044c",onSend:I}))}));function J(e){var t=e.children,a=void 0===t?[]:t,r=e.addNewComment,c=e.updateCommentLike,i=e.deleteComment,m=e.handleError,u=Object(n.a)(e,["children","addNewComment","updateCommentLike","deleteComment","handleError"]),p=Object(_.useState)(!1),f=Object(s.a)(p,2),d=f[0],v=f[1],E=Object(_.useState)(!1),g=Object(s.a)(E,2),N=g[0],y=g[1],x=Object(_.useRef)(),w=function(){var e=Object(l.a)(o.a.mark((function e(){var t,a;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,Object(k.a)(F.b+F.a+"/profile/delete_object",null,"POST","text",u.path);case 2:t=e.sent,a=Object(s.a)(t,1),null===a[0]?i(u.path):m("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0443\u0434\u0430\u043b\u0438\u0442\u044c \u043a\u043e\u043c\u043c\u0435\u043d\u0442\u0430\u0440\u0438\u0439");case 6:case"end":return e.stop()}}),e)})));return function(){return e.apply(this,arguments)}}(),I=function(){var e=Object(l.a)(o.a.mark((function e(){var t,a;return o.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!N){e.next=7;break}return e.next=3,Object(k.a)(F.b+F.a+"/profile/update_object_text",null,"POST","text",JSON.stringify({Path:u.path,Text:x.current.innerText}));case 3:t=e.sent,a=Object(s.a)(t,1),null!==a[0]&&m("\u041d\u0435\u0432\u043e\u0437\u043c\u043e\u0436\u043d\u043e \u0438\u0437\u043c\u0435\u043d\u0438\u0442\u044c \u043a\u043e\u043c\u043c\u0435\u043d\u0442\u0430\u0440\u0438\u0439");case 7:y(!N);case 8:case"end":return e.stop()}}),e)})));return function(){return e.apply(this,arguments)}}();return h.a.createElement("div",{className:"comment_container",style:{paddingLeft:"10px"}},h.a.createElement("div",{className:"record comment"},h.a.createElement("div",{className:"record__header"},h.a.createElement(b.a,{src:F.b+F.a+"/profile_bgs".concat(Object(O.a)(u.auth_id),"/profile_avatar.jpg"),name:u.first_name+" "+u.last_name,description:"\u0414\u043e\u0431\u0430\u0432\u043b\u0435\u043d: ".concat(u.creation_time)+(u.modification_time?" \u0438\u0437\u043c\u0435\u043d\u0435\u043d: ".concat(u.modification_time):"")})),h.a.createElement("div",{className:"record__text",ref:x,contentEditable:N},u.text),h.a.createElement("div",{className:"record__footer"},h.a.createElement("div",{className:"footer_item likes"},h.a.createElement("button",{type:"button",onClick:function(){c(u.me_liked,u.path)},className:"footer_icon default_img ".concat(u.me_liked?"liked_pic":"not_liked_pic")}),h.a.createElement("p",null,u.num_likes)),h.a.createElement("div",{className:"footer_item reply",onClick:function(){v(!0)}},h.a.createElement("button",{type:"button",className:"footer_icon default_img reply_pic"})),u.myId===u.auth_id&&h.a.createElement("div",{className:"footer_item post_comment__functional"},h.a.createElement("button",{type:"button",className:"footer_icon default_img post_comment_button__edit",onClick:I},N?"\u2713":"\ud83d\udd89"),h.a.createElement("button",{type:"button",className:"footer_icon default_img post_comment_button__delete",onClick:w},"\u274c"))),d&&h.a.createElement(j.a,{buttonMessage:"\u041e\u0442\u0432\u0435\u0442\u0438\u0442\u044c",onSend:function(e){v(!1),r(e,u.path)}}),a.map((function(e){return e&&h.a.createElement(J,Object.assign({updateCommentLike:c,key:e.path,addNewComment:r,deleteComment:i,handleError:m,myId:u.myId},e))}))))}t.default=Object(E.i)(R)},80:function(e,t,a){"use strict";var n=a(0),r=a.n(n),c=(a(81),a(57),Object(n.memo)((function(e){var t=e.placeholder,a=e.buttonMessage,c=e.onSend,o=Object(n.useRef)(),i=function(e){var t=e.target,a=o.current.value,n=o.current.selectionStart;o.current.value=a.slice(0,n)+t.innerText+a.slice(n,a.length),o.current.focus(),o.current.setSelectionRange(n+2,n+2)};return r.a.createElement("div",{className:"send_new_message__container"},r.a.createElement("textarea",{className:"default_search_input",placeholder:t,ref:o}),r.a.createElement("div",{className:"send_new_message__functional_area"},r.a.createElement("div",{className:"smile_area"},["\ud83d\ude10","\ud83d\ude2b","\ud83d\ude0e","\ud83d\ude02","\ud83d\ude21","\ud83d\ude2d","\ud83d\ude00","\ud83d\ude17","\ud83d\ude32","\ud83d\ude2c"].map((function(e){return r.a.createElement("span",{className:"smile",key:e,role:"img",onClick:i},e)}))),r.a.createElement("button",{className:"add_record__button",onClick:function(){c(o.current.value),o.current.value=""}},a)))})));t.a=c},81:function(e,t,a){},94:function(e,t,a){},95:function(e,t,a){e.exports=a.p+"static/media/f_add.7e42a7ae.ico"},96:function(e,t,a){e.exports=a.p+"static/media/f_friends.1639cc31.ico"},97:function(e,t,a){e.exports=a.p+"static/media/f_req_from.f634f58a.ico"},98:function(e,t,a){e.exports=a.p+"static/media/f_req_to.05fd42bf.ico"}}]);
//# sourceMappingURL=3.1c6f4305.chunk.js.map
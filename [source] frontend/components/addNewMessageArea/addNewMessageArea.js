import React, {memo, useRef} from "react";
import "./add_new_message_area.scss"
import "../../scss/default_blocks.scss"

const AddNewMessageArea = memo (({placeholder, buttonMessage, onSend}) => {
    const smiles = ['😐','😫','😎', '😂', '😡', '😭', '😀','😗','😲','😬']
    const ref = useRef()
    const _onSend = () => {
        onSend(ref.current.value);
        ref.current.value = "";
    }
    const onSmileAddHandler = ({target})=>{
        const oldMessage = ref.current.value;
        const selectionStart = ref.current.selectionStart;
        ref.current.value = oldMessage.slice(0, selectionStart) + target.innerText + oldMessage.slice(selectionStart, oldMessage.length);
        ref.current.focus();
        ref.current.setSelectionRange(selectionStart + 2, selectionStart + 2);
    }
    return (
        <div className={"send_new_message__container"}>
            <textarea className={"default_search_input"}
                      placeholder={placeholder}
                      ref={ref}
            />
            <div className={"send_new_message__functional_area"}>
                <div className={"smile_area"}>
                    {
                        smiles.map(smile =>
                            <span className={"smile"}
                                  key = {smile}
                                  role={"img"}
                                  onClick={onSmileAddHandler}
                            >
                                {smile}
                            </span>
                        )
                    }
                </div>
                <button className={"add_record__button"} onClick={_onSend}>{buttonMessage}</button>
            </div>
        </div>
    )
})

export default AddNewMessageArea;
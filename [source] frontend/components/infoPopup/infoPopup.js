import React from "react";
import {createPortal} from "react-dom";
import "./info_popup.scss";

export default class InfoPopup extends React.PureComponent {
    el = document.createElement('div');
    timer = setTimeout(this.props.handleClose, 5000);
    componentDidMount() {
        document.body.appendChild(this.el);
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps !== this.props){
            clearTimeout(this.timer);
        }
    }

    componentWillUnmount() {
        document.body.removeChild(this.el);
        clearTimeout(this.timer);
    }

    render() {
        return createPortal(
            <div className={"info_popup"} onClick={this.props.handleClose}>
                {this.props.text}
            </div>,
            this.el,
        );
    }
}

export function handleClose(){
    this.setState({error:null});
}

export function handleError(text) {
    this.setState({error: text});
}
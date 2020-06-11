import React, {createContext} from "react"
export const PlayerContext = createContext(undefined, undefined);

export class PlayerStoreHOC extends React.Component {
    state = {playlist: [],
             trackIndex: 0,
             done: false,
             playlistName: "UserMusic",
             withValue: "",
             userId: -1}

    updateMusic = (newState)=>{
        this.setState(newState)
    };

    render() {
        return (
            <PlayerContext.Provider
                value={{
                    playerStoreState : this.state,
                    updatePlayerStoreState: this.updateMusic
                }}
            >
                {this.props.children}
            </PlayerContext.Provider>
        );
    }
}
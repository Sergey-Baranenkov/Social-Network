import React from "react";
    import AudioPlayer, {RHAP_UI} from "react-h5-audio-player";
import "./header_audio_player.scss";
import PathFromIdGenerator from "../../../functools/PathFromIdGenerator";
import Fetcher from "../../../functools/Fetcher";
import {HTTP, ADDR} from "../../../address";

export default function HeaderAudioPlayer({playerStoreState, updatePlayerStoreState}){
    const onClickPrevious = () => {
        if (!isNaN(playerStoreState?.trackIndex) && Array.isArray(playerStoreState.playlist)){
            if (playerStoreState.trackIndex - 1 >= 0){
                updatePlayerStoreState({trackIndex: playerStoreState.trackIndex - 1})
            }
        }
    }

    const onClickNext = async () => {
        if (!isNaN(playerStoreState?.trackIndex) && Array.isArray(playerStoreState.playlist)){
            if (playerStoreState.trackIndex + 1 < playerStoreState.playlist.length){
                updatePlayerStoreState({trackIndex: playerStoreState.trackIndex + 1})
            }else if (playerStoreState.done === false){
                let url = HTTP + ADDR;
                if (playerStoreState.playlistName === "UserMusic"){
                    url += "/music/get_user_music";
                }else{
                    url += "/music/get_combined_music";
                }
                const params = {
                    limit: 3,
                    startFrom: playerStoreState.trackIndex + 1,
                    userId : playerStoreState.userId
                }

                if (playerStoreState.withValue !== ""){
                    params['withValue'] = playerStoreState.withValue
                }
                const [error, response] = await Fetcher(url,params);
                if (error === null){
                    if (response.Done){
                        updatePlayerStoreState({done: true, trackIndex: 0});
                    }else {
                        updatePlayerStoreState(
                            {playlist: [...playerStoreState.playlist, ...response[playerStoreState.playlistName]],
                            trackIndex: playerStoreState.trackIndex + 1,
                        })
                    }
                }else{console.log(error)}
            }else if (playerStoreState.done === true){
                updatePlayerStoreState({trackIndex: 0});
            }
        }
    }

    return (
        <div>
            <AudioPlayer
                autoPlay
                style={{width: "500px"}}
                showSkipControls={true}
                showJumpControls={false}
                src={
                    playerStoreState.playlist?.[playerStoreState.trackIndex]?
                        HTTP + ADDR + "/music_storage" +
                        PathFromIdGenerator(playerStoreState.playlist[playerStoreState.trackIndex].music_id) +
                        "/audio.mp3"
                        :
                        null
                    }

                onClickNext={onClickNext}
                onClickPrevious={onClickPrevious}
                onEnded={onClickNext}
                layout={"horizontal"}
                customAdditionalControls={[]}
                customProgressBarSection={[RHAP_UI.PROGRESS_BAR, RHAP_UI.DURATION]}
            />
        </div>

    )
}
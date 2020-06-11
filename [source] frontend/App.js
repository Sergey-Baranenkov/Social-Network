import React, {lazy, Suspense} from "react";
import {Route, Switch, Redirect} from "react-router-dom";
import "./components/static_sidebars/menu__static.scss";
import "./scss/page.scss"

import Loading from "./components/loader/LoadingPage";
import StaticSidebars from "./components/static_sidebars/staticSidebars";
import {PlayerContext} from "./PlayerContext";
import getCookie from "./functools/getCookie";
import AuthPage from "./components/auth_page/AuthPage";

const AboutMePage = lazy(()=> import("./components/about_me/AboutMePage"));
const GalleryPage = lazy(()=> import("./components/gallery/GalleryPage"));
const MessagesPage = lazy(()=> import("./components/messages/MessagesPage"));
const FriendsPage = lazy(()=> import("./components/friends/FriendsPage"));
const VideoPage = lazy(()=> import("./components/video/VideoPage"));
const AudioPage = lazy(()=> import("./components/audio/AudioPage"));
const WeatherPage = lazy(()=> import("./components/weather/WeatherPage"));
const SettingsPage = lazy(()=> import("./components/settings/SettingsPage"));
const ProfilePage = lazy(()=> import("./components/profile/ProfilePage"));

class App extends React.Component{
  render() {
    return (
        <>
            <Switch>

                <Route exact path="/">
                    <Redirect to ="/сообщения"/>
                </Route>

              <PrivateRoute path="/профиль/:id">
                  <Suspense fallback={<Loading/>}>
                    <ProfilePage/>
                  </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/about_me/:id">
                <Suspense fallback={<Loading/>}>
                  <AboutMePage />
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/настройки">
                <Suspense fallback={<Loading/>}>
                  <SettingsPage />
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/фотографии/:id">
                <Suspense fallback={<Loading/>}>
                  <GalleryPage />
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/погода">
                <Suspense fallback={<Loading/>}>
                  <WeatherPage />
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/музыка/:id">
                <Suspense fallback={<Loading/>}>
                  <PlayerContext.Consumer children={(props) => <AudioPage {...props}/>}/>
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/видео/:id">
                <Suspense fallback={<Loading/>}>
                  <VideoPage/>
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path={`/связи/:id`}>
                <Suspense fallback={<Loading/>}>
                  <FriendsPage/>
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/сообщения/:id">
                <Suspense fallback={<Loading/>}>
                  <MessagesPage/>
                </Suspense>
              </PrivateRoute>

              <PrivateRoute path="/сообщения">
                <Suspense fallback={<Loading/>}>
                  <MessagesPage/>
                </Suspense>
              </PrivateRoute>

              <Route path="/авторизация">
                <Suspense fallback={<Loading/>}>
                  <AuthPage/>
                </Suspense>
              </Route>
            </Switch>
        </>
    );
  }
}

function PrivateRoute({ children, ...rest }) {
  return (
      <Route
          {...rest}
          render={({ location }) =>
              !isNaN(+getCookie("userId")) ? (
                    <div className={"main_app"}>
                      <StaticSidebars/>
                      {children}
                    </div>
              ) : (
                  <Redirect
                      to={{
                        pathname: "/авторизация",
                        state: { from: location }
                      }}
                  />
              )
          }
      />
  );
}



export default App;
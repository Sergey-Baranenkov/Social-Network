class App extends React.Component {
    constructor(props){
        super(props);
        this.state = {test: false}
    }
    changeFunc = () => (
        this.setState((state, props) => ({test: !state.test}))
    )
    render() {
        return (
            <div>
                <button onClick = {this.changeFunc}>TEST</button>
                {this.state.test?<Registration/>:<Login/>}
            </div>

        );
    }
}

class Registration extends React.Component{
    constructor(props) {
        super(props);
    }
    render(){
        return(
            <form method="post" action = "/registration">
                <h1>REGISTRATION</h1>
                <input type="text" name ="email"/>
                <input type="text" name ="password" />
                <input type = "submit"/>
            </form>
        )
    }
}
class Login extends React.Component{
    constructor(props) {
        super(props);
    }
    render(){
        return(
            <form method="post" action = "/login">
                <h1>Login</h1>
                <input type="text" name ="email"/>
                <input type="text" name ="password" />
                <input type = "submit"/>
            </form>
        )
    }
}
ReactDOM.render( <App/>, document.querySelector("#root"));
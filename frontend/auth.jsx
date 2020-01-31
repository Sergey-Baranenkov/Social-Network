class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = { signIn: false };
  }
  changeFunc = () => {
    this.setState((state, props) => ({ signIn: !state.signIn }));
  };
  render() {
    return (
      <div className="containter">
        <div className="reg_form">
          <div className="auth_switcher">
            <label>
              <input
                defaultChecked="true"
                type="radio"
                name="auth_switcher"
                onChange={this.changeFunc}
              />
              <div>Sign Up</div>
            </label>
            <label>
              <input
                type="radio"
                name="auth_switcher"
                onChange={this.changeFunc}
              />
              <div>Log In</div>
            </label>
          </div>
          {this.state.signIn ? <Login /> : <Registration />}
        </div>
      </div>
    );
  }
}

class Registration extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    return (
      <div>
        <h1 className="auth_header">Sign Up For Free</h1>
        <form className="auth_form" method="post" action="/registration">
          <input
            type="text"
            name="first_name"
            placeholder="First Name"
            autoComplete="off"
          />
          <input
            type="text"
            name="last_name"
            placeholder="Last Name"
            autoComplete="off"
          />
          <input
            type="text"
            name="email"
            placeholder="Email Address"
            autoComplete="off"
            className="full_width"
          />
          <input
            type="text"
            name="password"
            placeholder="Password"
            autoComplete="off"
            className="full_width"
          />
          <input
            type="submit"
            className="full_width"
            id="signup_submit_button"
            value="GET STARTED"
          />
        </form>
      </div>
    );
  }
}
class Login extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    return (
      <div>
        <h1 className="auth_header">Welcome Back!</h1>
        <form className="auth_form" method="post" action="/login">
          <input
            type="text"
            name="email"
            placeholder="Email Address"
            autoComplete="off"
            className="full_width"
          />
          <input
            type="text"
            name="password"
            placeholder="Password"
            autoComplete="off"
            className="full_width"
          />
          <input
            type="submit"
            className="full_width"
            id="signup_submit_button"
            value="LOG IN"
          />
        </form>
      </div>
    );
  }
}
ReactDOM.render( <App/>, document.querySelector("#root"));
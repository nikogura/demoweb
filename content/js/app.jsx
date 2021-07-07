class App extends React.Component {
    render() {
        if (this.loggedIn) {
            return (<LoggedIn />);
        } else {
            return (<Home />);
        }
    }
}

class Home extends React.Component {
    render() {
        return (
            <div className="container">
                <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
                    <h1>Web Technology Demonstrator</h1>
                    <a href="/login" className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
                </div>
            </div>
        )
    }
}

class LoggedIn extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
        };

        this.logout = this.logout.bind(this);
    }

    logout() {
        localStorage.removeItem("id_token");
        localStorage.removeItem("access_token");
        localStorage.removeItem("profile");
        location.reload();
    }

    serverRequest() {
        fetch(window.location.href )
    }

    componentDidMount() {
      this.serverRequest();
    }

    render() {
        return (
            <div className="container">
                <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
                    <h1>Web Technology Demonstrator</h1>
                    <a onClick={this.logout} className="btn btn-primary btn-lg btn-login btn-block">Log Out</a>
                </div>
            </div>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('app'));

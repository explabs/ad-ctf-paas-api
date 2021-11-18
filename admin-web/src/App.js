import React, {Component} from "react";
import {connect} from "react-redux";
import {Route, Router, Switch} from "react-router-dom";

import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

import Login from "./components/login.component";
import HomePage from './pages/index';
import GeneratePage from './pages/generate';


import {logout} from "./actions/auth";
import {clearMessage} from "./actions/message";

import {history} from './helpers/history';

// import AuthVerify from "./common/auth-verify";
import EventBus from "./common/EventBus";
import {Nav, Navbar, NavDropdown} from "react-bootstrap";

class App extends Component {
  constructor(props) {
    super(props);
    this.logOut = this.logOut.bind(this);

    this.state = {
      currentUser: undefined,
    };

    history.listen((location) => {
      props.dispatch(clearMessage()); // clear message when changing location
    });
  }

    componentDidMount() {
        const user = this.props.user;

        if (user) {
            this.setState({
                currentUser: user,
            });
        }

        EventBus.on("logout", () => {
            this.logOut();
        });
    }

    componentWillUnmount() {
        EventBus.remove("logout");
    }

    logOut() {
        this.props.dispatch(logout());
        this.setState({
            currentUser: undefined,
        });
    }

    render() {
        const {currentUser} = this.state;

        return (
            <Router history={history}>
                <Navbar bg="light" expand="lg">
                    <Navbar.Brand className="ms-4">Admin Panel</Navbar.Brand>
                    <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                    <Navbar.Collapse id="basic-navbar-nav">
                        {currentUser ? (
                            <Nav className="me-auto container-fluid">
                                <Nav.Link href="/home">Teams</Nav.Link>
                                <NavDropdown title="Generate" id="basic-nav-dropdown">
                                    <NavDropdown.Item href="#action/3.1">Terraform</NavDropdown.Item>
                                    <NavDropdown.Item href="#action/3.2">VPN</NavDropdown.Item>
                                    <NavDropdown.Item href="generate-prom">Prometheus</NavDropdown.Item>
                                    <NavDropdown.Divider/>
                                    <NavDropdown.Item href="#action/3.5">All</NavDropdown.Item>
                                </NavDropdown>
                                <Nav.Link className="ms-auto" href="/login" onClick={this.logOut}>LogOut</Nav.Link>
                            </Nav>
                        ) : (
                            <Nav className="ms-auto">
                                <Nav.Link href="/login">Login</Nav.Link>
                            </Nav>
                        )}
                    </Navbar.Collapse>
                </Navbar>

                <div className="container mt-3">
                    <Switch>
                        <Route exact path={["/", "/home"]} component={HomePage}/>
                        <Route exact path="/login" component={Login}/>
                        <Route exact path="/generate-prom" component={GeneratePage}/>
                    </Switch>
                </div>

                {/* <AuthVerify logOut={this.logOut}/> */}
            </Router>
        );
    }
}

function mapStateToProps(state) {
    const {user} = state.auth;
    return {
        user,
    };
}

export default connect(mapStateToProps)(App);

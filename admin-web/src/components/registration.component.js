import React from "react";
import UserService from "../services/user.service"
import {Card, Form} from "react-bootstrap"

export class Registration extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            switch_state: true,
            items: {}
        };
        this.changeRegistrationStatus = this.changeRegistrationStatus.bind(this)
    }

    componentDidMount() {
        UserService.getPlatformInfo().then(
            response => {
                this.setState({items: response.data})
                if (response.data.reg_status === "open") {
                    this.setState({
                        switch_state: true,
                        isLoaded: true,
                    });
                } else {
                    this.setState({
                        switch_state: false,
                        isLoaded: true,
                    });
                }
            },
            error => {
                this.setState({
                    content:
                        (error.response &&
                            error.response.data &&
                            error.response.data.message) ||
                        error.message ||
                        error.toString()
                });
            }
        );
    }

    changeRegistrationStatus() {
        const {switch_state} = this.state;
        if (switch_state) {
            this.setState({switch_state: false})
            UserService.closeRegistration().then(
                response => {
                    if (response.data.reg_status === "open") {
                        this.setState({
                            switch_state: true,
                            isLoaded: true,
                        });
                    } else {
                        this.setState({
                            switch_state: false,
                            isLoaded: true,
                        });
                    }
                },
                error => {
                    this.setState({
                        content:
                            (error.response &&
                                error.response.data &&
                                error.response.data.message) ||
                            error.message ||
                            error.toString()
                    });
                }
            );
        } else {
            this.setState({switch_state: true})
            UserService.openRegistration().then(
                response => {
                    if (response.data.reg_status === "open") {
                        this.setState({
                            switch_state: true,
                            isLoaded: true,
                        });
                    } else {
                        this.setState({
                            switch_state: false,
                            isLoaded: true,
                        });
                    }
                },
                error => {
                    this.setState({
                        content:
                            (error.response &&
                                error.response.data &&
                                error.response.data.message) ||
                            error.message ||
                            error.toString()
                    });
                }
            );
        }
    }

    render() {
        const {error, isLoaded, items, switch_state} = this.state;
        if (error) {
            return <div>Ошибка: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Загрузка...</div>;
        } else {
            return (
                <  >
                    <div>
                        <Card
                            bg="light"
                            // className="mb-2"
                            style={{"width": "100%", /*margin: '15px'*/}}
                            border="primary"
                        >
                            <Card.Title>PlatformInfo</Card.Title>
                            <Card.Body>
                                <Card.Text>
                                    {Object.entries(items).map(([key, value]) => {
                                        return <li key={key}><b>{key}:</b> {value}</li>
                                    })}
                                </Card.Text>
                            </Card.Body>
                        </Card></div>
                    <Form>
                        <Form.Check
                            type="switch"
                            id="custom-switch"
                            label="Registration"
                            checked={switch_state}
                            onChange={this.changeRegistrationStatus}
                        />
                    </Form>
                </>
            );
        }
    }
}
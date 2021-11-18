import React from "react";
import UserService from "../services/user.service"
import {Form} from "react-bootstrap"

export class CompetitionComponent extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            switch_state: true
        };
        this.changeRegistrationStatus = this.changeRegistrationStatus.bind(this)
    }

    componentDidMount() {
        UserService.getPlatformInfo().then(
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

    changeRegistrationStatus() {
        const {switch_state} = this.state;
        if (switch_state) {
            this.setState({switch_state: false})
            UserService.startCompetition().then(
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
            UserService.stopCompetition().then(
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
        const {error, isLoaded, switch_state} = this.state;
        if (error) {
            return <div>Ошибка: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Загрузка...</div>;
        } else {
            return (
                <Form>
                    <Form.Check
                        type="switch"
                        id="custom-switch"
                        label="Checkers"
                        checked={switch_state}
                        onChange={this.changeRegistrationStatus}
                    />
                </Form>
            );
        }
    }
}
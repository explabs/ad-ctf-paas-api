import React from "react";
import UserService from "../services/user.service"
import {Button} from "react-bootstrap";

export class UpdateVPN extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            isLoaded: false,
            outputMessage: " ",
            items: {}
        };
        // this.updateVPN = this.updateVPN.bind(this)
    }

    updateVPNRoutes() {
        UserService.updateVPN().then(
            response => {
                console.log(response)
            },
            error => {
                console.log(error)
            }
        )

    }

    render() {
        return (
            <>
                <Button variant="primary" size="lg" onClick={this.updateVPNRoutes}>Update VPN</Button>{' '}
            </>
        );
    }
}
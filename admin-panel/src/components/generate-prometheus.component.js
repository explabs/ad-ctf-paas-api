import React from "react";
import {Button, ButtonGroup, Form} from "react-bootstrap";

function Generate() {
    return (
        <>
            <Form>
                <Form.Check
                    label="Checkers"
                    name="group1"
                    type="checkbox"
                    id="checkers"
                />
                <Form.Check
                    label="News"
                    name="group1"
                    type="checkbox"
                    id="news"
                />
                <Form.Check
                    label="Exploits"
                    name="group1"
                    type="checkbox"
                    id="exploits"
                />
            </Form>
            <ButtonGroup size="sm">
                <Button variant="outline-primary" size="sm">Autogenerate Prometheus Config</Button>
                <Button variant="outline-success" size="sm">Generate & Run</Button>
            </ButtonGroup>
        </>
    )
}

export default Generate;
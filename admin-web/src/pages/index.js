import React from 'react';

import {TeamList} from "../components/team-list.component";
import {Registration} from "../components/registration.component";
import {CompetitionComponent} from "../components/competition.component"


class HomePage extends React.Component {

    render() {
        return (
            <div>
                <TeamList/>
                <Registration/>
                <CompetitionComponent/>
            </div>
        );
    }
}

export default HomePage;
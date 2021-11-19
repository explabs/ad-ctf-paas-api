import React from 'react';

import {TeamList} from "../components/team-list.component";
import {Registration} from "../components/registration.component";
import {CompetitionComponent} from "../components/competition.component"
import {UpdateVPN} from "../components/update-vpn.component";


class HomePage extends React.Component {

    render() {
        return (
            <div>
                <TeamList/>
                <Registration/>
                <CompetitionComponent/>
                <UpdateVPN/>
            </div>
        );
    }
}

export default HomePage;
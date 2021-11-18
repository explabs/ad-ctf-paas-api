import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://192.168.100.105:8080/api/v1/";

class UserService {
    getTeamsList() {
        return axios.get(API_URL + "admin/teams", {headers: authHeader()});
    }

    getPlatformInfo() {
        return axios.get(API_URL + "services/system/info", );
    }

    closeRegistration() {
        return axios.get(API_URL + "admin/reg/close", {headers: authHeader()});
    }

    openRegistration() {
        return axios.get(API_URL + "admin/reg/open", {headers: authHeader()});
    }
    stopCompetition() {
        return axios.post(API_URL + "admin/prom/stop", {},{headers: authHeader()});
    }

    startCompetition() {
        return axios.post(API_URL + "admin/prom/start", {},{headers: authHeader()});
    }
}

export default new UserService();

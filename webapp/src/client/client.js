import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';

export default class Client {
    constructor() {
        this.url = '/plugins/bitbucket/api/v1';
    }

    getConnected = async (reminder = false) => {
        return this.doGet(`${this.url}/connected?reminder=` + reminder);
    };

    getReviews = async () => {
        return this.doGet(`${this.url}/reviews`);
    };

    getYourPrs = async () => {
        return this.doGet(`${this.url}/yourprs`);
    };

    getPrsDetails = async (prList) => {
        return this.doPost(`${this.url}/prsdetails`, prList);
    };

    getYourAssignments = async () => {
        return this.doGet(`${this.url}/yourassignments`);
    };

    getBitbucketUser = async (userID) => {
        return this.doPost(`${this.url}/user`, {user_id: userID});
    };

    getRepositories = async () => {
        return this.doGet(`${this.url}/repositories`);
    };

    createIssue = async (payload) => {
        return this.doPost(`${this.url}/createissue`, payload);
    };

    searchIssues = async (searchTerm) => {
        return this.doGet(`${this.url}/searchissues?term=${searchTerm}`);
    };

    attachCommentToIssue = async (payload) => {
        return this.doPost(`${this.url}/createissuecomment`, payload);
    };

    getIssue = async (owner, repo, issueId) => {
        return this.doGet(`${this.url}/issue?owner=${owner}&repo=${repo}&id=${issueId}`);
    };

    getPullRequest = async (owner, repo, prId) => {
        return this.doGet(`${this.url}/pr?owner=${owner}&repo=${repo}&id=${prId}`);
    };

    doGet = async (url, body, headers = {}) => {
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const options = {
            method: 'get',
            headers,
        };

        const response = await fetch(url, Client4.getOptions(options));

        if (response.ok) {
            return response.json();
        }

        const text = await response.text();

        throw new ClientError(Client4.url, {
            message: text || '',
            status_code: response.status,
            url,
        });
    };

    doPost = async (url, body, headers = {}) => {
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        const options = {
            method: 'post',
            body: JSON.stringify(body),
            headers,
        };

        const response = await fetch(url, Client4.getOptions(options));

        if (response.ok) {
            return response.json();
        }

        const text = await response.text();

        throw new ClientError(Client4.url, {
            message: text || '',
            status_code: response.status,
            url,
        });
    };
}

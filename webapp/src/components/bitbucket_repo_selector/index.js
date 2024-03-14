import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import manifest from 'manifest';

import {getRepos} from '../../actions';

import BitbucketRepoSelector from './bitbucket_repo_selector.jsx';

function mapStateToProps(state) {
    return {
        yourRepos: state[`plugins-${manifest.id}`].yourRepos,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getRepos,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(BitbucketRepoSelector);

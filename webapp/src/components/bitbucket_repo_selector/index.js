import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {id as pluginId} from 'manifest';
import {getRepos} from '../../actions';

import BitbucketRepoSelector from './bitbucket_repo_selector.jsx';

function mapStateToProps(state) {
    return {
        yourRepos: state[`plugins-${pluginId}`].yourRepos,
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

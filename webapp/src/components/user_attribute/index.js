import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getBitbucketUser} from '../../actions';

import manifest from '../../manifest';

import UserAttribute from './user_attribute.jsx';

function mapStateToProps(state, ownProps) {
    const id = ownProps.user ? ownProps.user.id : '';
    const {id: pluginId} = manifest;
    const user = state[`plugins-${pluginId}`].bitbucketUsers[id] || {};

    return {
        id,
        username: user.username,
        enterpriseURL: state[`plugins-${pluginId}`].enterpriseURL,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getBitbucketUser,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(UserAttribute);

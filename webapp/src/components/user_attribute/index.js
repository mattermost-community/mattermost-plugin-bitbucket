import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getBitbucketUser} from '../../actions';

import UserAttribute from './user_attribute.jsx';

function mapStateToProps(state, ownProps) {
    const id = ownProps.user ? ownProps.user.id : '';
    const user = state['plugins-bitbucket'].bitbucketUsers[id] || {};

    return {
        id,
        username: user.username,
        enterpriseURL: state['plugins-bitbucket'].enterpriseURL,
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

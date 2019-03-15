import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getReviews, getUnreads, getYourPrs, getYourAssignments} from '../../actions';

import SidebarButtons from './sidebar_buttons.jsx';

function mapStateToProps(state) {
    return {
        connected: state['plugins-bitbucket'].connected,
        username: state['plugins-bitbucket'].username,
        clientId: state['plugins-bitbucket'].clientId,
        reviews: state['plugins-bitbucket'].reviews,
        yourPrs: state['plugins-bitbucket'].yourPrs,
        yourAssignments: state['plugins-bitbucket'].yourAssignments,
        unreads: state['plugins-bitbucket'].unreads,
        enterpriseURL: state['plugins-bitbucket'].enterpriseURL,
        org: state['plugins-bitbucket'].organization,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getReviews,
            getUnreads,
            getYourPrs,
            getYourAssignments,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarButtons);

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getConnected, getReviews, getYourAssignments, getYourPrs, updateRhsState} from '../../actions';

import manifest from '../../manifest';

import SidebarButtons from './sidebar_buttons.jsx';

function mapStateToProps(state) {
    const {id: pluginId} = manifest;
    return {
        connected: state[`plugins-${pluginId}`].connected,
        reviews: state[`plugins-${pluginId}`].reviews,
        yourPrs: state[`plugins-${pluginId}`].yourPrs,
        yourAssignments: state[`plugins-${pluginId}`].yourAssignments,
        enterpriseURL: state[`plugins-${pluginId}`].enterpriseURL,
        showRHSPlugin: state[`plugins-${pluginId}`].rhsPluginAction,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getConnected,
            getReviews,
            getYourPrs,
            getYourAssignments,
            updateRhsState,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarButtons);

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getReviewsDetails, getYourPrsDetails} from '../../actions';
import {id as pluginId} from '../../manifest';

import SidebarRight from './sidebar_right.jsx';

function mapPrsToDetails(prs, details) {
    if (!prs) {
        return [];
    }

    return prs.map((pr) => {
        let foundDetails;
        if (details) {
            const repoUrl = pr.destination ? pr.destination.repository.links.html.href : pr.repository.links.html.href;
            foundDetails = details.find((prDetails) => {
                return (repoUrl === prDetails.url) && (pr.id === prDetails.id);
            });
        }

        if (!foundDetails) {
            return pr;
        }

        return {
            ...pr,
            participants: foundDetails.participants,
        };
    });
}

function mapStateToProps(state) {
    return {
        reviews: mapPrsToDetails(state[`plugins-${pluginId}`].reviews, state[`plugins-${pluginId}`].reviewsDetails),
        yourPrs: mapPrsToDetails(state[`plugins-${pluginId}`].yourPrs, state[`plugins-${pluginId}`].yourPrsDetails),
        yourAssignments: state[`plugins-${pluginId}`].yourAssignments || [],
        enterpriseURL: state[`plugins-${pluginId}`].enterpriseURL,
        rhsState: state[`plugins-${pluginId}`].rhsState,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getYourPrsDetails,
            getReviewsDetails,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(SidebarRight);

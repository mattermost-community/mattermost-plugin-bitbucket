import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {id as pluginId} from 'manifest';
import {closeCreateIssueModal, createIssue} from 'actions';

import CreateIssueModal from './create_issue';

const mapStateToProps = (state) => {
    const postId = state[`plugins-${pluginId}`].createIssueModalForPostId;
    const post = getPost(state, postId);

    return {
        visible: state[`plugins-${pluginId}`].isCreateIssueModalVisible,
        post,
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeCreateIssueModal,
    create: createIssue,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(CreateIssueModal);

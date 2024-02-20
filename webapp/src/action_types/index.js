import manifest from '../manifest';

const {id: pluginId} = manifest;
export default {
    RECEIVED_REPOSITORIES: pluginId + '_received_repositories',
    RECEIVED_REVIEWS: pluginId + '_received_reviews',
    RECEIVED_REVIEWS_DETAILS: pluginId + '_received_reviews_details',
    RECEIVED_YOUR_PRS: pluginId + '_received_your_prs',
    RECEIVED_YOUR_PRS_DETAILS: pluginId + '_received_your_prs_details',
    RECEIVED_YOUR_ASSIGNMENTS: pluginId + '_received_your_assignments',
    RECEIVED_MENTIONS: pluginId + '_received_mentions',
    RECEIVED_CONNECTED: pluginId + '_received_connected',
    RECEIVED_BITBUCKET_USER: pluginId + '_received_bitbucket_user',
    RECEIVED_SHOW_RHS_ACTION: pluginId + '_received_rhs_action',
    UPDATE_RHS_STATE: pluginId + '_update_rhs_state',
    CLOSE_CREATE_ISSUE_MODAL: pluginId + '_close_create_modal',
    OPEN_CREATE_ISSUE_MODAL: pluginId + '_open_create_modal',
    CLOSE_ATTACH_COMMENT_TO_ISSUE_MODAL: pluginId + '_close_attach_modal',
    OPEN_ATTACH_COMMENT_TO_ISSUE_MODAL: pluginId + '_open_attach_modal',
    RECEIVED_ATTACH_COMMENT_RESULT: pluginId + '_received_attach_comment',
};

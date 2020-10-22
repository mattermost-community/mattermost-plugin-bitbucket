import {connect} from 'react-redux';

import {id as pluginId} from 'manifest';

import {LinkTooltip} from './link_tooltip';

const mapStateToProps = (state) => {
    return {connected: state[`plugins-${pluginId}`].connected};
};

export default connect(mapStateToProps, null)(LinkTooltip);

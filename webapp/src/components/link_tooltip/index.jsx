import {connect} from 'react-redux';

import manifest from 'manifest';

import {LinkTooltip} from './link_tooltip';

const mapStateToProps = (state) => {
    return {connected: state[`plugins-${manifest.id}`].connected};
};

export default connect(mapStateToProps, null)(LinkTooltip);

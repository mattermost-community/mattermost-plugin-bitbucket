import React, {PureComponent} from 'react';
import PropTypes from 'prop-types';

import ReactSelectSetting from 'components/react_select_setting';

const initialState = {
    invalid: false,
    error: null,
};

export default class BitbucketRepoSelector extends PureComponent {
    static propTypes = {
        yourRepos: PropTypes.array.isRequired,
        theme: PropTypes.object.isRequired,
        onChange: PropTypes.func.isRequired,
        value: PropTypes.string,
        addValidate: PropTypes.func,
        removeValidate: PropTypes.func,
        actions: PropTypes.shape({
            getRepos: PropTypes.func.isRequired,
        }).isRequired,
    };

    constructor(props) {
        super(props);
        this.state = initialState;
    }

    componentDidMount() {
        this.props.actions.getRepos().then((result) => {
            if (result.error) {
                this.setState({
                    error: result.error.message,
                });
            }
        });
    }

    onChange = (name, newValue) => {
        this.props.onChange(newValue);
    };

    render() {
        const repoOptions = this.props.yourRepos.map((item) => ({value: item.name, label: item.full_name}));
        const {error} = this.state;

        let fetchingReposError = null;
        if (error) {
            fetchingReposError = (
                <p className='help-text error-text'>
                    <span>{error}</span>
                </p>
            );
        }

        return (
            <div className={'form-group margin-bottom x3'}>
                <ReactSelectSetting
                    name={'repo'}
                    label={'Repository'}
                    limitOptions={true}
                    required={true}
                    onChange={this.onChange}
                    options={repoOptions}
                    isMulti={false}
                    key={'repo'}
                    theme={this.props.theme}
                    addValidate={this.props.addValidate}
                    removeValidate={this.props.removeValidate}
                    value={repoOptions.find((option) => option.label === this.props.value)}
                />
                {fetchingReposError}
                <div className={'help-text'}>
                    {'Returns Bitbucket repositories connected to the user account'} <br/>
                </div>
            </div>
        );
    }
}

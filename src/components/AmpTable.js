import React from 'react';


class AmpTable extends React.Component {
    constructor() {
        super()

        this.displaySpecs = this.displaySpecs.bind(this);
    }

    displaySpecs() {
        let specs = this.props.specs.map(spec => {
            return(
                <tbody>
                    <tr key={spec.id}>
                        <td>{spec.band}</td>
                        <td>{spec.direction}</td>
                        <td>{spec.testName}</td>
                        <td>{spec.ip}</td>
                        <td>{spec.target}</td>
                        <td>{spec.lowLimit}</td>
                        <td>{spec.upLimit}</td>
                        <td>{spec.lEeprom}</td>
                        <td>{spec.uEeprom}</td>
                        <td>{spec.lowerFrequency}</td>
                        <td>{spec.upperFrequency}</td>
                    </tr>
                </tbody>
            )
        });
        return specs;
    }

    render() {
        return (
            <div>
                <table className="table table-hover table-bordered table-dark col-lg">
                    <thead>
                        <tr>
                            <th scope="col">Band</th>
                            <th scope="col">Direction</th>
                            <th scope="col">Test Name</th>
                            <th scope="col">Input Power</th>
                            <th scope="col">LowerLimit</th>
                            <th scope="col">UpperLimit</th>
                            <th scope="col">Target</th>
                            <th scope="col">E_LowerLimit</th>
                            <th scope="col">E_UpperLimit</th>
                            <th scope="col">StartFreq<div className="mhzSize">(MHz)</div></th>
                            <th scope="col">StopFreq<div className="mhzSize">(MHz)</div></th>
                        </tr>
                    </thead>
                    {this.displaySpecs()}
                </table>    
            </div>
        )
    }
}

export default AmpTable;
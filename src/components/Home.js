import React, { Component } from 'react';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';


import AmpTable from './AmpTable';

class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      modelBom: "",
       specs: [],
    };
    this.getAmpData = this.getAmpData.bind(this);
    this.setAmpModel = this.setAmpModel.bind(this);
  }

  getAmpData(e) {
    e.preventDefault();
    let params = this.refs.searchParams.value;
    axios({
      method:'get',
      url:`/search/${params}`,
      responseType:'json'
    })
    .then(response => {
      console.log(response);
      this.setState({ specs: response.data });
      this.refs.form.reset();
    })
    .catch(error => {
      alert("There was an issue" + error);
      console.log(error);
    });
  }

  setAmpModel(e) {
    let model = e.target.value;
    this.setState({ modelBom: model })
  }


  render() {
    return(
      <div>
        <div>
          <input type='text' ref='searchParams'></input>
        </div>
        <div>
          <button type='submit' onClick={this.getAmpData}>Submit</button>
        </div>
        <ol>
          {this.state.ampInfo.map(amp =>
            <li key={amp.id}>{amp.id.band}</li>
          )}
        </ol>
      </div>
    )
  }
};

export default Home;

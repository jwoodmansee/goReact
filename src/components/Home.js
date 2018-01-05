import React, { Component } from 'react';
import axios from 'axios';


class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      ampInfo: []
    }

    this.getAmpData = this.getAmpData.bind(this);


  }



  getAmpData() {
    let params = this.refs.searchParams.value;
    axios({
      method:'get',
      url:`http://localhost:8080/search/${params}`,
      responseType:'json'
    })
    .then(res =>{
      this.setState({ ampInfo: res.data });
      console.log(this.state)
    });
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

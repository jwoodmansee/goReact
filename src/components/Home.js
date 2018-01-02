import React, { Component } from 'react';
import axios from 'axios';

class Home extends Component {
  constructor(props) {
    super(props);
    this.getAmpData = this.getAmpData.bind(this);
  }



  getAmpData() {
    let params = this.refs.searchParams.value;
    axios({
      method:'get',
      url:`http://localhost:8080/search/${params}`,
      responseType:'json'
    })
    .then(function (response) {
      console.log(response)
    })
    .catch(function (error) {
      console.log(error)
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
      </div>
    )
  }
};

export default Home;

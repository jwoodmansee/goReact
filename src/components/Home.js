import React, { Component } from 'react';
import axios from 'axios';

class Home extends Component {

  

  getAmpData() {
    var apiUrl = 'http://localhost:8080/'
    axios.get(apiUrl + `{sreachParams}`)
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
          <input type='text' refs='sreachParams'></input>
        </div>
        <div>
          <button type='submit' onClick={this.getAmpData}>Submit</button>
        </div>
      </div>
    )
  }
};

export default Home;

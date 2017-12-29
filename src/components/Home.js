import React, { Component } from 'react';

class Home extends Component {

  buttonClicked() {
    console.log("Button was Clicked!")
  }

  render() {
    return(
      <div>
        <p>Hello World!</p>

        <div>
          <button type='submit' onClick={this.buttonClicked}>Submit</button>
        </div>
      </div>
    )
  }
};

export default Home;

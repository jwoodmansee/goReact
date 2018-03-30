import React from 'react';
import ReactDOM from 'react-dom';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import Home from './components/Home.js'

ReactDOM.render(
  <MuiThemeProvider>
    <Home />
  </MuiThemeProvider>,
  document.getElementById('app')
);

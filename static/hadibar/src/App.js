import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

var Welcome = React.createClass({
  getInitialState : function() {
    return {
      counter : 0
    };
  },

  render: function() {
    return <button onClick={this.handleClick}>{this.state.counter}</button>;
  },

  handleClick: function() {
    this.setState({
      counter: this.state.counter + 1
    })
  }
});

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to React</h2>
        </div>
        <Welcome />

        <p className="App-intro">
          Alles Klar
        </p>
      </div>
    );
  }
}

export default App;

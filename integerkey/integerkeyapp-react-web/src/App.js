import React from 'react';
import './App.css';
import IntegerKeyList from './components/IntegerKeyList';


function App() {


  const hostname =  window.location.hostname
  const port = 8080
  const url = 'http://' + hostname + ':' + port 

  console.log(url)


  return (
    <div className="integerkey-app">
      <h1>List of created Assets</h1>
     
      <IntegerKeyList url={url}/>
    </div>
  );
}

export default App;
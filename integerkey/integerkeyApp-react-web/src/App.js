import React, { useEffect, useState } from 'react';
import './App.css';
import IntegerKeyList from './components/IntegerKeyList';
import { FaSyncAlt } from 'react-icons/fa';
import {AiOutlineClear} from 'react-icons/ai';

function App() {

  const domain = "tally.tallysolutions.com"
  const hostname = "tbchlfdevpeer01"
  const port = 8080

  const url = 'http://' + hostname + '.' + domain + ':' + port 

  const [assets, setAssets] = useState([]); 

  const handleRefresh = async () => {
    const response = await fetch( url + '/integerKey/getAllAssets');
    const data = await response.json();
    setAssets(data);
  };

  useEffect(() => {
    handleRefresh();
  }, []);

  return (
    <div className="integerkey-app">
      <h1>List of created Assets</h1>
      <div className='buttons'>
                <button onClick={handleRefresh} className='refresh-button'>
                        <FaSyncAlt />
                </button>
                <button className='clearAll-button'>
                        <AiOutlineClear/>
                </button>
      </div>
      
      <IntegerKeyList url={url}/>
    </div>
  );
}

export default App;

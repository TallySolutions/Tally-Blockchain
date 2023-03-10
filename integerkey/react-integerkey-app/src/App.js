import React, { useEffect, useState } from 'react';
import './App.css';
import IntegerKeyList from './components/IntegerKeyList';
import { FaSyncAlt } from 'react-icons/fa';
import {AiOutlineClear} from 'react-icons/ai';

function App() {
        const [assets, setAssets] = useState([]); 

        const handleRefresh = async () => {
          const response = await fetch('http://20.219.112.54:8080/integerKey/getAllAssets');
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
            
            <IntegerKeyList assets={assets}/>
          </div>
        );
}

export default App;

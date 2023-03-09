import React, { useEffect, useState } from 'react';
import './App.css';
import IntegerKeyList from './components/IntegerKeyList';
import { FaSyncAlt } from 'react-icons/fa';

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
            <button onClick={handleRefresh} className='refresh-button'>
              <FaSyncAlt />
            </button>
            <IntegerKeyList assets={assets}/>
          </div>
        );
}

export default App;

import React, { useEffect, useState } from 'react';
import './App.css';
import IntegerKeyList from './components/IntegerKeyList';
import { FaSyncAlt } from 'react-icons/fa';
import {AiOutlineClear} from 'react-icons/ai';

function App() {


  const domain = "tally.tallysolutions.com"
  const hostname =  window.location.hostname
  const port = 8080
  const url = 'http://' + hostname + ':' + port 

  console.log(url)

  const [assets, setAssets] = useState([]); 

  const handleRefresh = async () => {
    const response = await fetch( url + '/integerKey/getAllAssets')
    .then(response => {
      if (response.ok){
              return response.json()
      }
      else{
          //asset.isComplete= true ;
         //  alert('The asset does not exist. Try reloading the list for the updated version.' );
          return console.error(response)
      }
      } )
      .then(data =>{
           
          //clear : assets.reamoveAll()
      
          //loop throug data array
        
          //create new asset objec, Name, Value and displayValue, add to assets
      
          //set Assets
    //        asset.Value= data["Value"]
    //        asset.displayValue= data["Name"] + " = " + data["Value"] 
    //        console.log(assets)
    //        setAssets(asssets)
    //   })

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
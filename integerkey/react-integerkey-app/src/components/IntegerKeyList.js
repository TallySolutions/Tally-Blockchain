import React, {useState} from 'react';
import IntegerKeyForm from './IntegerKeyForm';
import IntegerKey from './IntegerKey';


function IntegerKeyList() {

    const[assets, setAssets] = useState([]);

    const addAsset = asset => {
        
        if (!asset.assetname || /^\s*$/.test(asset.assetname) )   // making sure the name is valid
        {
            return
        }

        const newAssets = [asset, ...assets]
        setAssets(newAssets);
    }
    
    const incrementValue = asset =>{
        fetch('http://20.219.112.54:8080/integerKey/increaseValue',{  
                              method: 'POST',
                              headers: {
                                            'Content-Type': 'application/json' ,
                                            'Access-Control-Request-Method' : 'POST',
                                            'Access-Control-Request-Headers' : 'Content-Type'
                                        },
                              body: JSON.stringify({
                                Name: asset.assetname,
                                Value: "1"
                              })})
                              .then(response => {
                                            if (response.ok){
                                                    return response.json()
                                            }
                                            else{
                                                //asset.isComplete= true ;
                                                alert('You cannot go above 20' );
                                                return console.error(response)
                                            }
                              } )
                              .then(data =>{
                                    asset.Value= data["Value"]
                                    asset.displayValue= data["Name"] + " = " + data["Value"] 
                                    console.log(assets)
                                    // asset.isComplete = true ;
                                    // setAssets(assets) 
                                    updateAsset(asset.id, asset.Value)

                                    // completeAsset(asset.id)
                                })
                                     
                }
                const decrementValue = asset =>{
                     fetch('http://20.219.112.54:8080/integerKey/decreaseValue',{  
                                           method: 'POST',
                                           headers: {
                                                         'Content-Type': 'application/json' ,
                                                         'Access-Control-Request-Method' : 'POST',
                                                         'Access-Control-Request-Headers' : 'Content-Type'
                                                     },
                                           body: JSON.stringify({
                                             Name: asset.assetname,
                                             Value: "1"
                                           })})
                                           .then(response => {
                                                         if (response.ok){
                                                                 return response.json()
                                                         }
                                                         else{
                                                             //asset.isComplete= true ;
                                                             alert('You cannot go below 0' );
                                                             return console.error(response)
                                                         }
                                           } )
                                           .then(data =>{
                                                 asset.Value= data["Value"]
                                                 asset.displayValue= data["Name"] + " = " + data["Value"] 
                                                 console.log(assets)
                                                 updateAsset(asset.id, asset.Value)
                                             })
                                                  
                             }

    const removeAsset = assetname =>{
        const removeArr = [...assets].filter(asset=> asset.assetname !== assetname);
        fetch(`http://20.219.112.54:8080/integerKey/deleteAsset/${assetname}`, {
                                    method: 'DELETE',
                            })
                            .then(response => {
                                if (response.ok){
                                  return response.json()
                                }
                                else{
                                  alert('Asset was already removed! Try reloading the list to ensure there are no non-existing assets' );
                                  return console.error(response)
                                }
                              } )
                            .then(data => console.log(data))
                            .catch(error => console.error(error))
                setAssets(removeArr);
    }
    
    
     const updateAsset = (id, Value) =>{
        let updatedAssets= assets.map(asset =>{
            if(asset.id===id){
                asset.Value = Value
                return asset
            }
        });
        setAssets(updatedAssets)
     }

  return (
    <div>
        <h1>List of created Assets</h1>
        <IntegerKeyForm onSubmit={addAsset}/>
        <IntegerKey
            assets={assets}
            incrementValue={incrementValue}
            decrementValue = {decrementValue}
            removeAsset = {removeAsset}
        />
    </div>
  )
}

export default IntegerKeyList ;
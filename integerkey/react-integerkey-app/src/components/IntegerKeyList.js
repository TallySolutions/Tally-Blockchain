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


    const removeAsset = assetname =>{
        const removeArr = [...assets].filter(asset=> asset.assetname !== assetname);
        fetch(`http://20.219.112.54:8080/integerKey/deleteAsset/${assetname}`, {
                                    method: 'DELETE',
                            })
                            .then(response => response.json())
                            .then(data => console.log(data))
                            .catch(error => console.error(error))
                setAssets(removeArr);
    }


    const updateAsset = (assetId, newValue)=> {
        if (!newValue.assetname || /^\s*$/.test(newValue.assetname) )   // making sure the name is valid
        {
            return
        }

        setAssets(prev => prev.map(item=> (item.id === assetId ? newValue : item)))
    }
    
     const completeAsset = id =>{
        let updatedAssets= assets.map(asset =>{
            if(asset.id===id){
                asset.isComplete = !asset.isComplete
                return asset.assetname
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
            completeAsset = {completeAsset}
            removeAsset = {removeAsset}
            updateAsset = {updateAsset}
        />
    </div>
  )
}

export default IntegerKeyList
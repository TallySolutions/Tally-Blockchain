import React, {useState} from 'react'
import IntegerKeyForm from './IntegerKeyForm'
import {RiCloseCircleLine} from 'react-icons/ri'
import {TiArrowSortedUp} from 'react-icons/ti'

function IntegerKey({assets, completeAsset, removeAsset, updateAsset}) {

    const [edit, setEdit] = useState({
        id : null,
        value : ''
    })

    const submitUpdate = value => {
        updateAsset(edit.id, value)
        setEdit({
            id: null,
            value:''
        })
    }

    if (edit.id){
        return <IntegerKeyForm edit={edit} onSubmit={submitUpdate} />
    }


  return assets.map((asset, index)=> (
    <div className={asset.isComplete ? 'asset-row complete': 'asset-row'} key={index}>
            <div key={asset.id} onClick={()=>completeAsset(asset.id)}>

                {asset.assetname}
            </div>

            <div key={asset.id} onClick={()=>completeAsset(asset.id)}>
                {asset.value}
            </div> 

            <div className='icons'>
                <RiCloseCircleLine  
                
                onClick={
                    
                    () => removeAsset(asset.assetname) 
                    
                }
                        className='delete-icon'
                />
                {/* buttons for add and subtract */}
                <TiArrowSortedUp
 
                onClick={() => 
                    setEdit({id : asset.id , assetname : asset.assetname}
                        )}

                        className='edit-icon'
                />
            </div>
    </div>
  ));
};

export default IntegerKey ;
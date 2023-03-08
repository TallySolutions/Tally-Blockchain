import React, {useState} from 'react'
import IntegerKeyForm from './IntegerKeyForm'
import {RiCloseCircleLine} from 'react-icons/ri'
import {TiArrowSortedUp} from 'react-icons/ti'
import {TiArrowSortedDown} from 'react-icons/ti'
import {incrementValue} from './IntegerKeyList'
import {decrementValue} from './IntegerKeyList'

function IntegerKey({assets, completeAsset, removeAsset, updateAsset}) {

    // const [edit, setEdit] = useState({
    //     id : null,
    //     value : ''
    // })

    // const submitUpdate = value => {
    //     updateAsset(edit.id, value)
    //     setEdit({
    //         id: null,
    //         value:''
    //     })
    // }

    // if (edit.id){
    //     return <IntegerKeyForm edit={edit} onSubmit={submitUpdate} />
    // }

    const handleIncrementClick = (asset) => {
        incrementValue(asset);
      };
      const handleDecrementClick = (asset) => {
        decrementValue(asset);
      };


  return assets.map((asset, index)=> (
    <div className={asset.isComplete ? 'asset-row complete': 'asset-row'} key={index}>
            <div key={asset.id} onClick={()=>completeAsset(asset.id)}>

                {asset.assetname}
            </div>

            <div key={asset.id} onClick={()=>completeAsset(asset.id)}>
                {asset.value}
            </div> 

            <div className='icons'>
                {/* buttons for add and subtract */}
                <TiArrowSortedUp
 
                        onClick={() => handleIncrementClick(asset)}

                        className='edit-icon'
                />
                <TiArrowSortedDown
                        onClick={
                            () => handleDecrementClick(asset)
                                }
                            className='edit-icon'
                />
                <RiCloseCircleLine  
                            onClick={
                                () => removeAsset(asset.assetname) 
                            }
                            className='delete-icon'
                />
            </div>
    </div>
  ));
};

export default IntegerKey ;
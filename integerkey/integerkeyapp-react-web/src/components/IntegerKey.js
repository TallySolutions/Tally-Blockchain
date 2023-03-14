import React from 'react'
//import IntegerKeyForm from './IntegerKeyForm'
import { RiCloseCircleLine } from 'react-icons/ri'
import { TiArrowSortedUp } from 'react-icons/ti'
import { TiArrowSortedDown } from 'react-icons/ti'



function IntegerKey({ assets, incrementValue, decrementValue, removeAsset }) {


  var top = 0;
  return assets.list.map((asset, index) => (
      <div className='asset-row-container'>

        <div className='asset-row' key={index}>

          <div key={asset.Name}>

            {asset.displayValue}

          </div>

          <div className='icons'>
            {/* buttons for add and subtract, deletion */}
            <TiArrowSortedUp

              onClick={
                () => incrementValue(asset)
              }
              className='edit-icon'
            />
            <TiArrowSortedDown
              onClick={
                () => decrementValue(asset)
              }
              className='edit-icon'
            />
            <RiCloseCircleLine
              onClick={
                () => removeAsset(asset)
              }
              className='delete-icon'
            />
          </div>
        </div>
        <div className={asset.isUpdating ? 'asset-row overlay' : 'asset-row no-overlay'}></div>
      </div>
    )
  );
};

export default IntegerKey;
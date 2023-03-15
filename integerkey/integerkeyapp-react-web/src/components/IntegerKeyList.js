import React, { useState } from 'react';
import IntegerKeyForm from './IntegerKeyForm';
import IntegerKey from './IntegerKey';
import { FaSyncAlt } from 'react-icons/fa';


class IntegerKeyList extends React.Component {

  setAsset(assets) {
    var state = this.state;
    state.Assets.list = assets;

    this.setState(state);
  }

  setAssetUpdating(asset) {
    asset.isUpdating = true;
    var state = this.state;
    this.setState(state);
    this.operation_in_progress = true;
  }
  setAssetUpdated(asset) {
    asset.isUpdating = false;
    var state = this.state;
    this.setState(state);
    this.operation_in_progress = false;
  }


  constructor(props) {

    super(props);

    this.operation_in_progress = false;

    this.url = props.url;

    this.state = {
      Assets: {
        list: []
      }
    };



    this.addAsset = asset => {

      if (!asset.Name || /^\s*$/.test(asset.Name))   // making sure the name is valid
      {
        return
      }

      const newAssets = [asset, ...this.state.Assets.list]
      this.setAsset(newAssets);
    };

    this.incrementValue = asset => {
      this.setAssetUpdating(asset);
      fetch(this.url + '/integerKey/increaseValue', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Request-Method': 'POST',
          'Access-Control-Request-Headers': 'Content-Type'
        },
        body: JSON.stringify({
          Name: asset.Name,
          Value: "1" , 
          Owner: asset.Owner
        })
      })
        .then(response => {
          if (response.ok) {
             return response.json()
          }
          else {
            return {error: "Error increasing asset value."}
          }
        }).then(data => {
          if(data.hasOwnProperty('error')){
               alert('Error: ' +  data.error );  
               this.setAssetUpdated(asset);            
          }else{
            console.log(JSON.stringify(data));
            asset.Value = data.Value
            asset.displayValue = data.Name + " = " + data.Value + " | owned by " + data.Owner;
            this.updateAsset(asset.Name, asset.Value, asset.Owner)
          }
        });

    };
    this.decrementValue = asset => {
      this.setAssetUpdating(asset);
      fetch(this.url + '/integerKey/decreaseValue', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Request-Method': 'POST',
          'Access-Control-Request-Headers': 'Content-Type'
        },
        body: JSON.stringify({
          Name: asset.Name,
          Value: "1",
          Owner: asset.Owner
        })
      })
      .then(response => {
        if (response.ok) {
           return response.json()
        }
        else {
          return {error: "Error decreasing asset value."}
        }
      }).then(data => {
        if(data.hasOwnProperty('error')){
             alert('Error: ' +  data.error );  
             this.setAssetUpdated(asset);            
        }else{
          console.log(JSON.stringify(data));
          asset.Value = data.Value
          asset.displayValue = data.Name + " = " + data.Value + " | owned by " + data.Owner;
          this.updateAsset(asset.Name, asset.Value, asset.Owner)
        }
      });

    };

    this.removeAsset = asset => {
      var assetname = asset.Name;

      var intKeyList = this;

      const removeArr = [...this.state.Assets.list].filter(asset => asset.Name !== assetname);
      asset.isUpdating = true;
      this.setAssetUpdating(asset);

      fetch(this.url + '/integerKey/deleteAsset/' + assetname, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Request-Method': 'DELETE',
          'Access-Control-Request-Headers': 'Content-Type'
        },
      })
        .then(response => {
          asset.isUpdating = false;
          if (response.ok) {
            intKeyList.setAsset(removeArr);
          }
          else {
            alert('Error removing asset.' );
            this.setAssetUpdated(asset);            
          }
        });
    };


    this.updateAsset = (Name, Value, Owner) => {
      let updatedAssets = this.state.Assets.list.map(asset => {
        if (asset.Name === Name) {
          asset.Name = Name;
          asset.Value = Value;
          asset.Owner= Owner;
          asset.displayValue = Name + ' = ' + Value + " | owned by " + Owner;
          asset.isUpdating = false;
          return asset;
        }else{
          //asset.isUpdating = false;
          return asset;
        }
      });
      console.log(updatedAssets)
      this.setAsset(updatedAssets)
      this.operation_in_progress = false;
    };

    this.handleAutoRefresh = () => {
      this.handleRefresh(true);
    }
    this.handleManualRefresh = () => {
      this.handleRefresh(false);
    }
    this.handleRefresh = isAuto => {

      if(this.operation_in_progress) return;

      var intKeyList = this;
      fetch(this.url + '/integerKey/getAllAssets')
        .then(response => {
          if (response.ok) {
            return response.json();
          }
          else {
            return {error: 'Internal error'};
          }
        })
        .then(data => {
          if(data.hasOwnProperty('error')){
            if (!isAuto){
               alert('Error: ' +  data.error );  
            }            
          }else{
             var assets = [];
             //loop throug data array
             data.forEach(function (obj) {
               //create new asset objec, Name, Value, Owner and displayValue, add to assets
               var asset = { Name: obj.Name, Value: obj.Value, Owner: obj.Owner, displayValue: obj.Name + " = " + obj.Value + " | owned by "+ obj.Owner}
               assets.push(asset);
             });
   
             //set Assets
             intKeyList.setAsset(assets);
          }
        });
    };


  }


  componentDidMount() {
    console.log("Component loaded")
    this.handleAutoRefresh();
    //setInterval(this.handleAutoRefresh, 500);
  }

  render() {
    return (
      <div >
        <div className='buttons'>
          <button onClick={this.handleManualRefresh} className='refresh-button'>
            <FaSyncAlt />
          </button>
        </div>

        <IntegerKeyForm onSubmit={this.addAsset} url={this.url} />
        <IntegerKey
          assets={this.state.Assets}
          incrementValue={this.incrementValue}
          decrementValue={this.decrementValue}
          removeAsset={this.removeAsset}
        />
      </div>
    );
  }
}

export default IntegerKeyList;
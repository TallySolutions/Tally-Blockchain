import React, { useState } from 'react';
import IntegerKeyForm from './IntegerKeyForm';
import IntegerKey from './IntegerKey';
import { FaSyncAlt } from 'react-icons/fa';
import { AiOutlineClear } from 'react-icons/ai';

class IntegerKeyList extends React.Component {

  setAsset(assets) {
    var state = this.state;
    state.Assets.list = assets;

    this.setState(state);
  }

  constructor(props) {

    super(props);

    this.url = props.url;

    this.state = {
      Assets: {
        list: []
      }
    };



    this.addAsset = asset => {

      if (!asset.assetname || /^\s*$/.test(asset.assetname))   // making sure the name is valid
      {
        return
      }

      const newAssets = [asset, ...this.state.Assets.list]
      this.setAsset(newAssets);
    };

    this.incrementValue = asset => {
      // if (asset.Value >= 20){
      //     alert('You cannot have an asset with a value higher than 20.')
      //     return
      // }
      fetch(this.url + '/integerKey/increaseValue', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Request-Method': 'POST',
          'Access-Control-Request-Headers': 'Content-Type'
        },
        body: JSON.stringify({
          Name: asset.assetname,
          Value: "1"
        })
      })
        .then(response => {
          if (response.ok) {
            return response.json()
          }
          else {
            //asset.isComplete= true ;
            alert('The asset does not exist. Try reloading the list for the updated version.');
            return console.error(response)
          }
        })
        .then(data => {
          asset.Value = data["Value"]
          asset.displayValue = data["Name"] + " = " + data["Value"]
          console.log(this.state.Assets)
          // asset.isComplete = true ;
          // setAssets(assets) 
          this.updateAsset(asset.id, asset.Value)

          // completeAsset(asset.id)
        })

    };
    this.decrementValue = asset => {
      // if (asset.Value <= 0){
      //     alert('You cannot have an asset with a value lesser than 0.')
      // }
      fetch(this.url + '/integerKey/decreaseValue', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Request-Method': 'POST',
          'Access-Control-Request-Headers': 'Content-Type'
        },
        body: JSON.stringify({
          Name: asset.assetname,
          Value: "1"
        })
      })
        .then(response => {
          if (response.ok) {
            return response.json()
          }
          else {
            //asset.isComplete= true ;
            //  alert('The asset does not exist. Try reloading the list for the updated version.' );
            return console.error(response)
          }
        })
        .then(data => {
          asset.Value = data["Value"]
          asset.displayValue = data["Name"] + " = " + data["Value"]
          console.log(this.state.Assets)
          this.updateAsset(asset.id, asset.Value)
        })

    };

    this.removeAsset = assetname => {

      var intKeyList = this;

      const removeArr = [...this.state.Assets.list].filter(asset => asset.assetname !== assetname);
      fetch(this.url + '/integerKey/deleteAsset/' + assetname, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Request-Method': 'DELETE',
          'Access-Control-Request-Headers': 'Content-Type'
        },
      })
        .then(response => {
          if (response.ok) {
            return response.json()
          }
          else {
            //   alert('Asset was already removed! Try reloading the list to ensure there are no non-existing assets' );
            return console.error(response)
          }
        })
        .then(data => {
          console.log(data)
          intKeyList.setAsset(removeArr);
        })
        .catch(error => console.error(error))
    };


    this.updateAsset = (id, Value) => {
      let updatedAssets = this.state.Assets.list.map(asset => {
        if (asset.id === id) {
          asset.Value = Value
          return asset
        }
      });
      this.setAsset(updatedAssets)
    };

    this.handleRefresh = () => {

      var intKeyList = this;
      console.log(this);
      fetch(this.url + '/integerKey/getAllAssets')
        .then(response => {
          if (response.ok) {
            return response.json()
          }
          else {
            //asset.isComplete= true ;
            //  alert('The asset does not exist. Try reloading the list for the updated version.' );
            return console.error(response)
          }
        })
        .then(data => {
          var assets = [];
          console.log("DATA => " + data);
          //loop throug data array
          data.forEach(function (obj) {
            console.log(obj.Name);
            //create new asset objec, Name, Value and displayValue, add to assets
            var asset = { assetname: obj.Name, Value: obj.Value, displayValue: obj.Name + " = " + obj.Value }
            console.log(assets);
            assets.push(asset);
          });

          //set Assets
          intKeyList.setAsset(assets);
          console.log(assets)
        })

    };
  }


  componentDidMount() {
    console.log("Component loaded")
    this.handleRefresh();

    setInterval(this.handleRefresh, 500);
  }

  render() {
    return (
      <div >
        <div className='buttons'>
          <button onClick={this.handleRefresh} className='refresh-button'>
            <FaSyncAlt />
          </button>
          <button className='clearAll-button'>
            <AiOutlineClear />
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
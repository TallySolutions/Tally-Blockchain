import React, {useState} from 'react' // useState is imported in order to track data involved
import nextId from "react-id-generator";


function IntegerKeyForm(props) {

    const [input_asset, setInputAsset] = useState('')  // name is taken as input
    // useState({name :"" value:""}) ............ return (<>asset.name</h1> <p>asset.value</p>) -- when taking in >1 inputs

    const handleChange = e =>{
        setInputAsset(e.target.value);
    }



    const handleSubmit = e=>{
        e.preventDefault();  // preventing reloading of the page on clicking button
        fetch('http://20.219.112.54:8080/integerKey/createAsset',{  // handle error- server down
                              method: 'PUT',
                              headers: {
                                            'Content-Type': 'application/json' ,
                                            'Access-Control-Request-Method' : 'PUT',
                                            'Access-Control-Request-Headers' : 'Content-Type'
                                        },
                              body: JSON.stringify({
                                Name: input_asset
                              })
                })//test
                .then(response => {
                  if (response.ok){
                    return response.json()
                  }
                  else{
                    return console.error(response)
                  }
                } )
                .then(data =>{
                           props.onSubmit({
                            id: nextId("asset-id-"),
                            assetname: data["Name"],
                            assetvalue: data["Value"]
                           });
                           setInputAsset('');
                      });
                      // .error(e=>{
                      //   console.error(e)});
                            
};


  return (
    <form className='integerkey-form' onSubmit={handleSubmit}>

        <input type="text" placeholder="Add asset name" value={input_asset} name="assetname" className='integerkey-input' onChange={handleChange}/> 

        <button className='integerkey-button'>Add asset</button>
    </form>
  )
}

export default IntegerKeyForm;
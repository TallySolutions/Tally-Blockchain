import React, {useState} from 'react' // useState is imported in order to track data involved
import nextId from "react-id-generator";



function IntegerKeyForm(props) {
    const [input_asset, setInputAsset] = useState('')  // name is taken as input

    const handleChange = e =>{
        setInputAsset(e.target.value);
    }

    const handleSubmit = e=>{
        e.preventDefault();  // preventing reloading of the page on clicking button
        fetch('http://20.219.112.54:8080/integerKey/createAsset',{  
                              method: 'PUT',
                              headers: {
                                            'Content-Type': 'application/json' ,
                                            'Access-Control-Request-Method' : 'PUT',
                                            'Access-Control-Request-Headers' : 'Content-Type'
                                        },
                              body: JSON.stringify({
                                Name: input_asset
                              })
                })
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
                                id : nextId("asset-id:"),
                                assetname: data["Name"] ,
                                assetvalue: data["Value"] ,
                                displayValue: data["Name"] + " = " + data["Value"]
                           });
                           setInputAsset('');
                      }); 
                      // add error handler- to deal with asset creation error  (start with alert)                       
};


  return (
    <form className='integerkey-form' onSubmit={handleSubmit}>

        <input type="text" placeholder="Add asset name" value={input_asset} name="assetname" className='integerkey-input' onChange={handleChange}/> 

        <button className='integerkey-button'>Add asset</button>
    </form>
  )
};  

export default IntegerKeyForm;
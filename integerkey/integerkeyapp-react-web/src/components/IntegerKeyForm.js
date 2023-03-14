import React, {useState} from 'react' // useState is imported in order to track data involved

function IntegerKeyForm(props) {
    const [input_asset, setInputAsset] = useState('')  // name is taken as input

    const handleChange = e =>{
        setInputAsset(e.target.value);
    }

    const handleSubmit = e=>{
        e.preventDefault();  // preventing reloading of the page on clicking button
        fetch( props.url + '/integerKey/createAsset',{  
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
                    alert('Error in creating asset!!! Try a different name.' );
                    return console.error(response)
                  }
                } )
                .then(data =>{
                           props.onSubmit({
                                Name: data["Name"] ,
                                displayValue: data["Name"] + " = " + data["Value"],
                                isUpdating: false
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
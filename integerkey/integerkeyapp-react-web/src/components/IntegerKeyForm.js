import React, {useState} from 'react' // useState is imported in order to track data involved

function IntegerKeyForm(props) {
    const [input_asset, setInputAsset] = useState('')  // name is taken as input

    const [input_owner, setInputOwner]= useState('')  // owner is taken as input

    const handleChangeName = e =>{
        setInputAsset(e.target.value);
    }
    const handleChangeOwner = e =>{
      setInputOwner(e.target.value);
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
                                Name: input_asset,
                                Owner: input_owner
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
                                Name: data.Name ,
                                Owner: data.Owner,
                                displayValue: data["Name"] + " = " + data["Value"] + " | owned by " + data["Owner"],
                                isUpdating: false
                           });
                           setInputAsset('');
                      }); 
                      // add error handler- to deal with asset creation error  (start with alert)                       
};


// NOTE: FOR NOW, taking the owner name as input- instead of identifying the peer

  return (
    <form className='integerkey-form' onSubmit={handleSubmit}>
        <div className='integerkey-form-inputs'>
              <input type="text" placeholder="Add asset name" value={input_asset} 
                     name="assetname" className='integerkey-input' onChange={handleChangeName}
              />

              <input type="text" placeholder="Add owner name" value={input_owner} 
                     name="assetowner" className='integerkey-input' onChange={handleChangeOwner}
              /> 
        </div>
        <button className='integerkey-button'>Add asset</button>
    </form>
  )
};  

export default IntegerKeyForm;
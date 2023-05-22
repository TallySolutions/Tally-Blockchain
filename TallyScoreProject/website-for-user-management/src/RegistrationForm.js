import React, { useState } from 'react';

function RegistrationForm({ onNewRegistration }) {
  const [name, setName] = useState('');
  const [address,setAddress] = useState('');
  const [pan, setPan] = useState('');
  const [license, setLicense] = useState('');
  const [status, setStatus] = useState('');

  const handleFormSubmit = (e) => {
    e.preventDefault();
    const registration = {
      name: name,
      address: address,
      pan: pan,
      license: license,
      score: 500,
      status: status,
    };
    onNewRegistration(registration);
    setName('');
    setAddress('');
    setPan('');
    setLicense('');
    setStatus('');

    // add GSTN
  };

  return (
    <form onSubmit={handleFormSubmit}>
      <label>
        Name:
        <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        Address:
        <input type="text" value={address} onChange={(e)=>setAddress(e.target.value)}/>
      </label>
      <label>
        PAN:
        <input type="text" value={pan} onChange={(e) => setPan(e.target.value)} />
      </label>
      <label>
        License:
        <select value={license} onChange={(e) => setLicense(e.target.value)}>
          <option value="none">Select a license type</option>
          <option name="Silver"> Silver</option>
          <option name="Gold">Gold</option>
        </select>
      </label>
      <label>
        Status:
        <select value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="none">Select status</option>
          <option name="active"> active</option>
          <option name="inactive">inactive</option>
        </select>
      </label>
      <button type="submit">Submit</button>
    </form>
  );
}

export default RegistrationForm;

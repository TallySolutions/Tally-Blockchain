import React, { useState } from 'react';

function RegistrationForm({ onNewRegistration }) {
  const [pan, setPan] = useState('');
  const [name, setName] = useState('');
  const[phonenumber, setPhoneNumber] = useState('');
  const [address,setAddress] = useState('');
  const [license, setLicense] = useState('');
  const [status, setStatus] = useState('');

  const handleFormSubmit = (e) => {
    e.preventDefault();
    const registration = {
      pan: pan,
      name: name,
      phonenumber: phonenumber,
      address: address,
      license: license,
      score: 500,
      status: status,
    };
    onNewRegistration(registration);
    setPan('');
    setName('');
    setPhoneNumber('');
    setAddress('');
    setLicense('');
    setStatus('');

    // add GSTN
  };

  return (
    <form onSubmit={handleFormSubmit}>
      <label>
        Business's PAN:
        <input type="text" value={pan} onChange={(e) => setPan(e.target.value)} />
      </label>
      <label>
        Name:
        <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        Phone Number:
        <input type="tel" value={phonenumber} onChange={(e)=>setPhoneNumber(e.target.value)}/>
      </label>
      <label>
        Address:
        <input type="text" value={address} onChange={(e)=>setAddress(e.target.value)}/>
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

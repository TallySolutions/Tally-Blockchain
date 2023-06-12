import React, { useState } from 'react';

function RegistrationForm({ onNewRegistration }) {
  const [PAN, setPan] = useState('');
  const [Name, setName] = useState('');
  const[PhoneNo, setPhoneNumber] = useState('');
  const [Address,setAddress] = useState('');
  const [LicenseType, setLicense] = useState('');
  const[Score, setScore]= useState('');
  const [status, setStatus] = useState('');
  const [businessCertDetails,setbusinessCertDetails] = useState('');    // change to details structure
  const handleFormSubmit = (e) => {
    e.preventDefault();
    const registration = {
      PAN: PAN,
      Name: Name,
      PhoneNo: PhoneNo,
      Address: Address,
      LicenseType: LicenseType,
      Score: 500,
      // status: status,
    };
    onNewRegistration(registration);

    const forRequest = {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(registration),
    };
    fetch('http://43.204.226.103:8080/TallyScoreProject/performRegistration', forRequest)
    .then(response => {
      if (response.ok) {
         return response.json()
      }
      else {
        return {error: "Error in registration."}
      }
    }).then(data => {
      setbusinessCertDetails(data)
      console.log("businessCertDetails:" + JSON.stringify(data));
    }).catch(error => {
      console.error('Error:', error);
    });
    setPan('');
    setName('');
    setPhoneNumber('');
    setAddress('');
    setLicense('');
    setScore('');
    setStatus('');
    
  };

  return (
    <form onSubmit={handleFormSubmit}>
      <label>
        Business's PAN:
        <input type="text" value={PAN} onChange={(e) => setPan(e.target.value)} />
      </label>
      <label>
        Name:
        <input type="text" value={Name} onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        Phone Number:
        <input type="tel" value={PhoneNo} onChange={(e)=>setPhoneNumber(e.target.value)}/>
      </label>
      <label>
        Address:
        <input type="text" value={Address} onChange={(e)=>setAddress(e.target.value)}/>
      </label>
      <label>
        License:
        <select value={LicenseType} onChange={(e) => setLicense(e.target.value)}>
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

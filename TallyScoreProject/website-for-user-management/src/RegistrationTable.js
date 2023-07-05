import React, { useState } from 'react';
import OwnerVoucherDialog from './DialogBoxes/OwnerVoucherDialog';
import SupplierVoucherDialog from './DialogBoxes/SupplierVoucherDialog';

function RegistrationTable({ registrations, setRegistrations }) {

  const [ownerDialogOpen, setOwnerDialogOpen] = useState(false);
  const [ownerDialogPAN, setOwnerDialogPAN] = useState('');  
  const [supplierDialogOpen, setSupplierDialogOpen] = useState(false);
  const [supplierDialogPAN, setSupplierDialogPAN] = useState('');  

const handleOwnerVoucherButtonClick = (pan) => {
  setOwnerDialogOpen(true);
  setOwnerDialogPAN(pan);
};

const handleSupplierVoucherButtonClick = (pan) => {
  setSupplierDialogOpen(true);
  setSupplierDialogPAN(pan);
};

const handleOwnerDialogClose = () => {
  setOwnerDialogOpen(false);
};

const handleSupplierDialogClose = () => {
  setSupplierDialogOpen(false);
};

  const handleIncrementClick = (registration) => {
    const updateStruct = {
      PAN: registration.PAN,
      ChangeVal: "10"
    };

    const forRequest = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updateStruct),
    };

    fetch('http://43.204.226.103:8080/TallyScoreProject/increaseTallyScore', forRequest)
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Error in registration.');
        }
      })
      .then(data => {
        const updatedRegistrations = registrations.map(reg =>
          reg.PAN === registration.PAN ? { ...reg, Score: reg.Score + parseInt(updateStruct.ChangeVal) } : reg
        );
        setRegistrations(updatedRegistrations);
        console.log("businessCertDetails:", JSON.stringify(data));
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

  const handleDecrementClick = (registration) => {
    const updateStruct = {
      PAN: registration.PAN,
      ChangeVal: "10"
    };

    const forRequest = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updateStruct),
    };

    fetch('http://43.204.226.103:8080/TallyScoreProject/decreaseTallyScore', forRequest)
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Error in registration.');
        }
      })
      .then(data => {
        const updatedRegistrations = registrations.map(reg =>
          reg.PAN === registration.PAN ? { ...reg, Score: reg.Score - parseInt(updateStruct.ChangeVal) } : reg
        );
        setRegistrations(updatedRegistrations);
        console.log("businessCertDetails:", JSON.stringify(data));
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

  return (
    <table>
      <thead>
        <tr>
          <th>PAN</th>
          <th>Name</th>
          <th>Phone No.</th>
          <th>Address</th>
          <th>License</th>
          <th>Score</th>
          <th>Status</th>
          <th>Generate Voucher as Owner</th>
          <th>Generate Voucher as Supplier</th>
        </tr>
      </thead>
      <tbody>
        {registrations.map((registration) => (
          <tr key={registration.PAN}>
            <td id="PAN-column">{registration.PAN}</td>
            <td>{registration.Name}</td>
            <td>{registration.PhoneNo}</td>
            <td>{registration.Address}</td>
            <td>{registration.LicenseType}</td>
            <td><div id ="scorerow">
                    <button id="scorebutton" onClick={() => handleDecrementClick(registration)}>-</button>
                    {registration.Score}
                    <button id="scorebutton" onClick={() => handleIncrementClick(registration)}>+</button>
              </div>
            </td>
            <td>{registration.status}</td>
            <td>
        <div>
          <button id="owner-voucher-button"
            className='generate-voucher-button'
          onClick={() => handleOwnerVoucherButtonClick(registration.PAN)}>Generate</button>
        </div>
      </td>
      <td>
        <div>
          <button id="supplier-voucher-button" 
          className='generate-voucher-button' 
          onClick={() => handleSupplierVoucherButtonClick(registration.PAN)}>Generate</button>
        </div>
      </td>

      {ownerDialogOpen && (
        <div className="dialog-overlay">
          <div className="dialog-container">
            <OwnerVoucherDialog onClose={handleOwnerDialogClose} pan={ownerDialogPAN}/>
          </div>
        </div>
      )}

      {supplierDialogOpen && (
        <div className="dialog-overlay">
          <div className="dialog-container">
            <SupplierVoucherDialog onClose={handleSupplierDialogClose} pan={supplierDialogPAN}/>
          </div>
        </div>
      )}
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default RegistrationTable;

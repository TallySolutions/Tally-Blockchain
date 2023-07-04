import React, { useState } from 'react';

function UpdateVoucherDialog({ onClose, pan }) {
  const [formData, setFormData] = useState({
    voucherID: '',
    update: '',
    newValue: '',
  });

  const [voucherDetails, setVoucherDetails] = useState('');
  const [showUpdateButton, setShowUpdateButton] = useState(false);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (event) => { // to verify voucher existance
    event.preventDefault();
    console.log('Form submitted:', formData);
    fetch(`http://43.204.226.103:8080/TallyScoreProject/readVoucher/${pan}?voucherID=${formData.voucherID}`, {  // calling readvoucher endpoint- in order to verify voucher
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Error in reading voucher.');
        }
      })
      .then((data) => {
        console.log('Read Voucher Response:', data);
        setVoucherDetails(JSON.stringify(data));
        setShowUpdateButton(true);
      })
      .catch((error) => {
        alert('Error while reading voucher');
        console.error('Error:', error);
      });  
  };

  const handleUpdateVoucher = () => {                                              // THIS IS THE FINAL BUTTON THAT UPDATES THE VOUCHER
    fetch(`http://43.204.226.103:8080/TallyScoreProject/voucherUpdation/${pan}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body:  JSON.stringify({
        VoucherID: formData.voucherID.trim(),
        Parameter: formData.update,
        UpdatedValue: formData.newValue
      }),
    })
      .then((response) => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Error in voucher updation.');
        }
      })
      .then((data) => {
        console.log('Update Voucher Response:', data);
        alert('Update of voucher successful!');
        onClose();
      })
      .catch((error) => {
        alert("Error while updating voucher")
        console.error('Error:', error);
      });
    console.log('Update Voucher Submitted:', formData);
    onClose();
  };

  const handleButtonClick = (action) => {
    console.log('Update Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
  };

  return (
    <div>
      <button
        className="close-dialog-button"
        onClick={() => handleButtonClick('Back')}
      >
        Back
      </button>
      <form onSubmit={handleSubmit}>
        <div className="update-voucher-form-group">
          <label htmlFor="voucherID">Voucher ID:</label>
          <input
            type="text"
            id="voucherID"
            name="voucherID"
            value={formData.voucherID}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="update-voucher-form-group">
          <label htmlFor="update">Update:</label>
          <select
            id="update"
            name="update"
            value={formData.update}
            onChange={handleInputChange}
            required
          >
            <option value="">Select an option</option>
            <option value="Hashcode">Hashcode</option>
            <option value="TotalValue">TotalValue</option>
          </select>
        </div>
        <div className="update-voucher-form-group">
          <label htmlFor="newValue">New Value:</label>
          <input
            type="text"
            id="newValue"
            name="newValue"
            value={formData.newValue}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="update-voucher-form-buttons">
          <button type="submit">Submit</button>
        </div>
      </form>
      {voucherDetails && (
        <div className="display-voucher-field-update">
          <input
            className="update-voucher-details-display"
            type="text"
            value={voucherDetails}
            readOnly
          />
        </div>
      )}
      {showUpdateButton && (
        <div className="update-voucher-form-buttons">
          <button
            className="update-voucher-button"
            onClick={handleUpdateVoucher}
          >
            Complete Update
          </button>
        </div>
      )}
    </div>
  );
}

export default UpdateVoucherDialog;

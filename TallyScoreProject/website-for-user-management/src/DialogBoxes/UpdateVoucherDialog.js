import React, { useState } from 'react';

function UpdateVoucherDialog({ onClose }) {
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

  const handleSubmit = (event) => {                           // THIS IS THE VERIFICATION PART (submits the form and displays voucher in box for verification by user)
    event.preventDefault();
    console.log('Form submitted:', formData);
                                              // CALL THE READ VOUCHER ENDPOINT HERE AND DISPLAY THAT INSTEAD OF SAMPLE VOUCHER DETAILS
    setVoucherDetails('Sample voucher details');
    setShowUpdateButton(true);
  };

  const handleUpdateVoucher = () => {                                              // THIS IS THE FINAL BUTTON THAT UPDATES THE VOUCHER
                                          // CALL VOUCHER UPDATE ENDPOINT HERE
    console.log('Voucher updated:', formData);
    alert('Voucher has been updated.');
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

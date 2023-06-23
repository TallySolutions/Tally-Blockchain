import React, { useState } from 'react';

function CancelVoucherDialog({ onClose }) {
  const [formData, setFormData] = useState({
    voucherID: '',
  });

  const [voucherDetails, setVoucherDetails] = useState('');

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
                                                                // CALL READVOUCHER ENDPOINT HERE
    console.log('Form submitted:', formData);
    setVoucherDetails('o/p generated from read voucher'); 
  };

  const handleButtonClick = (action) => {
    console.log('Cancel Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
    if (action=== 'Cancel Voucher'){
        onClose();
                                                            // REPLACE THE ABOVE LINE WITH CANCEL VOUCHER PAI ENDPOINT CALL 
        alert('Voucher set to cancelled.')
    }
  };

  return (
    <div className="cancel-voucher-dialog">
      <button
        className="close-dialog-button"
        onClick={() => handleButtonClick('Back')}
      >
        Back
      </button>
      <form onSubmit={handleSubmit}>
        <div className="cancel-voucher-form-group">
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
        <div className="cancel-voucher-form-buttons">
          <button id='verify-voucher-button' type="submit">Verify Voucher</button>
        </div>
      </form>
      <div className='display-voucher-field'>
            <input
                className="voucher-details-display"
                type="text"
                value={voucherDetails}
                readOnly
            />
      </div>
      <div className="cancel-voucher-form-buttons">
        <button
          className="cancel-voucher-button"
          onClick={() => handleButtonClick('Cancel Voucher')}
        >
          Cancel Voucher
        </button>
      </div>
    </div>
  );
}

export default CancelVoucherDialog;

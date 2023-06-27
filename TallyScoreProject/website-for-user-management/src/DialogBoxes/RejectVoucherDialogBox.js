import React, { useState } from 'react';

function RejectVoucherDialogBox({ onClose }) {
  const [formData, setFormData] = useState({
    voucherID: '',
  });

  const [voucherDetails, setVoucherDetails] = useState('');
  const [isFormSubmitted, setIsFormSubmitted] = useState(false);
  const [showVoucherDetails, setShowVoucherDetails] = useState(false);
  const [showRejectVoucherButton, setShowRejectVoucherButton] = useState(false);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    setIsFormSubmitted(true);
    setShowVoucherDetails(true);
    // CALL READ VOUCHER ENDPOINT HERE
    console.log('Form submitted:', formData);
    setVoucherDetails('o/p generated from read voucher');
    setShowRejectVoucherButton(true);
  };

  const handleButtonClick = (action) => {
    console.log('Reject Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    } else if (action === 'Reject Voucher') {
      if (formData.voucherID.trim() !== '') {
        onClose();
        // REPLACE THE ABOVE LINE WITH REJECT VOUCHER API ENDPOINT CALL
        alert('Voucher rejected.');
      } else {
        alert('Please enter a voucher ID.');
      }
    }
  };

  return (
    <div className="reject-voucher-dialog">
      <button className="close-dialog-button" onClick={() => handleButtonClick('Back')}>
        Back
      </button>
      <form onSubmit={handleSubmit}>
        <div className="reject-voucher-form-group">
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
        <div className="reject-voucher-form-buttons">
          <button id="verify-voucher-button" type="submit">
            Verify Voucher
          </button>
        </div>
      </form>
      {showVoucherDetails && (
        <div className="display-voucher-field">
          <input className="voucher-details-display" type="text" value={voucherDetails} readOnly />
        </div>
      )}
      {showRejectVoucherButton && (
        <div className="reject-voucher-form-buttons">
          <button
            className="reject-voucher-button"
            onClick={() => handleButtonClick('Reject Voucher')}
          >
            Confirm Rejection
          </button>
        </div>
      )}
      {isFormSubmitted && formData.voucherID.trim() === '' && (
        <p className="error-message">Please enter a voucher ID.</p>
      )}
    </div>
  );
}

export default RejectVoucherDialogBox;

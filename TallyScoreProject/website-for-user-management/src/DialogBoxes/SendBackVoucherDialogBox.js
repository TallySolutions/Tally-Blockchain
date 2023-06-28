import React, { useState } from 'react';

function SendBackVoucherDialogBox({ onClose }) {
  const [formData, setFormData] = useState({
    voucherID: '',
  });

  const [voucherDetails, setVoucherDetails] = useState('');
  const [isFormSubmitted, setIsFormSubmitted] = useState(false);
  const [showVoucherDetails, setShowVoucherDetails] = useState(false);
  const [showSendBackVoucherButton, setShowSendBackVoucherButton] = useState(false);

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
    setShowSendBackVoucherButton(true);
  };

  const handleButtonClick = (action) => {
    console.log('Send Back Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
    if (action === 'Send Back Voucher') {
      if (formData.voucherID.trim() !== '') {
        onClose();
        // REPLACE THE ABOVE LINE WITH SEND BACK VOUCHER API ENDPOINT CALL
        alert('Voucher sent back successfully.');
      } else {
        alert('Please enter a voucher ID.');
      }
    }
  };

  return (
    <div>
      <button className="close-dialog-button" onClick={() => handleButtonClick('Back')}>
        Back
      </button>
      <form onSubmit={handleSubmit}>
        <div className="send-back-voucher-form-group">
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
        <div className="send-back-voucher-form-buttons">
          <button id="sendback-verify-voucher-button" type="submit">
            Verify Voucher
          </button>
        </div>
      </form>
      {showVoucherDetails && (
        <div className="display-voucher-field">
          <input className="voucher-details-display" type="text" value={voucherDetails} readOnly />
        </div>
      )}
      {showSendBackVoucherButton && (
        <div className="send-back-voucher-form-buttons">
          <button
            className="send-back-voucher-button"
            onClick={() => handleButtonClick('Send Back Voucher')}
          >
            Send Back Voucher
          </button>
        </div>
      )}
      {isFormSubmitted && formData.voucherID.trim() === '' && (
        <p className="error-message">Please enter a voucher ID.</p>
      )}
    </div>
  );
}

export default SendBackVoucherDialogBox;

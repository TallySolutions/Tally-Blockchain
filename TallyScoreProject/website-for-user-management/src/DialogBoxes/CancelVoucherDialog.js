import React, { useState } from 'react';

function CancelVoucherDialog({ onClose }) {
  const [formData, setFormData] = useState({
    voucherID: '',
  });

  const [voucherDetails, setVoucherDetails] = useState('');
  const [isFormSubmitted, setIsFormSubmitted] = useState(false);
  const [showVoucherDetails, setShowVoucherDetails] = useState(false);
  const [showCancelVoucherButton, setShowCancelVoucherButton] = useState(false);

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
                                            // CALL READVOUCHER ENDPOINT HERE
    console.log('Form submitted:', formData);
    setVoucherDetails('o/p generated from read voucher');
    setShowCancelVoucherButton(true);
  };

  const handleButtonClick = (action) => {
    console.log('Cancel Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
    if (action === 'Cancel Voucher') {
      if (formData.voucherID.trim() !== '') {
        onClose();
                                                    // REPLACE THE ABOVE LINE WITH CANCEL VOUCHER API ENDPOINT CALL
        alert('Voucher set to cancelled.');
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
        <div>
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
      {showCancelVoucherButton && (
        <div className="cancel-voucher-form-buttons">
          <button
            className="cancel-voucher-button"
            onClick={() => handleButtonClick('Cancel Voucher')}
          >
            Cancel Voucher
          </button>
        </div>
      )}
      {isFormSubmitted && formData.voucherID.trim() === '' && (
        <p className="error-message">Please enter a voucher ID.</p>
      )}
    </div>
  );
}

export default CancelVoucherDialog;
import React, { useState } from 'react';

function ApproveVoucherDialog({ onClose }) {
  const [formData, setFormData] = useState({
    voucherID: '',
  });

  const [voucherDetails, setVoucherDetails] = useState('');
  const [isFormSubmitted, setIsFormSubmitted] = useState(false);
  const [showVoucherDetails, setShowVoucherDetails] = useState(false);
  const [showApproveVoucherButton, setShowApproveVoucherButton] = useState(false);

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
    setShowApproveVoucherButton(true);
  };

  const handleButtonClick = (action) => {
    console.log('Approve Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
    if (action === 'Approve Voucher') {
      if (formData.voucherID.trim() !== '') {
        onClose();
        // REPLACE THE ABOVE LINE WITH APPROVE VOUCHER API ENDPOINT CALL
        alert('Voucher set to approved.');
      } else {
        alert('Please enter a voucher ID.');
      }
    }
  };

  return (
    <div className="approve-voucher-dialog">
      <button className="close-dialog-button" onClick={() => handleButtonClick('Back')}>
        Back
      </button>
      <form onSubmit={handleSubmit}>
        <div className="approve-voucher-form-group">
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
        <div className="approve-voucher-form-buttons">
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
      {showApproveVoucherButton && (
        <div className="approve-voucher-form-buttons">
          <button
            className="approve-voucher-button"
            onClick={() => handleButtonClick('Approve Voucher')}
          >
            Approve Voucher
          </button>
        </div>
      )}
      {isFormSubmitted && formData.voucherID.trim() === '' && (
        <p className="error-message">Please enter a voucher ID.</p>
      )}
    </div>
  );
}

export default ApproveVoucherDialog;

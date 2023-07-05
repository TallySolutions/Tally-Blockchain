import React, { useState } from 'react';

function CancelVoucherDialogBox({ onClose, pan }) {
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
        setShowCancelVoucherButton(true);
      })
      .catch((error) => {
        alert('Error while reading voucher');
        console.error('Error:', error);
      });
  };

  const handleButtonClick = (action) => {
    console.log('Cancel Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
    if (action === 'Cancel Voucher') {
      if (formData.voucherID.trim() !== '') {
        fetch(`http://43.204.226.103:8080/TallyScoreProject/voucherCancellation/${pan}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Request-Method': 'POST',
            'Access-Control-Request-Headers': 'Content-Type'
          },
          body: JSON.stringify({
            VoucherID: formData.voucherID.trim()
          })
        })
        .then(response =>{
          if(response.ok){
            return response.json()
          }
          else{
            return {error: "Error while cancelling voucher"}
          }
        }).then(data => {
          if (data.hasOwnProperty('error')){
            alert('Error' + data.error);
            console.error("Error:", data.error);
          }
          else{
            console.log(JSON.stringify(data))
            alert("Voucher has been cancelled successfully!")
            onClose();
          }
        }
        );
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

export default CancelVoucherDialogBox;

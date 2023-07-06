import React, { useState } from 'react';

function ApproveVoucherDialogBox({ onClose, pan }) {
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
    console.log('Form submitted:', formData);
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
        setShowApproveVoucherButton(true);
      })
      .catch((error) => {
        alert('Error while reading voucher');
        console.error('Error:', error);
      });
  };

  const handleButtonClick = (action) => {
    console.log('Approve Voucher Button Clicked:', action);
    if (action === 'Back') {
      onClose();
    }
    if (action === 'Approve Voucher') {
      if (formData.voucherID.trim() !== '') {
        fetch(`http://43.204.226.103:8080/TallyScoreProject/voucherApproval/${pan}`, {
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
            return {error: "Error while approving voucher"}
          }
        }).then(data => {
          if (data.hasOwnProperty('error')){
            alert('Error' + data.error);
            console.error("Error:", data.error);
          }
          else{
            console.log(JSON.stringify(data))
            alert("Voucher has been approved successfully!")
            onClose();
          }
        }
        );
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

export default ApproveVoucherDialogBox;

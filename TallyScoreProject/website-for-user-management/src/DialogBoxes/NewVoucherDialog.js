import React, { useState } from 'react';

function NewVoucherDialog({ onClose }) {
  const [formData, setFormData] = useState({
    voucherId: '',
    supplierId: '',
    voucherType: '',
    hashCode: '',
    totalValue: '',
    currency: '',
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({ ...prevData, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
                                                        // CALL ENDPOINT TO CREATE NEW VOUCHER
    console.log('New Voucher Submitted:', formData);
    alert('New voucher creation successful!')
    onClose();
  };

  return (
    <div className="new-voucher-dialog">
      <button className="close-dialog-button" onClick={onClose}>
        Back
      </button>
      <form className='new-voucher-form' onSubmit={handleSubmit}>
        <div className="new-voucher-form-group">
          <label htmlFor="voucherId">Voucher ID:</label>
          <input
            type="text"
            id="voucherId"
            name="voucherId"
            value={formData.voucherId}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="new-voucher-form-group">
          <label htmlFor="supplierId">Supplier ID:</label>
          <input
            type="text"
            id="supplierId"
            name="supplierId"
            value={formData.supplierId}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="new-voucher-form-group">
          <label htmlFor="voucherType">Voucher Type:</label>
          <select
            id="voucherType"
            name="voucherType"
            value={formData.voucherType}
            onChange={handleInputChange}
            required
          >
            <option value="">Select Voucher Type</option>
            <option value="contra">contra</option>
            <option value="purchase">purchase</option>
            <option value="receipt">receipt</option>
            <option value="payment">payment</option>
            <option value="sales">sales</option>
            <option value="credit">credit</option>
            <option value="debit">debit</option>
            <option value="journal">journal</option>
          </select>
        </div>
        <div className="new-voucher-form-group">
          <label htmlFor="hashCode">Hash Code:</label>
          <input
            type="text"
            id="hashCode"
            name="hashCode"
            value={formData.hashCode}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="new-voucher-form-group">
          <label htmlFor="totalValue">Total Value:</label>
          <input
            type="text"
            id="totalValue"
            name="totalValue"
            value={formData.totalValue}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="new-voucher-form-group">
          <label htmlFor="currency">Currency:</label>
          <input
            type="text"
            id="currency"
            name="currency"
            value={formData.currency}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="new-voucher-form-submit">
          <button type="submit">Submit</button>
        </div>
      </form>
    </div>
  );
}

export default NewVoucherDialog;

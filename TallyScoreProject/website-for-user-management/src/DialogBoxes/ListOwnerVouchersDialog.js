import React, { useEffect, useState } from 'react';

function ListOwnerVouchersDialog({ onClose, pan }) {
  const [vouchersList, setVouchersList] = useState([]);

  useEffect(() => {
    const fetchVouchersList = () => {
      fetch(`http://43.204.226.103:8080/TallyScoreProject/listOwnerVouchers/${pan}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then((response) => response.json())
        .then((data) => setVouchersList(data))
        .catch((error) => console.error('Error fetching vouchers list:', error));
    };

    fetchVouchersList();
  }, [pan]);

  return (
    <div className="list-vouchers-dialog">
      <h3 className="dialog-title">List of Vouchers as Owner</h3>
      <table className="vouchers-table">
        <thead>
          <tr>
            <th>Voucher ID</th>
            <th>Supplier ID</th>
            <th>Type</th>
            <th>State</th>
          </tr>
        </thead>
        <tbody>
          {vouchersList.map((voucher, index) => (
            <tr key={index}>
              <td>{voucher.VoucherID}</td>
              <td>{voucher.SupplierID}</td>
              <td>{voucher.VoucherType}</td>
              <td>{voucher.State}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <button className="close-list-dialog-button" onClick={onClose}>
        Close
      </button>
    </div>
  );
}

export default ListOwnerVouchersDialog;

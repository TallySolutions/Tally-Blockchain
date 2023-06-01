import React from 'react';

function RegistrationTable({ registrations }) {
  return (
    <table>
      <thead>
        <tr>
          <th>PAN</th>
          <th>Name</th>
          <th>Phone No.</th>
          <th>Address</th>
          <th>License</th>
          <th>Score</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {registrations.map((registration) => (
          <tr key={registration.PAN}>
            <td>{registration.PAN}</td>
            <td>{registration.Name}</td>
            <td>{registration.PhoneNo}</td>
            <td>{registration.Address}</td>
            <td>{registration.LicenseType}</td>
            <td>{registration.Score}</td>
            <td>{registration.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default RegistrationTable;

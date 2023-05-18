import React from 'react';

function RegistrationTable({ registrations }) {
  return (
    <table>
      <thead>
        <tr>
          <th>User ID</th>
          <th>Name</th>
          <th>PAN</th>
          <th>License</th>
          <th>Score</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {registrations.map((registration) => (
          <tr key={registration.userid}>
            <td>{registration.userid}</td>
            <td>{registration.name}</td>
            <td>{registration.pan}</td>
            <td>{registration.license}</td>
            <td>{registration.score}</td>
            <td>{registration.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default RegistrationTable;

import React, { useState } from 'react';
import RegistrationForm from './RegistrationForm';
import RegistrationTable from './RegistrationTable';

function App() {
  const [registrations, setRegistrations] = useState([]);

  const handleNewRegistration = (registration) => {
    const newRegistration = {
      ...registration,
    };
    setRegistrations([...registrations, newRegistration]);
  };

  return (
    <div className="App">
      <RegistrationForm onNewRegistration={handleNewRegistration} />
      <RegistrationTable
        registrations={registrations}
        setRegistrations={setRegistrations} // Pass setRegistrations as a prop
      />
    </div>
  );
}


export default App;

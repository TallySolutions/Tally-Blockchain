import React, { useState } from 'react';
import RegistrationForm from './RegistrationForm';
import RegistrationTable from './RegistrationTable';
import { v4 as uuid } from 'uuid';

function App() {
  const [registrations, setRegistrations] = useState([]);

  const handleNewRegistration = (registration) => {
    const newRegistration = {
      ...registration,
      userid: generateUniqueId(),
    };
    setRegistrations([...registrations, newRegistration]);
  };

  const generateUniqueId = () => {
    return uuid();
  };

  return (
    <div className="App">
      <RegistrationForm onNewRegistration={handleNewRegistration} />
      <RegistrationTable registrations={registrations} />
    </div>
  );
}

export default App;

import React from 'react';
import { Navigate } from 'react-router-dom';
import { authApi } from '../../services/auth';

const PrivateRoute: React.FC<{ children: React.ReactElement }> = ({ children }) => {
  return authApi.isAuthenticated() ? children : <Navigate to="/login" />;
};

export default PrivateRoute;
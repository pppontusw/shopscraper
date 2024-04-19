import { useApiKey } from './ApiKey';
import { Navigate } from 'react-router-dom';

function PrivateRoute({ children }) {
  const { apiKey } = useApiKey();
  return apiKey ? children : <Navigate to="/auth" replace />;
}

export default PrivateRoute
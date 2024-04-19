import { useApiKey } from './ApiKey';
import { ApiKeyInput } from '../components/ApiKeyInput';
import { Navigate } from 'react-router-dom';

function AuthRoute() {
  const { apiKey } = useApiKey();

  return apiKey ? <Navigate to="/" replace /> : <ApiKeyInput />;
}

export default AuthRoute;
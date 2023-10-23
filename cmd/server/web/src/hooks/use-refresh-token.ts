import axios from '@/api/axios';
import useAuth from "@/hooks/use-auth.ts";

const useRefreshToken = () => {
  const { setAuth } = useAuth();

  const refresh = async () => {
    const response = await axios.post('/api/v1/refresh_token', {
      withCredentials: true
    });
    setAuth((prev: any) => {
      console.log(JSON.stringify(prev));
      console.log(response.data.accessToken);
      return { ...prev, accessToken: response.data.accessToken }
    });
    return response.data.accessToken;
  }

  return refresh;
}

export default useRefreshToken;

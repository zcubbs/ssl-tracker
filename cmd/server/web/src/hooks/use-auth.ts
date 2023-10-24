import { useContext } from "react";
import AuthContext from "@/context/auth-provider.tsx";

const useAuth = () => {
  return useContext(AuthContext);
};

export default useAuth;

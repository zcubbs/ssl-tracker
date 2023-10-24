import { ReactElement } from "react";
import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "@/hooks/use-auth.ts";

interface RequireAuthProps {
  allowedRoles: string[];
}

const RequireAuth = ({ allowedRoles }: RequireAuthProps): ReactElement | null => {
  const { auth } = useAuth();
  const location = useLocation();

  console.log("RequireAuth: auth", auth);
  console.log("RequireAuth: allowedRoles", allowedRoles);

  // If the user is not authenticated, redirect to login
  if (!auth?.user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // If the user is authenticated but doesn't have a required role, redirect to unauthorized
  if (!allowedRoles.some((role) => auth.user?.role?.includes(role))) {
    return <Navigate to="/unauthorized" state={{ from: location }} replace />;
  }

  // If the user has a required role, render the outlet (component)
  return <Outlet />;
};

export default RequireAuth;

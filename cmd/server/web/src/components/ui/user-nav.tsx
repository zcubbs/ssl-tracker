import {Avatar, AvatarFallback, AvatarImage,} from "./avatar";
import {Button} from "./button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./dropdown-menu";
import {useContext} from "react";
import AuthContext, {Auth} from "@/context/auth-provider.tsx";
import axios from "@/api/axios.ts";

export function UserNav() {
  const {setAuth} = useContext(AuthContext);

  const logout = async () => {
    try {
      const savedAuthData = localStorage.getItem('authData');

// Check if savedAuthData is not null before parsing
      const authData: Auth = savedAuthData ? JSON.parse(savedAuthData) : null;
      const accessToken = authData ? authData.access_token : '';

      const response = await axios.post('/api/v1/logout_user',
        {
          session_id: authData?.session_id
        },
        {
          headers: {
            Authorization: `Bearer ${accessToken}`, // Authorization header
          }
        });

      if (response.status === 200) {
        // Clearing auth data from localStorage and context
        localStorage.removeItem('authData');
        if (setAuth) {
          setAuth(null);
        }

        // Redirect to the login page or another page as per your application flow
        location.href = '/login';
      } else {
        // Handle unsuccessful logout attempt
        localStorage.removeItem('authData');
        console.error('Logout failed:', response);
      }
    } catch (error) {
      console.error('An error occurred during logout:', error);
    }
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="relative h-8 w-8 rounded-full">
          <Avatar className="h-8 w-8">
            <AvatarImage src={`https://avatar.vercel.sh/1.svg?size=30&text=A`} alt="anonymous"/>
            <AvatarFallback>A</AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" align="end" forceMount>
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">anonymous</p>
            <p className="text-xs leading-none text-muted-foreground">
              m@example.com
            </p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator/>
        <DropdownMenuGroup>
          <DropdownMenuItem>
            Profile
          </DropdownMenuItem>
          <DropdownMenuItem>
            Billing
          </DropdownMenuItem>
          <DropdownMenuItem>
            Settings
          </DropdownMenuItem>
          <DropdownMenuItem>New Space</DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator/>
        <DropdownMenuItem onSelect={logout}>
          Log out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

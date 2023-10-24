import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "./avatar";
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
import AuthContext from "@/context/auth-provider.tsx";
import {useNavigate} from "react-router-dom";
import axios from "@/api/axios.ts";

export function UserNav() {
  const {setAuth} = useContext(AuthContext);
  const navigate = useNavigate();

  const logout = async () => {
    try {
      axios.post('/api/v1/logout', {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
      }).then((response) => {
        if (response.status >= 200 && response.status < 300) {
          // Clearing the authentication context
          if (setAuth) {
            setAuth(undefined);
          }
        } else {
          console.error(response)
        }
      });

      navigate('/');
    } catch (error) {
      // TODO: Handle error, e.g., show a notification to the user
      console.error("An error occurred during logout:", error);
    }
  };

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

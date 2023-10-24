import * as React from "react"

import {cn} from "@/lib/utils"
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {Icons} from "@/components/ui/icons.tsx";
import {SyntheticEvent, useState} from "react";
import {useLocation, useNavigate} from "react-router-dom";
import useAuth from "@/hooks/use-auth.ts";

interface UserLoginFormProps extends React.HTMLAttributes<HTMLDivElement> {
}

export function UserLoginForm({className, ...props}: UserLoginFormProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || "/";
  const [isLoading, setIsLoading] = useState<boolean>(false)

  const [username, setUsername] = useState<string>("")
  const [password, setPassword] = useState<string>("")

  const {setAuth} = useAuth();

  const onSubmit = async (event: SyntheticEvent) => {
    event.preventDefault()
    setIsLoading(true)

    let request = {
      username: username,
      password: password
    }

    await fetch('http://localhost:8000/api/v1/login_user', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(request),
    }).then((response) => {
      if (!response.ok) {
        response.json().then((data) => {
          let responseErrors = []
          let details = data?.details
          for (let key in details) {
            for (let violation of details[key]?.field_violations) {
              responseErrors.push(violation?.field + ": " + violation?.description)
            }
          }
        })
      }

      if (response.ok) {
        response.json().then((data) => {
          if (setAuth) {
            setAuth(data)
          }
        });
        navigate(from, { replace: true });
      }

    }).catch((error) => {
      console.error(error)
    })

    setIsLoading(false)
  }

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <div className="flex flex-col space-y-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">
          Login
        </h1>
        <p className="text-sm text-muted-foreground">
          Connect to your account
        </p>
      </div>
      <form onSubmit={onSubmit}>
        <div className="grid gap-2">
          <div className="grid gap-1">
            <Label className="sr-only" htmlFor="email">
              Username
            </Label>
            <Input
              id="username"
              placeholder="username or email"
              type="text"
              autoCapitalize="username"
              autoComplete="text"
              autoCorrect="off"
              disabled={isLoading}
              onChange={(e) => setUsername(e.target.value)}
            />
            <Label className="sr-only" htmlFor="email">
              Password
            </Label>
            <Input
              id="password"
              placeholder="password"
              type="password"
              autoCapitalize="none"
              autoComplete="password"
              autoCorrect="off"
              disabled={isLoading}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <Button disabled={isLoading}>
            {isLoading && (
              <Icons.spinner className="mr-2 h-4 w-4 animate-spin"/>
            )}
            Login
          </Button>
        </div>
      </form>
    </div>
  )
}

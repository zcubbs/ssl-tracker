import * as React from "react"

import {cn} from "@/lib/utils"
import {Button} from "@/components/ui/button"
import {Input} from "@/components/ui/input"
import {Label} from "@/components/ui/label"
import {Icons} from "@/components/ui/icons.tsx";
import {useState} from "react";

interface UserRegisterFormProps extends React.HTMLAttributes<HTMLDivElement> {
}

export function UserRegisterForm({className, ...props}: UserRegisterFormProps) {
    const [isLoading, setIsLoading] = useState<boolean>(false)

    const [name, setName] = useState<string>("")
    const [email, setEmail] = useState<string>("")
    const [password, setPassword] = useState<string>("")

    async function onSubmit(event: React.SyntheticEvent) {
        event.preventDefault()
        setIsLoading(true)

        let request = {
            full_name: name,
            username: email,
            email: email,
            password: password
        }

        const response = await fetch('http://localhost:8000/api/v1/create_user', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(request),
        })

        const data = await response.json()
        console.log(data)

        setTimeout(() => {
            setIsLoading(false)
        }, 3000)
    }

    return (
        <div className={cn("grid gap-6", className)} {...props}>
            <div className="flex flex-col space-y-2 text-center">
                <h1 className="text-2xl font-semibold tracking-tight">
                    Create an account
                </h1>
                <p className="text-sm text-muted-foreground">
                    Enter your email below to create your account
                </p>
            </div>
            <form onSubmit={onSubmit}>
                <div className="grid gap-2">
                    <div className="grid gap-1">
                        <Label className="sr-only" htmlFor="email">
                            Name
                        </Label>
                        <Input
                            id="name"
                            placeholder="Amy Smith"
                            type="text"
                            autoCapitalize="on"
                            autoComplete="name"
                            autoCorrect="off"
                            disabled={isLoading}
                            onChange={(e) => setName(e.target.value)}
                        />
                        <Label className="sr-only" htmlFor="email">
                            Email
                        </Label>
                        <Input
                            id="email"
                            placeholder="name@example.com"
                            type="email"
                            autoCapitalize="none"
                            autoComplete="email"
                            autoCorrect="off"
                            disabled={isLoading}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <Label className="sr-only" htmlFor="email">
                            Password
                        </Label>
                        <Input
                            id="password"
                            placeholder="**********"
                            type="password"
                            autoCapitalize="none"
                            autoCorrect="off"
                            disabled={isLoading}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <Button disabled={isLoading}>
                        {isLoading && (
                            <Icons.spinner className="mr-2 h-4 w-4 animate-spin"/>
                        )}
                        Sign In with Email
                    </Button>
                </div>
            </form>
            <div className="relative">
                <div className="absolute inset-0 flex items-center">
                    <span className="w-full border-t"/>
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-background px-2 text-muted-foreground">
                    Or continue with
                  </span>
                </div>
            </div>
            <Button variant="outline" type="button" disabled>
                {isLoading ? (
                    <Icons.spinner className="mr-2 h-4 w-4 animate-spin"/>
                ) : (
                    <Icons.gitHub className="mr-2 h-4 w-4"/>
                )}{" "}
                Github
            </Button>
            <Button variant="outline" type="button" disabled>
                {isLoading ? (
                    <Icons.spinner className="mr-2 h-4 w-4 animate-spin"/>
                ) : (
                    <Icons.google className="mr-2 h-4 w-4"/>
                )}{" "}
                Google
            </Button>
        </div>
    )
}

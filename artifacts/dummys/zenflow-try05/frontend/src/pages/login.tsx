// fullend:gen ssot=frontend/login.html contract=2946ce8
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { api } from '../api'

export default function Login() {
  const queryClient = useQueryClient()

  const loginForm = useForm()
  const loginMutation = useMutation({
    mutationFn: (data: any) => api.Login(data),
    onSuccess: () => {
      queryClient.invalidateQueries()
    },
  })

  const registerForm = useForm()
  const registerMutation = useMutation({
    mutationFn: (data: any) => api.Register(data),
    onSuccess: () => {
      queryClient.invalidateQueries()
    },
  })

  return (
    <div>
      <title>ZenFlow Login</title>
      <form onSubmit={loginForm.handleSubmit((data) => loginMutation.mutate(data))}>
        <input type="email" placeholder="Email" {...loginForm.register('email')} />
        <input type="password" placeholder="Password" {...loginForm.register('password')} />
        <button type="submit">Login</button>
      </form>
      <form onSubmit={registerForm.handleSubmit((data) => registerMutation.mutate(data))}>
        <input type="text" placeholder="Organization Name" {...registerForm.register('org_name')} />
        <input type="email" placeholder="Email" {...registerForm.register('email')} />
        <input type="password" placeholder="Password" {...registerForm.register('password')} />
        <button type="submit">Register</button>
      </form>
    </div>
  )
}

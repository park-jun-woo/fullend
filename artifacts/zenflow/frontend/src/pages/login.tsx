// fullend:gen ssot=frontend/login.html contract=f8df1fe
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

  return (
    <div>
      <title>ZenFlow - Login</title>
      <h1>Login</h1>
      <form onSubmit={loginForm.handleSubmit((data) => loginMutation.mutate(data))}>
        <input type="email" placeholder="Email" {...loginForm.register('email')} />
        <input type="password" placeholder="Password" {...loginForm.register('password')} />
        <button type="submit">Login</button>
      </form>
    </div>
  )
}

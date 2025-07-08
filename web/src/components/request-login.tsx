import { LogIn } from 'lucide-react'

export const RequestLogin = () => (
  <div className="text-center py-12">
    <LogIn className="h-12 w-12 text-gray-400 mx-auto mb-4" />
    <h3 className="text-lg font-medium text-gray-900 mb-2">
      You're not logged in
    </h3>
    <p className="text-gray-500">
      Log in with your GitHub account to view your repositories
    </p>
  </div>
)

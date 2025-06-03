import { Button } from "@/components/ui/button"
import { Shield, Github } from "lucide-react"

export const Header = () => {
  return (
    <header className="backdrop-blur-xl bg-black/20 border-b border-white/10">
      <div className="max-w-7xl mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg flex items-center justify-center">
              <Shield className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                OSS Health
              </h1>
              <p className="text-sm text-gray-400">Dependency Security & Health Monitoring</p>
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <Button className="px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg text-white transition-all duration-300 backdrop-blur-sm">
              <Github className="w-4 h-4 inline mr-2" />
              Connect GitHub
            </Button>
          </div>
        </div>
      </div>
    </header>
  )
}

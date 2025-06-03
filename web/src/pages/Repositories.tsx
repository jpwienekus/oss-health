import { useGitHubPopupLogin } from "@/auth/useGitHubPopupLogin"
// import { RequestGithubLogin } from "@/components/repositories/RequestGithubLogin"
import { getClient } from "@/graphql/client"
// import { SYNC_REPOS } from "@/graphql/mutations"
import { GET_REPOS } from "@/graphql/queries"
import type { Repository } from "@/types"
import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { AlertTriangle, Clock, Database, Download, Github, Star } from "lucide-react"

export const Repositories = () => {
  const { jwt } = useGitHubPopupLogin()
  const [data, setData] = useState<Repository[]>([])
  const [syncDate, setSyncDate] = useState<Date | undefined>(undefined)

  const fetchRepositories = async () => {
    if (!jwt) {
      return
    }

    const client = getClient(jwt)
    const response = await client.request<{ repositories: { repositories: Repository[], syncDate: string } }>(GET_REPOS)
    setData(response.repositories.repositories)
    setSyncDate(response.repositories.syncDate ? new Date(response.repositories.syncDate) : undefined)

  }

  // const syncRepositories = async () => {
  //   if (!jwt) {
  //     return
  //   }
  //
  //   const client = getClient(jwt)
  //   const response = await client.request<{ syncRepositories: { repositories: Repository[], syncDate: string } }>(SYNC_REPOS)
  //   setData(response.syncRepositories.repositories)
  //   setSyncDate(response.syncRepositories.syncDate ? new Date(response.syncRepositories.syncDate) : undefined)
  // }

  useEffect(() => {
    fetchRepositories()
  }, []);

  const getScoreColor = (score: number) => {
    if (score >= 80) {
      return 'text-green-400'
    }
    if (score >= 60) {
      return 'text-yellow-400'
    }
    return 'text-red-400'
  }

  const getScoreBg = (score: number) => {
    if (score >= 80) {
      return 'bg-green-500/20 border-green-500/30'
    }
    if (score >= 60) {
      return 'bg-yellow-500/20 border-yellow-500/30'
    }
    return 'bg-red-500/20 border-red-500/30'
  }

  const formatDate = (date: Date | undefined) => {
    if (!date) {
      return '-'
    }

    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    })
  }

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold text-white">Health Reports</h2>
          <p className="text-gray-400">Monitor your project dependencies</p>
        </div>
        <div className="flex space-x-3">
          <Button className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors">
            <Download className="w-4 h-4 inline mr-2" />
            Export HTML
          </Button>
          <Button className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors">
            <Download className="w-4 h-4 inline mr-2" />
            Export MD
          </Button>
        </div>
      </div>

      <div className="grid gap-6">
        {data.map((repository, i) => (
          <div key={i} className="bg-black/30 backdrop-blur-xl border border-white/10 rounded-2xl p-6 hover:border-white/20 transition-all duration-300">
            <div className="flex justify-between items-start mb-6">
              <div>
                <h3 className="text-xl font-semibold text-white mb-2">{repository.name}</h3>
                <p className="text-gray-400">Scanned on {formatDate(new Date())}</p>
              </div>
              <div className={`px-4 py-2 rounded-xl border ${getScoreBg(2)}`}>
                <span className={`text-2xl font-bold ${getScoreColor(2)}`}>2</span>
                <span className="text-gray-400 text-sm ml-1">/100</span>
              </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-6 md-6">
              <div className="text-center">
                <div className="w-12  h-12 bg-blue-500/20 rounded-lg flex items-center justify-center mx-auto mb-2">
                  <Database className="w-6 h-6 text-blue-400" />
                </div>
                <div className="text-2xl font-bold text-white">33</div>
                <div className="text-sm text-gray-400">Dependencies</div>
              </div>

              <div className="text-center">
                <div className="w-12  h-12 bg-red-500/20 rounded-lg flex items-center justify-center mx-auto mb-2">
                  <AlertTriangle className="w-6 h-6 text-red-400" />
                </div>
                <div className="text-2xl font-bold text-white">3</div>
                <div className="text-sm text-gray-400">Vulnerabilities</div>
              </div>

              <div className="text-center">
                <div className="w-12  h-12 bg-yellow-500/20 rounded-lg flex items-center justify-center mx-auto mb-2">
                  <Star className="w-6 h-6 text-yellow-400" />
                </div>
                <div className="text-2xl font-bold text-white">3</div>
                <div className="text-sm text-gray-400">Stars</div>
              </div>

              <div className="text-center">
                <div className="w-12  h-12 bg-green-500/20 rounded-lg flex items-center justify-center mx-auto mb-2">
                  <Clock className="w-6 h-6 text-green-400" />
                </div>
                <div className="text-2xl font-bold text-white">3</div>
                <div className="text-sm text-gray-400">Last Commit</div>
              </div>
            </div>

            <div className="flex space-x-4 mt-3">
              <Button className="flex-1 py-3 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white rounded-lg font-medium transition-all duration-300">
                View Detailed Report
              </Button>
              <Button className="px-6 py-3 bg-white/10 hover:bg-white/20 border border-white/20 text-white rounder-lg transition-all duration-300">
                <Github className="w-4 h-4" />
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
